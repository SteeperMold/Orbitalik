use tonic::{Request, Response, Status};

use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::position::PositionServiceApi;
use crate::transport::grpc::converters::ToChrono;
use crate::transport::grpc::server::TrajectoryGrpcServer;
use crate::transport::grpc::server::trajectory_grpc::{
    LookAnglesRequest, LookAnglesResponse, PositionRequest, PositionResponse,
};

impl<P, T, Pa> TrajectoryGrpcServer<P, T, Pa>
where
    P: PositionServiceApi,
    T: Sync,
    Pa: Sync,
{
    pub async fn handle_get_position(
        &self,
        request: Request<PositionRequest>,
    ) -> Result<Response<PositionResponse>, Status> {
        let req = request.into_inner();

        let identifier = req
            .identifier
            .ok_or_else(|| Status::invalid_argument("Missing satellite identifier"))?
            .try_into()?;

        let datetime = req
            .datetime
            .ok_or_else(|| Status::invalid_argument("Missing datetime"))?
            .to_chrono()
            .map_err(|_| Status::invalid_argument("Invalid datetime"))?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(PositionComputation::default, PositionComputation::from);

        let (position, metadata) = self
            .position
            .get_position(identifier, datetime, &compute)
            .await?;

        let response = PositionResponse::from_position(&position, metadata, req.units)?;
        Ok(Response::new(response))
    }

    pub async fn handle_look_angles(
        &self,
        request: Request<LookAnglesRequest>,
    ) -> Result<Response<LookAnglesResponse>, Status> {
        let req = request.into_inner();

        let identifier = req
            .identifier
            .ok_or_else(|| {
                Status::invalid_argument("Missing satellite identifier (norad_id or name)")
            })?
            .try_into()?;

        let datetime = req
            .datetime
            .ok_or_else(|| Status::invalid_argument("Missing datetime"))?
            .to_chrono()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(LookAnglesComputation::default, LookAnglesComputation::from);

        let (look_angles, metadata) = self
            .position
            .get_look_angles(identifier, datetime, &observer, &compute)
            .await?;

        let response = LookAnglesResponse::from_look_angles(&look_angles, metadata, req.units)?;
        Ok(Response::new(response))
    }
}

#[allow(clippy::unwrap_used)]
#[cfg(test)]
mod tests {
    use chrono::{TimeZone, Utc};
    use prost_types::{FieldMask, Timestamp};
    use tonic::{Code, Request, Status};

    use crate::astro::models::{
        LookAngles, SatelliteIdentifier as DomainSatelliteIdentifier, SatellitePosition,
    };
    use crate::astro::propagation::position::PositionComputation;
    use crate::domain::models::TrajectoryComputationMetadata;
    use crate::domain::passes::MockPassesServiceApi;
    use crate::domain::position::MockPositionServiceApi;
    use crate::domain::trajectory::MockTrajectoryServiceApi;
    use crate::transport::grpc::server::TrajectoryGrpcServer;
    use crate::transport::grpc::server::trajectory_grpc::geodetic_input::{Alt, Lat, Lon};
    use crate::transport::grpc::server::trajectory_grpc::{
        GeodeticInput, LookAnglesRequest, PositionRequest, UnitSettings,
        unit_settings::{AngleUnit, DistanceUnit},
    };
    use crate::transport::grpc::server::trajectory_grpc::{
        SatelliteIdentifier, satellite_identifier,
    };

    fn norad_identifier(id: u32) -> SatelliteIdentifier {
        SatelliteIdentifier {
            kind: Some(satellite_identifier::Kind::NoradId(id)),
        }
    }

    fn valid_timestamp() -> Timestamp {
        Timestamp {
            seconds: 1_700_000_000,
            nanos: 0,
        }
    }

