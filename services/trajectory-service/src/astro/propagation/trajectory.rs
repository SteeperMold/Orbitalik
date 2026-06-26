use chrono::Duration;
use rayon::prelude::{IntoParallelIterator, ParallelIterator};

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{ObserverTrajectory, Trajectory};
use crate::astro::models::{Sampling, TimeRange};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::astro::propagation::propagator::Propagator;
use crate::astro::errors::PropagationError;

impl Propagator {
    pub fn trajectory_at(
        &self,
        range: TimeRange,
        sampling: Sampling,
        compute: &PositionComputation,
    ) -> Result<Trajectory, PropagationError> {
        let step = i64::from(sampling.step_seconds);
        let total_secs = (range.end - range.start).num_seconds();
        let steps = (total_secs / step) + 1;

        let samples: Result<Vec<_>, _> = (0..steps)
            .into_par_iter()
            .map(|i| {
                let t = range.start + Duration::seconds(i * step);
                self.position_at(t, compute)
            })
            .collect();

        Ok(Trajectory {
            start: range.start,
            end: range.end,
            step_seconds: sampling.step_seconds,
            samples: samples?,
        })
    }

    pub fn observer_trajectory_at(
        &self,
        range: TimeRange,
        sampling: Sampling,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<ObserverTrajectory, PropagationError> {
        let step = i64::from(sampling.step_seconds);
        let total_secs = (range.end - range.start).num_seconds();
        let steps = (total_secs / step) + 1;

        let samples: Result<Vec<_>, _> = (0..steps)
            .into_par_iter()
            .map(|i| {
                let t = range.start + Duration::seconds(i * step);
                self.look_angles_at(t, observer, compute)
            })
            .collect();

        Ok(ObserverTrajectory {
            start: range.start,
            end: range.end,
            step_seconds: sampling.step_seconds,
            samples: samples?,
        })
    }
}
