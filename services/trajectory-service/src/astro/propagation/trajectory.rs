use chrono::Duration;
use rayon::prelude::{IntoParallelIterator, ParallelIterator};

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::errors::PropagationError;
use crate::astro::models::{ObserverTrajectory, Trajectory};
use crate::astro::models::{Sampling, TimeRange};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::astro::propagation::propagator::Propagator;

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

#[allow(clippy::unwrap_used, clippy::cast_possible_wrap)]
#[cfg(test)]
mod tests {
    use chrono::Duration;

    use crate::astro::models::{Sampling, TimeRange};
    use crate::astro::propagation::look_angles::LookAnglesComputation;
    use crate::astro::propagation::position::PositionComputation;
    use crate::astro::test_utils::{assert_close, test_datetime, test_observer, test_propagator};

    use uom::si::angle::degree;
    use uom::si::length::kilometer;

    fn test_range() -> TimeRange {
        let start = test_datetime();

        TimeRange {
            start,
            end: start + Duration::minutes(10),
        }
    }

    fn sampling(step_seconds: u32) -> Sampling {
        Sampling { step_seconds }
    }

    #[test]
    fn trajectory_returns_correct_metadata() {
        let propagator = test_propagator();
        let range = test_range();
        let sampling = sampling(60);

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let result = propagator.trajectory_at(range, sampling, &compute).unwrap();

        assert_eq!(result.start, range.start);
        assert_eq!(result.end, range.end);
        assert_eq!(result.step_seconds, 60);
    }

