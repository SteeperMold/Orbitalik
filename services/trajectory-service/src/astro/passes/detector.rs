use chrono::{DateTime, Duration, Utc};
use std::num::TryFromIntError;
use std::sync::atomic::{AtomicUsize, Ordering};
use uom::si::angle::radian;
use uom::si::f64::Angle;

use crate::astro::errors::PropagationError;
use crate::astro::models::{Pass, SatelliteIdentifier};
use crate::astro::passes::predictor::PassPredictionOptions;
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::propagator::Propagator;

#[derive(Clone, Copy)]
struct ElevationSample {
    time: DateTime<Utc>,
    azimuth_rad: f64,
    elevation_rad: f64,
}

struct PassState {
    aos: ElevationSample,
    max: ElevationSample,
}

const COARSE_STEP_SECONDS: i64 = 60;
const CROSSING_TOLERANCE_MILLISECONDS: i64 = 10;

impl Propagator {
    pub fn predict_passes(
        &self,
        options: &PassPredictionOptions,
        remaining: &AtomicUsize,
    ) -> Result<Vec<Pass>, PropagationError> {
        let min_el_rad = options.min_elevation.get::<radian>();
        let min_peak_el_rad = options.min_peak_elevation.get::<radian>();

        let mut passes = Vec::new();
        let mut prev_sample: Option<ElevationSample> = None;
        let mut current_pass: Option<PassState> = None;

        let mut t = options.range.start;

        while t <= options.range.end && remaining.load(Ordering::Relaxed) > 0 {
            let Some(sample) = self.sample(t, options)? else {
                t += Duration::seconds(COARSE_STEP_SECONDS);
                continue;
            };

            // first sample is already inside a pass
            if prev_sample.is_none() && sample.elevation_rad >= min_el_rad {
                current_pass = Some(PassState {
                    aos: sample,
                    max: sample,
                });
            }

            // normal aos detection
            if let (Some(prev), None) = (&prev_sample, &current_pass)
                && Self::crossed_horizon(prev, &sample, min_el_rad)
            {
                let aos = self.find_horizon_crossing(prev, &sample, min_el_rad, options)?;

                current_pass = Some(PassState { aos, max: sample });
            }

            // inside pass
            if let Some(pass) = &mut current_pass {
                if sample.elevation_rad >= min_el_rad {
                    Self::update_peak(pass, &sample);
                } else if let Some(prev) = prev_sample {
                    // los reached
                    let los = self.find_horizon_crossing(&sample, &prev, min_el_rad, options)?;

                    let final_pass = Self::build_pass_from_state(
                        SatelliteIdentifier::NoradId(self.norad_id),
                        pass,
                        los,
                    )?;

                    // skip passes with low peak elevation
                    if final_pass.max_elevation.get::<radian>() >= min_peak_el_rad {
                        if Self::decrement_remaining(remaining).is_err() {
                            break;
                        }

                        passes.push(final_pass);
                    }

                    current_pass = None;
                }
            }

            prev_sample = Some(sample);
            t += Duration::seconds(COARSE_STEP_SECONDS);
        }

        Ok(passes)
    }

    fn sample(
        &self,
        t: DateTime<Utc>,
        options: &PassPredictionOptions,
    ) -> Result<Option<ElevationSample>, PropagationError> {
        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: true,
            range: true,
        };

        let la = self.look_angles_at(t, options.observer, &compute)?;

