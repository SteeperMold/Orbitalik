use std::num::TryFromIntError;
use chrono::{DateTime, Duration, Utc};
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

            // aos detection
            if let (Some(prev), None) = (&prev_sample, &current_pass)
                && Self::crossed_horizon(prev, &sample, min_el_rad)
            {
                current_pass = Some(PassState {
                    aos: *prev,
                    max: sample,
                });
            }

            // inside pass
            if let Some(pass) = &mut current_pass {
                if sample.elevation_rad >= min_el_rad {
                    Self::update_peak(pass, &sample);
                } else {
                    // los reached
                    let final_pass = Self::build_pass_from_state(
                        SatelliteIdentifier::NoradId(self.norad_id),
                        pass,
                        sample,
                    )?;

                    if final_pass.max_elevation.get::<radian>() < min_peak_el_rad {
                        current_pass = None;
                        prev_sample = Some(sample);
                        t += Duration::seconds(COARSE_STEP_SECONDS);
                        continue;
                    }

                    if Self::decrement_remaining(remaining).is_err() {
                        break;
                    }

                    passes.push(final_pass);
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
