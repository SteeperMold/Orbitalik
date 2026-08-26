use async_trait::async_trait;
use chrono::Utc;
use std::sync::Arc;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::ObserverTrajectory;
use crate::astro::models::Sampling;
use crate::astro::models::SatelliteIdentifier;
use crate::astro::models::TimeRange;
use crate::astro::models::Trajectory;
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::astro::propagation::propagator::Propagator;
use crate::domain::errors::ServiceError;
use crate::domain::models::TrajectoryComputationMetadata;
use crate::domain::tle_provider::TleProvider;
use crate::domain::trajectory::TrajectoryServiceApi;

pub struct TrajectoryService {
    tle_grpc_client: Arc<dyn TleProvider>,
}

impl TrajectoryService {
    pub const fn new(tle_grpc_client: Arc<dyn TleProvider>) -> Self {
        Self { tle_grpc_client }
    }
}

#[async_trait]
impl TrajectoryServiceApi for TrajectoryService {
    async fn get_trajectory(
        &self,
        satellite_identifier: SatelliteIdentifier,
        range: TimeRange,
        sampling: Sampling,
        compute: &PositionComputation,
    ) -> Result<(Trajectory, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let trajectory = Propagator::from_tle(&tle)?.trajectory_at(range, sampling, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((trajectory, metadata))
    }

    async fn get_observer_trajectory(
        &self,
        satellite_identifier: SatelliteIdentifier,
        range: TimeRange,
        sampling: Sampling,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(ObserverTrajectory, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let observer_trajectory = Propagator::from_tle(&tle)?
            .observer_trajectory_at(range, sampling, observer, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((observer_trajectory, metadata))
    }
}

#[allow(
    clippy::unwrap_used,
    clippy::cast_sign_loss,
    clippy::cast_possible_truncation
)]
#[cfg(test)]
mod tests {
    use super::*;

    use chrono::Duration;
    use mockall::predicate::function;
    use uom::si::angle::{Angle, degree};
    use uom::si::length::kilometer;
    
    use crate::astro::test_utils::{test_datetime, test_tle};
    use crate::domain::tle_provider::MockTleProvider;

    fn test_satellite() -> SatelliteIdentifier {
        SatelliteIdentifier::NoradId(25544)
    }

    fn test_range() -> TimeRange {
        let start = test_datetime();

        TimeRange {
            start,
            end: start + Duration::minutes(10),
        }
    }

    fn test_sampling() -> Sampling {
        Sampling { step_seconds: 60 }
    }

    fn test_observer() -> Geodetic {
        Geodetic {
            lat: Angle::new::<degree>(0.0),
            lon: Angle::new::<degree>(0.0),
            alt: uom::si::f64::Length::new::<kilometer>(0.0),
        }
    }

    fn test_provider(tle: crate::astro::models::Tle) -> MockTleProvider {
        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        provider
    }

    #[tokio::test]
    async fn get_trajectory_returns_trajectory_and_metadata() {
        let tle = test_tle();
        let expected_norad_id = tle.norad_id;
        let expected_name = tle.satellite_name.clone();
        let expected_epoch = tle.epoch;

        let range = test_range();
        let sampling = test_sampling();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let (trajectory, metadata) = service
            .get_trajectory(test_satellite(), range, sampling, &compute)
            .await
            .unwrap();

        assert_eq!(trajectory.start, range.start);
        assert_eq!(trajectory.end, range.end);
        assert_eq!(trajectory.step_seconds, sampling.step_seconds);

        assert!(!trajectory.samples.is_empty());

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_id, expected_norad_id);
        assert_eq!(metadata.satellite_name, expected_name);
        assert_eq!(metadata.tle_epoch, expected_epoch);
    }

    #[tokio::test]
    async fn get_trajectory_requests_correct_satellite() {
        let tle = test_tle();
        let expected_satellite = test_satellite();

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .with(function({
                let expected = expected_satellite.clone();

                move |actual| {
                    // Adapt this to your SatelliteIdentifier definition
                    matches!(
                        (actual, &expected),
                        (
                            SatelliteIdentifier::NoradId(actual_id),
                            SatelliteIdentifier::NoradId(expected_id)
                        ) if actual_id == expected_id
                    )
                }
            }))
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = TrajectoryService::new(Arc::new(provider));

        service
            .get_trajectory(
                expected_satellite,
                test_range(),
                test_sampling(),
                &PositionComputation::default(),
            )
            .await
            .unwrap();
    }

    #[tokio::test]
    async fn get_trajectory_respects_computation_options() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let (trajectory, _) = service
            .get_trajectory(test_satellite(), test_range(), test_sampling(), &compute)
            .await
            .unwrap();

        assert!(!trajectory.samples.is_empty());

