use chrono::{DateTime, Duration, Utc};
use uom::si::angle::radian;
use uom::si::f64::Angle;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::errors::PropagationError;
use crate::astro::models::TimeRange;
use crate::astro::models::{Pass, SatelliteIdentifier};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::propagator::Propagator;

#[derive(Clone, Copy)]
struct ElevationSample {
    time: DateTime<Utc>,
    azimuth_rad: f64,
    elevation_rad: f64,
}

pub struct PassPredictionOptions<'a> {
    pub range: TimeRange,
    pub observer: &'a Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub max_results: Option<usize>,
}

const COARSE_STEP_SECONDS: i64 = 60;

impl Propagator {
    pub fn predict_passes(
        &self,
        options: &PassPredictionOptions,
    ) -> Result<Vec<Pass>, PropagationError> {
        let min_el = options.min_elevation.get::<radian>();
        let min_peak_el = options.min_peak_elevation.get::<radian>();

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: true,
            range: true,
        };
        let mut samples: Vec<ElevationSample> = Vec::new();

        let mut t = options.range.start;
        while t <= options.range.end {
            let la = self.look_angles_at(t, options.observer, &compute)?;

            let (Some(az), Some(el)) = (la.azimuth, la.elevation) else {
                t += Duration::seconds(COARSE_STEP_SECONDS);
                continue;
            };

            samples.push(ElevationSample {
                time: t,
                azimuth_rad: az.get::<radian>(),
                elevation_rad: el.get::<radian>(),
            });

            t += Duration::seconds(COARSE_STEP_SECONDS);
        }

        let mut passes = Vec::new();
        let mut i = 1;

        while i < samples.len() {
            let prev = samples[i - 1];
            let curr = samples[i];

            if !(prev.elevation_rad < min_el && curr.elevation_rad >= min_el) {
                i += 1;
                continue;
            }

            // AOS: crossing upward through min elevation
            let aos_idx = i - 1;

            let mut max_idx = aos_idx;
            let mut j = i;

            // follow pass until LOS
            while j < samples.len() && samples[j].elevation_rad >= min_el {
                if samples[j].elevation_rad > samples[max_idx].elevation_rad {
                    max_idx = j;
                }
                j += 1;
            }

            if j >= samples.len() {
                break;
            }

            let los_idx = j;

            if samples[max_idx].elevation_rad < min_peak_el {
                i = j;
                continue;
            }

            let pass = Self::build_pass(
                SatelliteIdentifier::NoradId(self.norad_id),
                &samples,
                aos_idx,
                max_idx,
                los_idx,
            );

            passes.push(pass);

            if let Some(max) = options.max_results
                && passes.len() >= max
            {
                break;
            }

            i = j;
        }

        Ok(passes)
    }

    fn build_pass(
        satellite: SatelliteIdentifier,
        s: &[ElevationSample],
        aos_idx: usize,
        max_idx: usize,
        los_idx: usize,
    ) -> Pass {
        let aos = s[aos_idx];
        let max = s[max_idx];
        let los = s[los_idx];

        let duration = (los.time - aos.time).num_seconds().max(0) as u32;

        Pass {
            satellite,

            aos: aos.time,
            aos_azimuth: aos.azimuth_rad.to_degrees(),

            max_elevation_time: max.time,
            max_elevation: max.elevation_rad.to_degrees(),
            max_elevation_azimuth: max.azimuth_rad.to_degrees(),

            los: los.time,
            los_azimuth: los.azimuth_rad.to_degrees(),

            duration_seconds: duration,
        }
    }
}
