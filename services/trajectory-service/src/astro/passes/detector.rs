use crate::astro::errors::PropagationError;
use crate::astro::models::{Pass, SatelliteIdentifier};
use crate::astro::passes::predictor::PassPredictionOptions;
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::propagator::Propagator;
use chrono::{DateTime, Duration, Utc};
use std::sync::atomic::{AtomicUsize, Ordering};
use uom::si::angle::radian;
use uom::si::f64::Angle;

#[derive(Clone, Copy)]
struct ElevationSample {
    time: DateTime<Utc>,
    azimuth_rad: f64,
    elevation_rad: f64,
}

struct PassState {
    aos: ElevationSample,
    max: ElevationSample,
    max_el_rad: f64,
}

const COARSE_STEP_SECONDS: i64 = 60;

impl Propagator {
    pub fn predict_passes(
        &self,
        options: &PassPredictionOptions,
        remaining: &AtomicUsize, // global quota
    ) -> Result<Vec<Pass>, PropagationError> {
        let min_el_rad = options.min_elevation.get::<radian>();
        let min_peak_el_rad = options.min_peak_elevation.get::<radian>();

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: true,
            range: true,
        };

        let mut passes = Vec::new();

        let mut prev_sample: Option<ElevationSample> = None;
        let mut current_pass: Option<PassState> = None;

        let mut t = options.range.start;

        while t <= options.range.end {
            if remaining.load(Ordering::Relaxed) == 0 {
                break;
            }

            let la = self.look_angles_at(t, options.observer, &compute)?;

            let (Some(az), Some(el)) = (la.azimuth, la.elevation) else {
                t += Duration::seconds(COARSE_STEP_SECONDS);
                continue;
            };

            let sample = ElevationSample {
                time: t,
                azimuth_rad: az.get::<radian>(),
                elevation_rad: el.get::<radian>(),
            };

            // aos detection
            if let Some(prev) = prev_sample
                && current_pass.is_none()
                && prev.elevation_rad < min_el_rad
                && sample.elevation_rad >= min_el_rad
            {
                current_pass = Some(PassState {
                    aos: prev,
                    max: sample,
                    max_el_rad: sample.elevation_rad,
                });
            }

            // inside pass tracking
            if let Some(pass) = &mut current_pass {
                if sample.elevation_rad >= min_el_rad {
                    if sample.elevation_rad > pass.max_el_rad {
                        pass.max = sample;
                        pass.max_el_rad = sample.elevation_rad;
                    }
                } else {
                    // los reached
                    let final_pass = Self::build_pass_from_state(
                        SatelliteIdentifier::NoradId(self.norad_id),
                        pass,
                        sample,
                    );

                    if final_pass.max_elevation.get::<radian>() < min_peak_el_rad {
                        current_pass = None;
                        prev_sample = Some(sample);
                        t += Duration::seconds(COARSE_STEP_SECONDS);
                        continue;
                    }

                    let old = remaining.fetch_update(Ordering::Relaxed, Ordering::Relaxed, |x| {
                        if x > 0 { Some(x - 1) } else { None }
                    });

                    if old.is_err() {
                        break; // no quota left
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

    fn build_pass_from_state(
        satellite: SatelliteIdentifier,
        state: &PassState,
        los: ElevationSample,
    ) -> Pass {
        let duration = (los.time - state.aos.time).num_seconds().max(0) as u32;

        Pass {
            satellite,

            aos: state.aos.time,
            aos_azimuth: Angle::new::<radian>(state.aos.azimuth_rad),

            max_elevation_time: state.max.time,
            max_elevation: Angle::new::<radian>(state.max.elevation_rad),
            max_elevation_azimuth: Angle::new::<radian>(state.max.azimuth_rad),

            los: los.time,
            los_azimuth: Angle::new::<radian>(los.azimuth_rad),

            duration_seconds: duration,
        }
    }
}