        Ok(match (la.azimuth, la.elevation) {
            (Some(az), Some(el)) => Some(ElevationSample {
                time: t,
                azimuth_rad: az.get::<radian>(),
                elevation_rad: el.get::<radian>(),
            }),
            _ => None,
        })
    }

    fn crossed_horizon(prev: &ElevationSample, curr: &ElevationSample, min_el_rad: f64) -> bool {
        prev.elevation_rad < min_el_rad && curr.elevation_rad >= min_el_rad
    }

    fn update_peak(pass: &mut PassState, sample: &ElevationSample) {
        if sample.elevation_rad > pass.max.elevation_rad {
            pass.max = *sample;
        }
    }

    fn decrement_remaining(remaining: &AtomicUsize) -> Result<usize, usize> {
        remaining.fetch_update(Ordering::Relaxed, Ordering::Relaxed, |x| {
            if x > 0 { Some(x - 1) } else { None }
        })
    }

    fn find_horizon_crossing(
        &self,
        before: &ElevationSample,
        after: &ElevationSample,
        min_el_rad: f64,
        options: &PassPredictionOptions,
    ) -> Result<ElevationSample, PropagationError> {
        let rising = before.elevation_rad < min_el_rad && after.elevation_rad >= min_el_rad;

        let mut low = *before;
        let mut high = *after;

        while (high.time - low.time).num_milliseconds().abs() > CROSSING_TOLERANCE_MILLISECONDS {
            let midpoint = low.time + (high.time - low.time) / 2;

            let Some(mid) = self.sample(midpoint, options)? else {
                break;
            };

            if rising {
                if mid.elevation_rad >= min_el_rad {
                    high = mid;
                } else {
                    low = mid;
                }
            } else if mid.elevation_rad < min_el_rad {
                high = mid;
            } else {
                low = mid;
            }
        }

        Ok(if rising { high } else { low })
    }

    fn build_pass_from_state(
        satellite: SatelliteIdentifier,
        state: &PassState,
        los: ElevationSample,
    ) -> Result<Pass, TryFromIntError> {
        let duration = u32::try_from((los.time - state.aos.time).num_seconds().max(0))?;

        Ok(Pass {
            satellite,

            aos: state.aos.time,
            aos_azimuth: Angle::new::<radian>(state.aos.azimuth_rad),

            max_elevation_time: state.max.time,
            max_elevation: Angle::new::<radian>(state.max.elevation_rad),
            max_elevation_azimuth: Angle::new::<radian>(state.max.azimuth_rad),

            los: los.time,
            los_azimuth: Angle::new::<radian>(los.azimuth_rad),

            duration_seconds: duration,
        })
    }
}

#[allow(
    clippy::unwrap_used,
    clippy::panic,
    clippy::similar_names,
    clippy::cast_precision_loss
)]
#[cfg(test)]
mod tests {
    use super::*;

    use chrono::{DateTime, TimeZone, Utc};
    use uom::si::angle::degree;

    use crate::astro::models::SatelliteIdentifier;
    use crate::astro::test_utils::assert_close;

    fn time(seconds: i64) -> DateTime<Utc> {
        Utc.timestamp_opt(seconds, 0).single().unwrap()
    }

    fn sample(seconds: i64, azimuth_deg: f64, elevation_deg: f64) -> ElevationSample {
        ElevationSample {
            time: time(seconds),
            azimuth_rad: azimuth_deg.to_radians(),
            elevation_rad: elevation_deg.to_radians(),
        }
    }

    #[test]
    fn crossed_horizon_when_crossing_upward() {
        let prev = sample(0, 100.0, -1.0);
        let curr = sample(60, 110.0, 1.0);

        assert!(Propagator::crossed_horizon(
            &prev,
            &curr,
            0.0_f64.to_radians(),
        ));
    }

    #[test]
    fn crossed_horizon_when_both_below() {
        let prev = sample(0, 100.0, -10.0);
        let curr = sample(60, 110.0, -1.0);

        assert!(!Propagator::crossed_horizon(
            &prev,
            &curr,
            0.0_f64.to_radians(),
        ));
    }

    #[test]
    fn crossed_horizon_when_both_above() {
        let prev = sample(0, 100.0, 1.0);
        let curr = sample(60, 110.0, 10.0);

        assert!(!Propagator::crossed_horizon(
            &prev,
            &curr,
            0.0_f64.to_radians(),
        ));
    }

    #[test]
    fn crossing_exactly_at_threshold_counts_as_crossing() {
        let prev = sample(0, 100.0, -1.0);
        let curr = sample(60, 110.0, 0.0);

        assert!(Propagator::crossed_horizon(
            &prev,
            &curr,
            0.0_f64.to_radians(),
        ));
    }

    #[test]
    fn crossing_downward_does_not_count_as_aos() {
        let prev = sample(0, 100.0, 1.0);
        let curr = sample(60, 110.0, -1.0);

        assert!(!Propagator::crossed_horizon(
            &prev,
            &curr,
            0.0_f64.to_radians(),
        ));
    }

    #[test]
    fn update_peak_replaces_lower_peak() {
        let aos = sample(0, 90.0, 1.0);
        let initial_max = sample(60, 100.0, 20.0);

        let mut state = PassState {
            aos,
            max: initial_max,
        };

        let new_peak = sample(120, 110.0, 40.0);

        Propagator::update_peak(&mut state, &new_peak);

        assert_close(state.max.elevation_rad, 40.0_f64.to_radians(), 1e-12);

        assert_eq!(state.max.time, new_peak.time);
        assert_close(state.max.azimuth_rad, 110.0_f64.to_radians(), 1e-12);
    }

