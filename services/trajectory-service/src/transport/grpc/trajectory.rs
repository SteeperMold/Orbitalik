use tonic::{Request, Response, Status};

use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::models::UnitContext;
use crate::domain::trajectory::TrajectoryServiceApi;
use crate::transport::grpc::server::TrajectoryGrpcServer;
use crate::transport::grpc::server::trajectory_grpc::{
    ObserverTrajectoryRequest, ObserverTrajectoryResponse, TrajectoryRequest, TrajectoryResponse,
};

impl<P, T, Pa> TrajectoryGrpcServer<P, T, Pa>
where
    P: Sync,
    T: TrajectoryServiceApi,
    Pa: Sync,
{
    pub async fn handle_get_trajectory(
        &self,
        request: Request<TrajectoryRequest>,
    ) -> Result<Response<TrajectoryResponse>, Status> {
        let req = request.into_inner();

        let identifier = req
            .identifier
            .ok_or_else(|| Status::invalid_argument("Missing satellite identifier"))?
            .try_into()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let sampling = req
            .sampling
            .ok_or_else(|| Status::invalid_argument("Missing sampling"))?
            .try_into()?;

        let unit_context = UnitContext::try_from(req.units)?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(PositionComputation::default, PositionComputation::from);

        let (trajectory, metadata) = self
            .trajectory
            .get_trajectory(identifier, range, sampling, &compute)
            .await?;

        let response = TrajectoryResponse::build(&trajectory, metadata, &unit_context)?;

        Ok(Response::new(response))
    }

    pub async fn handle_get_observer_trajectory(
        &self,
        request: Request<ObserverTrajectoryRequest>,
    ) -> Result<Response<ObserverTrajectoryResponse>, Status> {
        let req = request.into_inner();

        let identifier = req
            .identifier
            .ok_or_else(|| Status::invalid_argument("Missing satellite identifier"))?
            .try_into()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let sampling = req
            .sampling
            .ok_or_else(|| Status::invalid_argument("Missing sampling"))?
            .try_into()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let unit_context = UnitContext::try_from(req.units)?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(LookAnglesComputation::default, LookAnglesComputation::from);

        let (observer_trajectory, metadata) = self
            .trajectory
            .get_observer_trajectory(identifier, range, sampling, &observer, &compute)
            .await?;

        let response =
            ObserverTrajectoryResponse::build(&observer_trajectory, metadata, &unit_context)?;

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
        ObserverTrajectory, SatelliteIdentifier as DomainSatelliteIdentifier, Trajectory,
    };
    use crate::astro::propagation::position::PositionComputation;
    use crate::domain::errors::ServiceError;
    use crate::domain::models::TrajectoryComputationMetadata;
    use crate::domain::trajectory::MockTrajectoryServiceApi;
    use crate::transport::grpc::server::TrajectoryGrpcServer;
    use crate::transport::grpc::server::trajectory_grpc;
    use crate::transport::grpc::server::trajectory_grpc::{
        GeodeticInput, ObserverTrajectoryRequest, TrajectoryRequest, UnitSettings,
        unit_settings::{AngleUnit, DistanceUnit},
    };
    use crate::transport::grpc::server::trajectory_grpc::{SatelliteIdentifier, TimeRange};

    fn norad_identifier(id: u32) -> SatelliteIdentifier {
        SatelliteIdentifier {
            kind: Some(trajectory_grpc::satellite_identifier::Kind::NoradId(id)),
        }
    }

    fn valid_timestamp() -> Timestamp {
        Timestamp {
            seconds: Utc
                .with_ymd_and_hms(2024, 1, 1, 0, 0, 0)
                .single()
                .unwrap()
                .timestamp(),
            nanos: 0,
        }
    }

    fn valid_units() -> UnitSettings {
        UnitSettings {
            distance_unit: DistanceUnit::Kilometers as i32,
            angle_unit: AngleUnit::Degrees as i32,
        }
    }

    fn valid_range() -> TimeRange {
        TimeRange {
            start: Some(valid_timestamp()),
            end: Some(Timestamp {
                seconds: valid_timestamp().seconds + 3600,
                nanos: 0,
            }),
        }
    }

    fn valid_sampling() -> trajectory_grpc::SamplingOptions {
        trajectory_grpc::SamplingOptions {
            step_seconds: 60.0 as u32,
        }
    }

    fn valid_observer() -> GeodeticInput {
        GeodeticInput {
            lat: Some(trajectory_grpc::geodetic_input::Lat::LatDeg(50.0)),
            lon: Some(trajectory_grpc::geodetic_input::Lon::LonDeg(10.0)),
            alt: Some(trajectory_grpc::geodetic_input::Alt::AltM(100.0)),
        }
    }

    fn valid_metadata() -> TrajectoryComputationMetadata {
        TrajectoryComputationMetadata {
            propagation_model: "".to_string(),
            computation_time: Default::default(),
            norad_id: 0,
            satellite_name: "".to_string(),
            tle_epoch: Default::default(),
        }
    }

    fn valid_trajectory_request() -> TrajectoryRequest {
        TrajectoryRequest {
            identifier: Some(norad_identifier(25544)),
            range: Some(valid_range()),
            sampling: Some(valid_sampling()),
            units: Some(valid_units()),
            output_mask: None,
        }
    }

    fn valid_observer_trajectory_request() -> ObserverTrajectoryRequest {
        ObserverTrajectoryRequest {
            identifier: Some(norad_identifier(25544)),
            range: Some(valid_range()),
            sampling: Some(valid_sampling()),
            observer: Some(valid_observer()),
            units: Some(valid_units()),
            output_mask: None,
        }
    }

    fn test_server(
        trajectory: MockTrajectoryServiceApi,
    ) -> TrajectoryGrpcServer<
        crate::domain::position::MockPositionServiceApi,
        MockTrajectoryServiceApi,
        crate::domain::passes::MockPassesServiceApi,
    > {
        TrajectoryGrpcServer::new(
            crate::domain::position::MockPositionServiceApi::new(),
            trajectory,
            crate::domain::passes::MockPassesServiceApi::new(),
        )
    }

    #[tokio::test]
    async fn get_trajectory_rejects_missing_identifier() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.identifier = None;

        let error = server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing satellite identifier");
    }

    #[tokio::test]
    async fn get_trajectory_rejects_missing_range() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.range = None;

        let error = server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing range");
    }

    #[tokio::test]
    async fn get_trajectory_rejects_missing_sampling() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.sampling = None;

        let error = server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing sampling");
    }

    #[tokio::test]
    async fn get_trajectory_rejects_missing_units() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.units = None;

        let error = server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_trajectory_rejects_unspecified_distance_unit() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.units.as_mut().unwrap().distance_unit = DistanceUnit::Unspecified as i32;

        let error = server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_trajectory_rejects_unspecified_angle_unit() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.units.as_mut().unwrap().angle_unit = AngleUnit::Unspecified as i32;

        let error = server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_trajectory_uses_default_position_computation() {
        let mut trajectory = MockTrajectoryServiceApi::new();

        trajectory
            .expect_get_trajectory()
            .withf(|identifier, _, _, compute| {
                identifier == &DomainSatelliteIdentifier::NoradId(25544)
                    && compute == &PositionComputation::default()
            })
            .times(1)
            .returning(|_, _, _, _| {
                Ok((
                    Trajectory {
                        start: Default::default(),
                        end: Default::default(),
                        step_seconds: 0,
                        samples: vec![],
                    },
                    valid_metadata(),
                ))
            });

        let server = test_server(trajectory);

        server
            .handle_get_trajectory(Request::new(valid_trajectory_request()))
            .await
            .expect("get_trajectory failed");
    }

    #[tokio::test]
    async fn get_trajectory_uses_output_mask() {
        let mut trajectory = MockTrajectoryServiceApi::new();

        let mask = FieldMask {
            paths: vec!["position".to_string()],
        };

        let expected_compute = PositionComputation::from(&mask);

        trajectory
            .expect_get_trajectory()
            .withf(move |identifier, _, _, compute| {
                identifier == &DomainSatelliteIdentifier::NoradId(25544)
                    && compute == &expected_compute
            })
            .times(1)
            .returning(|_, _, _, _| {
                Ok((
                    Trajectory {
                        start: Default::default(),
                        end: Default::default(),
                        step_seconds: 0,
                        samples: vec![],
                    },
                    valid_metadata(),
                ))
            });

        let server = test_server(trajectory);

        let mut request = valid_trajectory_request();
        request.output_mask = Some(mask);

        server
            .handle_get_trajectory(Request::new(request))
            .await
            .expect("get_trajectory failed");
    }

    #[tokio::test]
    async fn get_trajectory_propagates_service_error() {
        let mut trajectory = MockTrajectoryServiceApi::new();

        trajectory
            .expect_get_trajectory()
            .times(1)
            .returning(|_, _, _, _| {
                Err(ServiceError::from(Status::internal(
                    "trajectory service failed",
                )))
            });

        let server = test_server(trajectory);

        let error = server
            .handle_get_trajectory(Request::new(valid_trajectory_request()))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::Internal);
        assert_eq!(error.message(), "trajectory service failed");
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_missing_identifier() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.identifier = None;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing satellite identifier");
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_missing_range() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.range = None;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing range");
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_missing_sampling() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.sampling = None;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing sampling");
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_missing_observer() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.observer = None;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing observer");
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_missing_units() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.units = None;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_unspecified_distance_unit() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.units.as_mut().unwrap().distance_unit = DistanceUnit::Unspecified as i32;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_observer_trajectory_rejects_unspecified_angle_unit() {
        let trajectory = MockTrajectoryServiceApi::new();
        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.units.as_mut().unwrap().angle_unit = AngleUnit::Unspecified as i32;

        let error = server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_observer_trajectory_uses_default_look_angles_computation() {
        let mut trajectory = MockTrajectoryServiceApi::new();

        trajectory
            .expect_get_observer_trajectory()
            .times(1)
            .returning(|_, _, _, _, _| {
                Ok((
                    ObserverTrajectory {
                        start: Default::default(),
                        end: Default::default(),
                        step_seconds: 0,
                        samples: vec![],
                    },
                    valid_metadata(),
                ))
            });

        let server = test_server(trajectory);

        server
            .handle_get_observer_trajectory(Request::new(valid_observer_trajectory_request()))
            .await
            .expect("get_observer_trajectory failed");
    }

    #[tokio::test]
    async fn get_observer_trajectory_uses_output_mask() {
        let mut trajectory = MockTrajectoryServiceApi::new();

        let mask = FieldMask {
            paths: vec!["azimuth".to_string()],
        };

        trajectory
            .expect_get_observer_trajectory()
            .times(1)
            .returning(|_, _, _, _, _| {
                Ok((
                    ObserverTrajectory {
                        start: Default::default(),
                        end: Default::default(),
                        step_seconds: 0,
                        samples: vec![],
                    },
                    valid_metadata(),
                ))
            });

        let server = test_server(trajectory);

        let mut request = valid_observer_trajectory_request();
        request.output_mask = Some(mask);

        server
            .handle_get_observer_trajectory(Request::new(request))
            .await
            .expect("get_observer_trajectory failed");
    }

    #[tokio::test]
    async fn get_observer_trajectory_propagates_service_error() {
        let mut trajectory = MockTrajectoryServiceApi::new();

        trajectory
            .expect_get_observer_trajectory()
            .times(1)
            .returning(|_, _, _, _, _| {
                Err(ServiceError::from(Status::internal(
                    "trajectory service failed",
                )))
            });

        let server = test_server(trajectory);

        let error = server
            .handle_get_observer_trajectory(Request::new(valid_observer_trajectory_request()))
            .await
            .expect_err("request should fail");

        assert_eq!(error.code(), Code::Internal);
        assert_eq!(error.message(), "trajectory service failed");
    }
}
