use rayon::prelude::*;
use std::sync::atomic::AtomicUsize;
use uom::si::f64::Angle;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::errors::PropagationError;
use crate::astro::models::{Pass, TimeRange};
use crate::astro::propagation::propagator::Propagator;

pub struct PassPredictor;

pub struct PassPredictionOptions<'a> {
    pub range: TimeRange,
    pub observer: &'a Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub max_results: Option<usize>,
}

impl PassPredictor {
    pub fn predict_many(
        satellites: &[Propagator],
        options: &PassPredictionOptions,
    ) -> Result<Vec<Pass>, PropagationError> {
        let remaining = AtomicUsize::new(options.max_results.unwrap_or(usize::MAX));

        let mut all: Vec<Pass> = satellites
            .par_iter()
            .map(|propagator| propagator.predict_passes(options, &remaining))
            .collect::<Result<Vec<_>, _>>()?
            .into_par_iter()
            .flatten()
            .collect();

        all.sort_by_key(|p| p.aos);

        Ok(all)
    }
}

#[allow(clippy::unwrap_used, clippy::panic)]
#[cfg(test)]
mod tests {
    use super::*;

    use chrono::Duration;
    use uom::si::angle::degree;

    use crate::astro::models::SatelliteIdentifier;
    use crate::astro::test_utils::{test_datetime, test_observer, test_propagator};

    fn test_options(
        observer: &'_ Geodetic,
        max_results: Option<usize>,
    ) -> PassPredictionOptions<'_> {
        let start = test_datetime();

