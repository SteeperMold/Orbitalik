use async_trait::async_trait;
use chrono::{DateTime, Utc};
use std::sync::Arc;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{LookAngles, SatelliteIdentifier, SatellitePosition};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::astro::propagation::propagator::Propagator;
use crate::domain::errors::ServiceError;
use crate::domain::models::TrajectoryComputationMetadata;
use crate::domain::position::PositionServiceApi;
use crate::domain::tle_provider::TleProvider;

pub struct PositionService {
    tle_grpc_client: Arc<dyn TleProvider>,
}

impl PositionService {
    pub const fn new(tle_grpc_client: Arc<dyn TleProvider>) -> Self {
        Self { tle_grpc_client }
    }
}

#[async_trait]
impl PositionServiceApi for PositionService {
    async fn get_position(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        compute: &PositionComputation,
    ) -> Result<(SatellitePosition, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let position = Propagator::from_tle(&tle)?.position_at(datetime, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((position, metadata))
    }

    async fn get_look_angles(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(LookAngles, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let look_angles =
            Propagator::from_tle(&tle)?.look_angles_at(datetime, observer, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((look_angles, metadata))
    }
}

#[allow(clippy::unwrap_used)]
#[cfg(test)]
mod tests {
    use super::*;

    use chrono::Utc;
    use std::sync::Arc;
    use uom::si::angle::degree;
    use uom::si::length::kilometer;

    use crate::astro::coords::geodetic::Geodetic;
    use crate::astro::models::SatelliteIdentifier;
    use crate::astro::propagation::look_angles::LookAnglesComputation;
    use crate::astro::propagation::position::PositionComputation;
    use crate::astro::test_utils::test_tle;
    use crate::domain::tle_provider::MockTleProvider;

    fn test_satellite_identifier() -> SatelliteIdentifier {
        SatelliteIdentifier::NoradId(25544)
    }

    fn test_observer() -> Geodetic {
        Geodetic {
            lat: uom::si::f64::Angle::new::<degree>(0.0),
            lon: uom::si::f64::Angle::new::<degree>(0.0),
            alt: uom::si::f64::Length::new::<kilometer>(0.0),
        }
    }

    fn service(mock: MockTleProvider) -> PositionService {
        PositionService::new(Arc::new(mock))
    }

    #[tokio::test]
    async fn get_position_returns_position_and_metadata() {
        let tle = test_tle();
        let satellite = test_satellite_identifier();
        let datetime = tle.epoch;

        let mut mock = MockTleProvider::new();

        let expected_satellite = satellite.clone();

        mock.expect_get_tle()
            .withf(move |identifier| {
                matches!(
                    (identifier, &expected_satellite),
                    (
                        SatelliteIdentifier::NoradId(actual),
                        SatelliteIdentifier::NoradId(expected)
                    ) if actual == expected
                )
            })
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(mock);

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let result = service.get_position(satellite, datetime, &compute).await;

        assert!(result.is_ok());

        let (position, metadata) = result.unwrap();

        assert_eq!(position.time, datetime);

        assert!(position.teme.is_some());
        assert!(position.ecef.is_some());
        assert!(position.geodetic.is_some());

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_id, 25544);
        assert!(!metadata.satellite_name.is_empty());
        assert_eq!(metadata.tle_epoch, datetime);
    }

    #[tokio::test]
    async fn get_position_requests_tle_using_given_satellite() {
        let tle = test_tle();
        let expected_tle = test_tle();

        let satellite = test_satellite_identifier();

        let mut mock = MockTleProvider::new();

        let expected_satellite = satellite.clone();

        mock.expect_get_tle()
            .withf(move |identifier| {
                matches!(
                    (identifier, &expected_satellite),
                    (
                        SatelliteIdentifier::NoradId(actual),
                        SatelliteIdentifier::NoradId(expected)
                    ) if actual == expected
                )
            })
            .times(1)
            .return_once(move |_| Ok(expected_tle));

        let service = service(mock);

        let compute = PositionComputation::default();

        let _ = service
            .get_position(satellite, tle.epoch, &compute)
            .await
            .unwrap();
    }

    #[tokio::test]
    async fn get_position_respects_computation_options() {
        let tle = test_tle();
        let datetime = tle.epoch;

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let (position, _) = service
            .get_position(test_satellite_identifier(), datetime, &compute)
            .await
            .unwrap();

        assert!(position.teme.is_some());
        assert!(position.ecef.is_none());
        assert!(position.geodetic.is_none());
    }

    #[tokio::test]
    async fn get_position_with_no_computation_returns_empty_position() {
        let tle = test_tle();
        let datetime = tle.epoch;

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let compute = PositionComputation::default();

        let (position, _) = service
            .get_position(test_satellite_identifier(), datetime, &compute)
            .await
            .unwrap();

        assert!(position.teme.is_none());
        assert!(position.ecef.is_none());
        assert!(position.geodetic.is_none());
    }

    #[tokio::test]
    async fn get_position_propagation_errors_are_returned() {
        let tle = test_tle();
        let datetime = tle.epoch;

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = service
            .get_position(test_satellite_identifier(), datetime, &compute)
            .await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn get_look_angles_returns_look_angles_and_metadata() {
        let tle = test_tle();
        let datetime = tle.epoch;
        let observer = test_observer();

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: true,
            range: true,
        };

        let result = service
            .get_look_angles(test_satellite_identifier(), datetime, &observer, &compute)
            .await;

        assert!(result.is_ok());

        let (look_angles, metadata) = result.unwrap();

        assert_eq!(look_angles.time, datetime);

        assert!(look_angles.azimuth.is_some());
        assert!(look_angles.elevation.is_some());
        assert!(look_angles.range.is_some());

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_id, 25544);
        assert_eq!(metadata.tle_epoch, datetime);
    }

    #[tokio::test]
    async fn get_look_angles_respects_computation_options() {
        let tle = test_tle();
        let datetime = tle.epoch;
        let observer = test_observer();

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: false,
            range: false,
        };

        let (result, _) = service
            .get_look_angles(test_satellite_identifier(), datetime, &observer, &compute)
            .await
            .unwrap();

        assert!(result.azimuth.is_some());
        assert!(result.elevation.is_none());
        assert!(result.range.is_none());
    }

    #[tokio::test]
    async fn metadata_contains_tle_information() {
        let tle = test_tle();
        let expected_name = tle.satellite_name.clone();
        let expected_norad = tle.norad_id;
        let expected_epoch = tle.epoch;
        let datetime = tle.epoch;

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let (position, metadata) = service
            .get_position(
                test_satellite_identifier(),
                datetime,
                &PositionComputation::default(),
            )
            .await
            .unwrap();

        assert_eq!(position.time, datetime);

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_id, expected_norad);
        assert_eq!(metadata.satellite_name, expected_name);
        assert_eq!(metadata.tle_epoch, expected_epoch);
    }

    #[tokio::test]
    async fn computation_time_is_populated() {
        let tle = test_tle();
        let datetime = tle.epoch;

        let before = Utc::now();

        let mut mock = MockTleProvider::new();

        mock.expect_get_tle().times(1).return_once(move |_| Ok(tle));

        let service = service(mock);

        let (_, metadata) = service
            .get_position(
                test_satellite_identifier(),
                datetime,
                &PositionComputation::default(),
            )
            .await
            .unwrap();

        let after = Utc::now();

        assert!(metadata.computation_time >= before);
        assert!(metadata.computation_time <= after);
    }
}