    fn valid_units() -> UnitSettings {
        UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Degrees as i32,
        }
    }

    fn valid_metadata() -> TrajectoryComputationMetadata {
        TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc.timestamp_opt(1_700_000_000, 0).unwrap(),
            norad_id: 25544,
            satellite_name: "ISS (ZARYA)".to_string(),
            tle_epoch: Utc.timestamp_opt(1_700_000_000, 0).unwrap(),
        }
    }

    fn valid_position() -> SatellitePosition {
        SatellitePosition {
            time: Utc.timestamp_opt(1_700_000_000, 0).unwrap(),
            teme: None,
            ecef: None,
            geodetic: None,
        }
    }

    fn valid_observer() -> GeodeticInput {
        GeodeticInput {
            lat: Some(Lat::LatDeg(56.9496)),
            lon: Some(Lon::LonDeg(24.1052)),
            alt: Some(Alt::AltM(0.0)),
        }
    }

    fn valid_look_angles() -> LookAngles {
        LookAngles {
            time: Utc.timestamp_opt(1_700_000_000, 0).unwrap(),
            azimuth: None,
            elevation: None,
            range: None,
        }
    }

    fn test_server(
        position: MockPositionServiceApi,
    ) -> TrajectoryGrpcServer<MockPositionServiceApi, MockTrajectoryServiceApi, MockPassesServiceApi>
    {
        TrajectoryGrpcServer::new(
            position,
            MockTrajectoryServiceApi::new(),
            MockPassesServiceApi::new(),
        )
    }

    #[tokio::test]
    async fn get_position_rejects_missing_identifier() {
        let mut position = MockPositionServiceApi::new();

        position.expect_get_position().never();

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: None,
            datetime: Some(valid_timestamp()),
            output_mask: None,
            units: None,
        });

        let error = server.handle_get_position(request).await.unwrap_err();

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing satellite identifier");
    }

    #[tokio::test]
    async fn get_position_rejects_missing_datetime() {
        let mut position = MockPositionServiceApi::new();

        position.expect_get_position().never();

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: None,
            output_mask: None,
            units: None,
        });

        let error = server.handle_get_position(request).await.unwrap_err();

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing datetime");
    }

    #[tokio::test]
    async fn get_position_rejects_invalid_timestamp() {
        let mut position = MockPositionServiceApi::new();

        position.expect_get_position().never();

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(Timestamp {
                seconds: 0,
                nanos: -1,
            }),
            output_mask: None,
            units: None,
        });

        let error = server.handle_get_position(request).await.unwrap_err();

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_position_uses_default_computation_without_output_mask() {
        let mut position = MockPositionServiceApi::new();

        position
            .expect_get_position()
            .withf(|identifier, datetime, compute| {
                *identifier == DomainSatelliteIdentifier::NoradId(25544)
                    && *datetime == Utc.timestamp_opt(1_700_000_000, 0).unwrap()
                    && *compute == PositionComputation::default()
            })
            .times(1)
            .returning(|_, _, _| Ok((valid_position(), valid_metadata())));

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            output_mask: None,
            units: Some(valid_units()),
        });

        let response = server
            .handle_get_position(request)
            .await
            .expect("get_position failed")
            .into_inner();

        assert!(response.time.is_some());
        assert!(response.metadata.is_some());
    }

    #[tokio::test]
    async fn get_position_uses_output_mask() {
        let mut position = MockPositionServiceApi::new();

        position
            .expect_get_position()
            .withf(|identifier, datetime, compute| {
                *identifier == DomainSatelliteIdentifier::NoradId(25544)
                    && *datetime == Utc.timestamp_opt(1_700_000_000, 0).unwrap()
                    && compute.teme == false
                    && compute.ecef
                    && compute.geodetic == false
            })
            .times(1)
            .returning(|_, _, _| Ok((valid_position(), valid_metadata())));

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            output_mask: Some(FieldMask {
                paths: vec!["ecef".to_string()],
            }),
            units: Some(valid_units()),
        });

        let response = server
            .handle_get_position(request)
            .await
            .expect("get_position failed")
            .into_inner();

        assert!(response.time.is_some());
    }

    #[tokio::test]
    async fn get_position_uses_all_requested_output_fields() {
        let mut position = MockPositionServiceApi::new();

        position
            .expect_get_position()
            .withf(|_, _, compute| compute.teme && compute.ecef && compute.geodetic)
            .times(1)
            .returning(|_, _, _| Ok((valid_position(), valid_metadata())));

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            output_mask: Some(FieldMask {
                paths: vec![
                    "eci".to_string(),
                    "ecef".to_string(),
                    "geodetic".to_string(),
                ],
            }),
            units: Some(valid_units()),
        });

        server
            .handle_get_position(request)
            .await
            .expect("get_position failed");
    }

    #[tokio::test]
    async fn get_position_propagates_service_error() {
        let mut position = MockPositionServiceApi::new();

        position
            .expect_get_position()
            .times(1)
            .returning(|_, _, _| Err(Status::internal("position calculation failed").into()));

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            output_mask: None,
            units: Some(valid_units()),
        });

        let error = server
            .handle_get_position(request)
            .await
            .expect_err("expected get_position to fail");

        assert_eq!(error.code(), Code::Internal);
        assert_eq!(error.message(), "position calculation failed");
    }

    #[tokio::test]
    async fn get_position_rejects_missing_units() {
        let mut position = MockPositionServiceApi::new();

        position
            .expect_get_position()
            .times(1)
            .returning(|_, _, _| Ok((valid_position(), valid_metadata())));

        let server = test_server(position);

        let request = Request::new(PositionRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            output_mask: None,
            units: None,
        });

        let error = server
            .handle_get_position(request)
            .await
            .expect_err("expected missing units to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }
    #[tokio::test]
    async fn get_look_angles_rejects_missing_identifier() {
        let mut position = MockPositionServiceApi::new();

        position.expect_get_look_angles().never();

        let server = test_server(position);

        let request = Request::new(LookAnglesRequest {
            identifier: None,
            datetime: Some(valid_timestamp()),
            observer: None,
            output_mask: None,
            units: None,
        });

        let error = server
            .handle_look_angles(request)
            .await
            .expect_err("expected request to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(
            error.message(),
            "Missing satellite identifier (norad_id or name)"
        );
    }
    #[tokio::test]
    async fn get_look_angles_rejects_missing_datetime() {
        let mut position = MockPositionServiceApi::new();

        position.expect_get_look_angles().never();

        let server = test_server(position);

        let request = Request::new(LookAnglesRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: None,
            observer: None,
            output_mask: None,
            units: None,
        });

        let error = server
            .handle_look_angles(request)
            .await
            .expect_err("expected request to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing datetime");
    }
    #[tokio::test]
    async fn get_look_angles_rejects_missing_observer() {
        let mut position = MockPositionServiceApi::new();

        position.expect_get_look_angles().never();

        let server = test_server(position);

        let request = Request::new(LookAnglesRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            observer: None,
            output_mask: None,
            units: None,
        });

        let error = server
            .handle_look_angles(request)
            .await
            .expect_err("expected request to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing observer");
    }
    #[tokio::test]
    async fn get_look_angles_uses_output_mask() {
        let mut position = MockPositionServiceApi::new();

        position
            .expect_get_look_angles()
            .withf(|identifier, datetime, observer, compute| {
                *identifier == DomainSatelliteIdentifier::NoradId(25544)
                    && *datetime == Utc.timestamp_opt(1_700_000_000, 0).unwrap()
                    && compute.azimuth
                    && compute.elevation
                    && !compute.range
                    && observer.lat.value > -std::f64::consts::PI
            })
            .times(1)
            .returning(|_, _, _, _| Ok((valid_look_angles(), valid_metadata())));

        let server = test_server(position);

        let request = Request::new(LookAnglesRequest {
            identifier: Some(norad_identifier(25544)),
            datetime: Some(valid_timestamp()),
            observer: Some(valid_observer()),
            output_mask: Some(FieldMask {
                paths: vec!["azimuth".to_string(), "elevation".to_string()],
            }),
            units: Some(valid_units()),
        });

        server
            .handle_look_angles(request)
            .await
            .expect("get_look_angles failed");
    }
}