        PassPredictionOptions {
            range: TimeRange {
                start,
                end: start + Duration::hours(6),
            },
            observer,
            min_elevation: Angle::new::<degree>(0.0),
            min_peak_elevation: Angle::new::<degree>(0.0),
            max_results,
        }
    }

    #[test]
    fn no_satellites_returns_empty_result() {
        let observer = test_observer();
        let options = test_options(&observer, None);

        let result = PassPredictor::predict_many(&[], &options).unwrap();

        assert!(result.is_empty());
    }

    #[test]
    fn single_satellite_returns_its_passes() {
        let propagator = test_propagator();
        let observer = test_observer();

        let options = test_options(&observer, None);

        let result = PassPredictor::predict_many(&[propagator], &options).unwrap();

        assert!(!result.is_empty());

        for pass in &result {
            assert!(pass.aos <= pass.los);
            assert!(pass.max_elevation.get::<degree>() >= 0.0);
        }
    }

    #[test]
    fn max_results_zero_returns_empty_result() {
        let propagator = test_propagator();
        let observer = test_observer();

        let options = test_options(&observer, Some(0));

        let result = PassPredictor::predict_many(&[propagator], &options).unwrap();

        assert!(result.is_empty());
    }

    #[test]
    fn max_results_limits_number_of_passes() {
        let propagator = test_propagator();
        let observer = test_observer();

        let options = test_options(&observer, Some(2));

        let result = PassPredictor::predict_many(&[propagator], &options).unwrap();

        assert!(result.len() <= 2);
    }

    #[test]
    fn max_results_one_returns_at_most_one_pass() {
        let propagator = test_propagator();
        let observer = test_observer();

        let options = test_options(&observer, Some(1));

        let result = PassPredictor::predict_many(&[propagator], &options).unwrap();

        assert!(result.len() <= 1);
    }

    #[test]
    fn no_max_results_returns_all_available_passes() {
        let propagator = test_propagator();
        let observer = test_observer();

        let unlimited = test_options(&observer, None);
        let limited = test_options(&observer, Some(usize::MAX));

        let all = PassPredictor::predict_many(&[propagator], &unlimited).unwrap();

        let max_usize = PassPredictor::predict_many(&[test_propagator()], &limited).unwrap();

        assert_eq!(all.len(), max_usize.len());
    }

    #[test]
    fn results_are_sorted_by_aos() {
        let satellites = vec![test_propagator(), test_propagator(), test_propagator()];

        let observer = test_observer();
        let options = test_options(&observer, None);

        let result = PassPredictor::predict_many(&satellites, &options).unwrap();

        assert!(
            result
                .windows(2)
                .all(|window| window[0].aos <= window[1].aos)
        );
    }

    #[test]
    fn results_are_deterministic() {
        let satellites = vec![test_propagator(), test_propagator(), test_propagator()];

        let observer = test_observer();
        let options = test_options(&observer, None);

        let first = PassPredictor::predict_many(&satellites, &options).unwrap();

        let second = PassPredictor::predict_many(&satellites, &options).unwrap();

        assert_eq!(first.len(), second.len());

        for (a, b) in first.iter().zip(second.iter()) {
            assert_eq!(a.aos, b.aos);
            assert_eq!(a.los, b.los);
            assert_eq!(a.max_elevation_time, b.max_elevation_time);

            match (&a.satellite, &b.satellite) {
                (SatelliteIdentifier::NoradId(a_id), SatelliteIdentifier::NoradId(b_id)) => {
                    assert_eq!(a_id, b_id);
                }

                _ => panic!("unexpected satellite identifier"),
            }

            assert_eq!(a.duration_seconds, b.duration_seconds);
        }
    }

    #[test]
    fn multiple_satellites_are_combined_into_one_result() {
        let satellites = vec![test_propagator(), test_propagator()];

        let observer = test_observer();
        let options = test_options(&observer, None);

        let result = PassPredictor::predict_many(&satellites, &options).unwrap();

        assert!(!result.is_empty());

        assert!(
            result
                .windows(2)
                .all(|window| window[0].aos <= window[1].aos)
        );
    }

    #[test]
    fn min_peak_elevation_filters_passes() {
        let propagator = test_propagator();
        let observer = test_observer();

        let all_options = test_options(&observer, None);

        let all = PassPredictor::predict_many(&[propagator], &all_options).unwrap();

        let high_options = PassPredictionOptions {
            range: all_options.range,
            observer: all_options.observer,
            min_elevation: all_options.min_elevation,
            min_peak_elevation: Angle::new::<degree>(80.0),
            max_results: None,
        };

        let high = PassPredictor::predict_many(&[test_propagator()], &high_options).unwrap();

        assert!(high.len() <= all.len());

        for pass in high {
            assert!(pass.max_elevation.get::<degree>() >= 80.0);
        }
    }

    #[test]
    fn higher_min_elevation_cannot_create_more_passes() {
        let propagator = test_propagator();
        let observer = test_observer();

        let low_options = PassPredictionOptions {
            range: test_options(&observer, None).range,
            observer: &observer,
            min_elevation: Angle::new::<degree>(0.0),
            min_peak_elevation: Angle::new::<degree>(0.0),
            max_results: None,
        };

        let high_options = PassPredictionOptions {
            range: low_options.range,
            observer: &observer,
            min_elevation: Angle::new::<degree>(20.0),
            min_peak_elevation: Angle::new::<degree>(20.0),
            max_results: None,
        };

        let low = PassPredictor::predict_many(&[propagator], &low_options).unwrap();

        let high = PassPredictor::predict_many(&[test_propagator()], &high_options).unwrap();

        assert!(high.len() <= low.len());
    }

    #[test]
    fn results_respect_max_results_across_multiple_satellites() {
        let satellites = vec![test_propagator(), test_propagator(), test_propagator()];

        let observer = test_observer();

        let options = test_options(&observer, Some(3));

        let result = PassPredictor::predict_many(&satellites, &options).unwrap();

        assert!(result.len() <= 3);
    }

    #[test]
    fn max_results_is_global_not_per_satellite() {
        let satellites = vec![test_propagator(), test_propagator(), test_propagator()];

        let observer = test_observer();

        let options = test_options(&observer, Some(1));

        let result = PassPredictor::predict_many(&satellites, &options).unwrap();

        assert!(result.len() <= 1);
    }
}