    #[test]
    fn update_peak_does_not_replace_higher_peak() {
        let aos = sample(0, 90.0, 1.0);
        let initial_max = sample(60, 100.0, 40.0);

        let mut state = PassState {
            aos,
            max: initial_max,
        };

        let lower_peak = sample(120, 110.0, 20.0);

        Propagator::update_peak(&mut state, &lower_peak);

        assert_close(state.max.elevation_rad, 40.0_f64.to_radians(), 1e-12);

        assert_eq!(state.max.time, initial_max.time);
    }

    #[test]
    fn decrement_remaining_decrements_positive_value() {
        let remaining = AtomicUsize::new(3);

        let result = Propagator::decrement_remaining(&remaining);

        assert_eq!(result, Ok(3));
        assert_eq!(remaining.load(Ordering::Relaxed), 2);
    }

    #[test]
    fn decrement_remaining_reaches_zero() {
        let remaining = AtomicUsize::new(1);

        let result = Propagator::decrement_remaining(&remaining);

        assert_eq!(result, Ok(1));
        assert_eq!(remaining.load(Ordering::Relaxed), 0);
    }

    #[test]
    fn decrement_remaining_does_not_underflow() {
        let remaining = AtomicUsize::new(0);

        let result = Propagator::decrement_remaining(&remaining);

        assert_eq!(result, Err(0));
        assert_eq!(remaining.load(Ordering::Relaxed), 0);
    }

    #[test]
    fn build_pass_contains_correct_values() {
        let aos = sample(100, 90.0, 5.0);
        let max = sample(160, 120.0, 45.0);
        let los = sample(220, 150.0, -1.0);

        let state = PassState { aos, max };

        let pass =
            Propagator::build_pass_from_state(SatelliteIdentifier::NoradId(12345), &state, los)
                .unwrap();

        assert_eq!(pass.aos, time(100));
        assert_eq!(pass.max_elevation_time, time(160));
        assert_eq!(pass.los, time(220));

        assert_close(pass.aos_azimuth.get::<degree>(), 90.0, 1e-12);
        assert_close(pass.max_elevation.get::<degree>(), 45.0, 1e-12);
        assert_close(pass.max_elevation_azimuth.get::<degree>(), 120.0, 1e-12);
        assert_close(pass.los_azimuth.get::<degree>(), 150.0, 1e-12);

        assert_eq!(pass.duration_seconds, 120);

        match pass.satellite {
            SatelliteIdentifier::NoradId(id) => assert_eq!(id, 12345),
            SatelliteIdentifier::Name(_) => panic!("expected NORAD identifier"),
        }
    }

    // TODO return this test
    // #[test]
    // fn predicts_reference_pass() {
    //     let propagator = test_propagator();
    //     let observer = test_observer();
    //
    //     let range = TimeRange {
    //         start: test_datetime(),
    //         end: test_datetime() + TimeDelta::hours(6),
    //     };
    //
    //     let options = PassPredictionOptions {
    //         range,
    //         observer: &observer,
    //         min_elevation: Angle::new::<degree>(0.0),
    //         min_peak_elevation: Angle::new::<degree>(0.0),
    //         max_results: None,
    //     };
    //
    //     let remaining = AtomicUsize::new(10);
    //
    //     let passes = propagator.predict_passes(&options, &remaining).unwrap();
    //
    //     assert!(!passes.is_empty());
    //     assert_eq!(passes.len(), 3);
    //
    //     let pass = &passes[0];
    //
    //     let reference_aos: DateTime<Utc> = "2026-01-01T00:51:42.061040+00:00".parse().unwrap();
    //     let reference_los: DateTime<Utc> = "2026-01-01T01:02:30.158873+00:00".parse().unwrap();
    //
    //     assert_close((pass.aos - reference_aos).num_milliseconds() as f64, 0.0, 10.0);
    //     assert_close((pass.los - reference_los).num_milliseconds() as f64, 0.0, 10.0);
    //
    //     assert_close(
    //         pass.max_elevation.get::<degree>(),
    //         33.550_290_464_261_14,
    //         0.1,
    //     );
    // }
}