        for sample in trajectory.samples {
            assert!(sample.teme.is_some());
            assert!(sample.ecef.is_none());
            assert!(sample.geodetic.is_none());
        }
    }

    #[tokio::test]
    async fn get_trajectory_with_no_computation_returns_empty_samples() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let compute = PositionComputation::default();

        let (trajectory, _) = service
            .get_trajectory(test_satellite(), test_range(), test_sampling(), &compute)
            .await
            .unwrap();

        assert!(!trajectory.samples.is_empty());

        for sample in trajectory.samples {
            assert!(sample.teme.is_none());
            assert!(sample.ecef.is_none());
            assert!(sample.geodetic.is_none());
        }
    }

    #[tokio::test]
    async fn get_trajectory_preserves_range_and_sampling() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let range = test_range();

        let sampling = Sampling { step_seconds: 30 };

        let (trajectory, _) = service
            .get_trajectory(
                test_satellite(),
                range,
                sampling,
                &PositionComputation::default(),
            )
            .await
            .unwrap();

        assert_eq!(trajectory.start, range.start);
        assert_eq!(trajectory.end, range.end);
        assert_eq!(trajectory.step_seconds, 30);
    }

    #[tokio::test]
    async fn get_trajectory_sample_count_matches_sampling() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let range = test_range();

        let sampling = Sampling { step_seconds: 60 };

        let (trajectory, _) = service
            .get_trajectory(
                test_satellite(),
                range,
                sampling,
                &PositionComputation::default(),
            )
            .await
            .unwrap();

        let expected_count = ((range.end - range.start).num_seconds() / 60) + 1;

        assert_eq!(trajectory.samples.len(), expected_count as usize);
    }

    #[tokio::test]
    async fn get_observer_trajectory_returns_trajectory_and_metadata() {
        let tle = test_tle();

        let expected_norad_id = tle.norad_id;
        let expected_name = tle.satellite_name.clone();
        let expected_epoch = tle.epoch;

        let range = test_range();
        let sampling = test_sampling();
        let observer = test_observer();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: true,
            range: true,
        };

        let (trajectory, metadata) = service
            .get_observer_trajectory(test_satellite(), range, sampling, &observer, &compute)
            .await
            .unwrap();

        assert_eq!(trajectory.start, range.start);
        assert_eq!(trajectory.end, range.end);
        assert_eq!(trajectory.step_seconds, sampling.step_seconds);

        assert!(!trajectory.samples.is_empty());

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_id, expected_norad_id);
        assert_eq!(metadata.satellite_name, expected_name);
        assert_eq!(metadata.tle_epoch, expected_epoch);
    }

    #[tokio::test]
    async fn get_observer_trajectory_respects_computation_options() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: false,
            range: false,
        };

        let (trajectory, _) = service
            .get_observer_trajectory(
                test_satellite(),
                test_range(),
                test_sampling(),
                &test_observer(),
                &compute,
            )
            .await
            .unwrap();

        assert!(!trajectory.samples.is_empty());

        for sample in trajectory.samples {
            assert!(sample.azimuth.is_some());
            assert!(sample.elevation.is_none());
            assert!(sample.range.is_none());
        }
    }

    #[tokio::test]
    async fn get_observer_trajectory_preserves_range_and_sampling() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let range = test_range();

        let sampling = Sampling { step_seconds: 30 };

        let (trajectory, _) = service
            .get_observer_trajectory(
                test_satellite(),
                range,
                sampling,
                &test_observer(),
                &LookAnglesComputation::default(),
            )
            .await
            .unwrap();

        assert_eq!(trajectory.start, range.start);
        assert_eq!(trajectory.end, range.end);
        assert_eq!(trajectory.step_seconds, 30);
    }

    #[tokio::test]
    async fn get_observer_trajectory_sample_count_matches_sampling() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let range = test_range();

        let sampling = Sampling { step_seconds: 30 };

        let (trajectory, _) = service
            .get_observer_trajectory(
                test_satellite(),
                range,
                sampling,
                &test_observer(),
                &LookAnglesComputation::default(),
            )
            .await
            .unwrap();

        let expected_count = ((range.end - range.start).num_seconds() / 30) + 1;

        assert_eq!(trajectory.samples.len(), expected_count as usize);
    }

    #[tokio::test]
    async fn get_trajectory_propagates_tle_error() {
        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(|_| Err(tonic::Status::not_found("TLE not found")));

        let service = TrajectoryService::new(Arc::new(provider));

        let result = service
            .get_trajectory(
                test_satellite(),
                test_range(),
                test_sampling(),
                &PositionComputation::default(),
            )
            .await;

        assert!(result.is_err());
    }

    #[tokio::test]
    async fn get_observer_trajectory_propagates_tle_error() {
        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(|_| Err(tonic::Status::not_found("TLE not found")));

        let service = TrajectoryService::new(Arc::new(provider));

        let result = service
            .get_observer_trajectory(
                test_satellite(),
                test_range(),
                test_sampling(),
                &test_observer(),
                &LookAnglesComputation::default(),
            )
            .await;

        assert!(result.is_err());
    }

    #[tokio::test]
    async fn metadata_computation_time_is_populated() {
        let tle = test_tle();

        let provider = test_provider(tle);
        let service = TrajectoryService::new(Arc::new(provider));

        let before = Utc::now();

        let (_, metadata) = service
            .get_trajectory(
                test_satellite(),
                test_range(),
                test_sampling(),
                &PositionComputation::default(),
            )
            .await
            .unwrap();

        let after = Utc::now();

        assert!(metadata.computation_time >= before);
        assert!(metadata.computation_time <= after);
    }
}