    #[test]
    fn trajectory_contains_expected_number_of_samples() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 11);
    }

    #[test]
    fn trajectory_contains_single_sample_for_zero_length_range() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let range = TimeRange {
            start: datetime,
            end: datetime,
        };

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 1);
        assert_eq!(result.samples[0].time, datetime);
    }

    #[test]
    fn trajectory_samples_are_at_correct_times() {
        let propagator = test_propagator();
        let start = test_datetime();

        let range = TimeRange {
            start,
            end: start + Duration::minutes(5),
        };

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 6);

        for (i, sample) in result.samples.iter().enumerate() {
            let expected = start + Duration::seconds(i as i64 * 60);

            assert_eq!(sample.time, expected);
        }
    }

    #[test]
    fn trajectory_samples_are_ordered() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator
            .trajectory_at(range, sampling(30), &compute)
            .unwrap();

        assert!(
            result
                .samples
                .windows(2)
                .all(|window| window[0].time < window[1].time)
        );
    }

    #[test]
    fn trajectory_with_different_sampling_steps() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result_60 = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        let result_120 = propagator
            .trajectory_at(range, sampling(120), &compute)
            .unwrap();

        let result_300 = propagator
            .trajectory_at(range, sampling(300), &compute)
            .unwrap();

        assert_eq!(result_60.samples.len(), 11);
        assert_eq!(result_120.samples.len(), 6);
        assert_eq!(result_300.samples.len(), 3);
    }

    #[test]
    fn trajectory_samples_are_separated_by_sampling_interval() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator
            .trajectory_at(range, sampling(90), &compute)
            .unwrap();

        for window in result.samples.windows(2) {
            let delta = window[1].time - window[0].time;

            assert_eq!(delta, Duration::seconds(90));
        }
    }

    #[test]
    fn trajectory_does_not_modify_requested_coordinates() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: false,
            ecef: false,
            geodetic: true,
        };

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert!(!result.samples.is_empty());

        for sample in &result.samples {
            assert!(sample.teme.is_none());
            assert!(sample.ecef.is_none());
            assert!(sample.geodetic.is_some());
        }
    }

    #[test]
    fn trajectory_with_no_computation_returns_empty_samples() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation::default();

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 11);

        for sample in &result.samples {
            assert!(sample.teme.is_none());
            assert!(sample.ecef.is_none());
            assert!(sample.geodetic.is_none());
        }
    }

    #[test]
    fn trajectory_with_all_computations_returns_all_coordinates() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        for sample in &result.samples {
            assert!(sample.teme.is_some());
            assert!(sample.ecef.is_some());
            assert!(sample.geodetic.is_some());
        }
    }

    #[test]
    fn trajectory_positions_match_individual_position_at_calls() {
        let propagator = test_propagator();
        let start = test_datetime();

        let range = TimeRange {
            start,
            end: start + Duration::minutes(2),
        };

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let trajectory = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert_eq!(trajectory.samples.len(), 3);

        for (i, sample) in trajectory.samples.iter().enumerate() {
            let datetime = start + Duration::seconds(i as i64 * 60);

            let expected = propagator.position_at(datetime, &compute).unwrap();

            assert_eq!(sample.time, expected.time);

            let actual_teme = sample.teme.as_ref().unwrap();
            let expected_teme = expected.teme.as_ref().unwrap();

            assert_close(
                actual_teme.x.get::<kilometer>(),
                expected_teme.x.get::<kilometer>(),
                1e-10,
            );

            assert_close(
                actual_teme.y.get::<kilometer>(),
                expected_teme.y.get::<kilometer>(),
                1e-10,
            );

            assert_close(
                actual_teme.z.get::<kilometer>(),
                expected_teme.z.get::<kilometer>(),
                1e-10,
            );

            let actual_ecef = sample.ecef.as_ref().unwrap();
            let expected_ecef = expected.ecef.as_ref().unwrap();

            assert_close(
                actual_ecef.x.get::<kilometer>(),
                expected_ecef.x.get::<kilometer>(),
                1e-10,
            );

            assert_close(
                actual_ecef.y.get::<kilometer>(),
                expected_ecef.y.get::<kilometer>(),
                1e-10,
            );

            assert_close(
                actual_ecef.z.get::<kilometer>(),
                expected_ecef.z.get::<kilometer>(),
                1e-10,
            );

            let actual_geodetic = sample.geodetic.as_ref().unwrap();
            let expected_geodetic = expected.geodetic.as_ref().unwrap();

            assert_close(
                actual_geodetic.lat.get::<degree>(),
                expected_geodetic.lat.get::<degree>(),
                1e-10,
            );

            assert_close(
                actual_geodetic.lon.get::<degree>(),
                expected_geodetic.lon.get::<degree>(),
                1e-10,
            );

            assert_close(
                actual_geodetic.alt.get::<kilometer>(),
                expected_geodetic.alt.get::<kilometer>(),
                1e-10,
            );
        }
    }

    #[test]
    fn trajectory_positions_change_over_time() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        let first = result.samples[0].teme.as_ref().unwrap();
        let last = result.samples.last().unwrap().teme.as_ref().unwrap();

        assert!(
            (first.x.get::<kilometer>() - last.x.get::<kilometer>()).abs() > 1e-6
                || (first.y.get::<kilometer>() - last.y.get::<kilometer>()).abs() > 1e-6
                || (first.z.get::<kilometer>() - last.z.get::<kilometer>()).abs() > 1e-6
        );
    }

    #[test]
    fn trajectory_is_deterministic() {
        let propagator = test_propagator();
        let range = test_range();

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let first = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        let second = propagator
            .trajectory_at(range, sampling(60), &compute)
            .unwrap();

        assert_eq!(first.samples.len(), second.samples.len());

        for (a, b) in first.samples.iter().zip(second.samples.iter()) {
            assert_eq!(a.time, b.time);

            let a_teme = a.teme.as_ref().unwrap();
            let b_teme = b.teme.as_ref().unwrap();

            assert_close(
                a_teme.x.get::<kilometer>(),
                b_teme.x.get::<kilometer>(),
                1e-12,
            );

            assert_close(
                a_teme.y.get::<kilometer>(),
                b_teme.y.get::<kilometer>(),
                1e-12,
            );

            assert_close(
                a_teme.z.get::<kilometer>(),
                b_teme.z.get::<kilometer>(),
                1e-12,
            );
        }
    }

    #[test]
    fn observer_trajectory_returns_correct_metadata() {
        let propagator = test_propagator();
        let range = test_range();

        let observer = test_observer();

        let compute = LookAnglesComputation::default();

        let result = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        assert_eq!(result.start, range.start);
        assert_eq!(result.end, range.end);
        assert_eq!(result.step_seconds, 60);
    }

    #[test]
    fn observer_trajectory_contains_expected_number_of_samples() {
        let propagator = test_propagator();
        let range = test_range();

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let result = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 11);
    }

    #[test]
    fn observer_trajectory_contains_single_sample_for_zero_length_range() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let range = TimeRange {
            start: datetime,
            end: datetime,
        };

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let result = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 1);
        assert_eq!(result.samples[0].time, datetime);
    }

    #[test]
    fn observer_trajectory_samples_are_at_correct_times() {
        let propagator = test_propagator();
        let start = test_datetime();

        let range = TimeRange {
            start,
            end: start + Duration::minutes(5),
        };

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let result = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        assert_eq!(result.samples.len(), 6);

        for (i, sample) in result.samples.iter().enumerate() {
            let expected = start + Duration::seconds(i as i64 * 60);

            assert_eq!(sample.time, expected);
        }
    }

    #[test]
    fn observer_trajectory_samples_are_ordered() {
        let propagator = test_propagator();
        let range = test_range();

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let result = propagator
            .observer_trajectory_at(range, sampling(30), &observer, &compute)
            .unwrap();

        assert!(
            result
                .samples
                .windows(2)
                .all(|window| window[0].time < window[1].time)
        );
    }

    #[test]
    fn observer_trajectory_respects_sampling_interval() {
        let propagator = test_propagator();
        let range = test_range();

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let result = propagator
            .observer_trajectory_at(range, sampling(90), &observer, &compute)
            .unwrap();

        for window in result.samples.windows(2) {
            let delta = window[1].time - window[0].time;

            assert_eq!(delta, Duration::seconds(90));
        }
    }

    #[test]
    fn observer_trajectory_matches_individual_look_angle_calls() {
        let propagator = test_propagator();
        let start = test_datetime();

        let range = TimeRange {
            start,
            end: start + Duration::minutes(2),
        };

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let trajectory = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        assert_eq!(trajectory.samples.len(), 3);

        for (i, sample) in trajectory.samples.iter().enumerate() {
            let datetime = start + Duration::seconds(i as i64 * 60);

            let expected = propagator
                .look_angles_at(datetime, &observer, &compute)
                .unwrap();

            assert_eq!(sample.time, expected.time);
        }
    }

    #[test]
    fn observer_trajectory_is_deterministic() {
        let propagator = test_propagator();
        let range = test_range();

        let observer = test_observer();
        let compute = LookAnglesComputation::default();

        let first = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        let second = propagator
            .observer_trajectory_at(range, sampling(60), &observer, &compute)
            .unwrap();

        assert_eq!(first.samples.len(), second.samples.len());

        for (a, b) in first.samples.iter().zip(second.samples.iter()) {
            assert_eq!(a.time, b.time);
        }
    }

    #[test]
    fn parallel_trajectory_matches_sequential_trajectory() {
        let propagator = test_propagator();
        let start = test_datetime();

        let range = TimeRange {
            start,
            end: start + Duration::minutes(30),
        };

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let parallel = propagator
            .trajectory_at(range, sampling(10), &compute)
            .unwrap();

        let step = i64::from(10_u32);
        let total_secs = (range.end - range.start).num_seconds();
        let steps = (total_secs / step) + 1;

        let sequential = (0..steps)
            .map(|i| {
                let t = range.start + Duration::seconds(i * step);
                propagator.position_at(t, &compute).unwrap()
            })
            .collect::<Vec<_>>();

        assert_eq!(parallel.samples.len(), sequential.len());

        for (parallel, sequential) in parallel.samples.iter().zip(sequential.iter()) {
            assert_eq!(parallel.time, sequential.time);

            let parallel_teme = parallel.teme.as_ref().unwrap();
            let sequential_teme = sequential.teme.as_ref().unwrap();

            assert_close(
                parallel_teme.x.get::<kilometer>(),
                sequential_teme.x.get::<kilometer>(),
                1e-12,
            );

            assert_close(
                parallel_teme.y.get::<kilometer>(),
                sequential_teme.y.get::<kilometer>(),
                1e-12,
            );

            assert_close(
                parallel_teme.z.get::<kilometer>(),
                sequential_teme.z.get::<kilometer>(),
                1e-12,
            );
        }
    }
}
