use tonic::{Request, Response, Status};
use uom::si::angle::degree;
use uom::si::f64::Angle;

use crate::astro::models::SatelliteIdentifier;
use crate::domain::models::UnitContext;
use crate::domain::passes::PassesServiceApi;
use crate::service::passes::{GetPassesOptions, NextPassesOptions};
use crate::transport::grpc::conversions::requests::IntoAngle;
use crate::transport::grpc::server::TrajectoryGrpcServer;
use crate::transport::grpc::server::trajectory_grpc::{
    GetPassesResponse, NextPassesRequest, NextPassesResponse, PassPredictionRequest,
};

impl<P, T, Pa> TrajectoryGrpcServer<P, T, Pa>
where
    P: Sync,
    T: Sync,
    Pa: PassesServiceApi,
{
    pub async fn handle_get_passes(
        &self,
        request: Request<PassPredictionRequest>,
    ) -> Result<Response<GetPassesResponse>, Status> {
        let req = request.into_inner();

        let satellites: Vec<SatelliteIdentifier> = req
            .satellites
            .into_iter()
            .map(TryInto::try_into)
            .collect::<Result<_, _>>()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let unit_context = UnitContext::try_from(req.units)?;

        let min_elevation = req
            .min_elevation
            .map_or_else(|| Angle::new::<degree>(0.0), IntoAngle::into_angle);

        let min_peak_elevation = req
            .min_peak_elevation
            .map_or_else(|| Angle::new::<degree>(0.0), IntoAngle::into_angle);

        let prediction_options = GetPassesOptions {
            range,
            observer,
            min_elevation,
            min_peak_elevation,
            max_results: req.max_results.map(|n| n as usize),
        };

        let (passes, metadata) = self
            .passes
            .get_passes(satellites, &prediction_options)
            .await?;

        let response = GetPassesResponse::build(&passes, metadata, range, &unit_context)?;

        Ok(Response::new(response))
    }

    pub async fn handle_get_next_passes(
        &self,
        request: Request<NextPassesRequest>,
    ) -> Result<Response<NextPassesResponse>, Status> {
        let req = request.into_inner();

        let satellites: Vec<SatelliteIdentifier> = req
            .satellites
            .into_iter()
            .map(TryInto::try_into)
            .collect::<Result<_, _>>()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let unit_context = UnitContext::try_from(req.units)?;

        let min_elevation = req
            .min_elevation
            .map_or_else(|| Angle::new::<degree>(0.0), IntoAngle::into_angle);

        let min_peak_elevation = req
            .min_peak_elevation
            .map_or_else(|| Angle::new::<degree>(0.0), IntoAngle::into_angle);

        let prediction_options = NextPassesOptions {
            observer,
            min_elevation,
            min_peak_elevation,
            passes_count: req.count as usize,
        };

        let (passes, metadata) = self
            .passes
            .next_passes(satellites, &prediction_options)
            .await?;

        let response = NextPassesResponse::build(&passes, metadata, &unit_context)?;

        Ok(Response::new(response))
    }
}

#[allow(clippy::unwrap_used)]
#[cfg(test)]
mod tests {
    use chrono::{TimeZone, Utc};
    use prost_types::Timestamp;
    use tonic::{Code, Request, Status};
    use uom::si::angle::{degree, radian};

    use crate::astro::models::SatelliteIdentifier as DomainSatelliteIdentifier;
    use crate::domain::models::PassesComputationMetadata;
    use crate::domain::passes::MockPassesServiceApi;
    use crate::domain::position::MockPositionServiceApi;
    use crate::domain::trajectory::MockTrajectoryServiceApi;
    use crate::transport::grpc::server::TrajectoryGrpcServer;
    use crate::transport::grpc::server::trajectory_grpc::{
        NextPassesRequest, PassPredictionRequest, SatelliteIdentifier, TimeRange, UnitSettings,
        next_passes_request, pass_prediction_request, satellite_identifier,
        unit_settings::{AngleUnit, DistanceUnit},
    };

    fn norad_identifier(id: u32) -> SatelliteIdentifier {
        SatelliteIdentifier {
            kind: Some(satellite_identifier::Kind::NoradId(id)),
        }
    }

    fn valid_units() -> UnitSettings {
        UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Degrees as i32,
        }
    }

    fn valid_radian_units() -> UnitSettings {
        UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Radians as i32,
        }
    }

    fn valid_metadata() -> PassesComputationMetadata {
        PassesComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc.timestamp_opt(1_700_000_000, 0).unwrap(),
            norad_ids: vec![25544],
            tle_epoch: Utc.timestamp_opt(1_700_000_000, 0).unwrap(),
            satellites_evaluated: 0,
            passes_found: 0,
            satellite_names: vec![],
            computation_ms: 0,
        }
    }

    fn valid_range() -> TimeRange {
        TimeRange {
            start: Some(Timestamp {
                seconds: 1_700_000_000,
                nanos: 0,
            }),
            end: Some(Timestamp {
                seconds: 1_700_360_000,
                nanos: 0,
            }),
        }
    }

    fn valid_pass_prediction_request() -> PassPredictionRequest {
        PassPredictionRequest {
            satellites: vec![norad_identifier(25544)],
            range: Some(valid_range()),
            observer: None,
            min_elevation: None,
            min_peak_elevation: None,
            max_results: None,
            units: Some(valid_units()),
        }
    }

    fn valid_next_passes_request() -> NextPassesRequest {
        NextPassesRequest {
            satellites: vec![norad_identifier(25544)],
            observer: None,
            min_elevation: None,
            min_peak_elevation: None,
            count: 5,
            units: Some(valid_units()),
        }
    }

    fn test_get_passes_server(
        passes: MockPassesServiceApi,
    ) -> TrajectoryGrpcServer<MockPositionServiceApi, MockTrajectoryServiceApi, MockPassesServiceApi>
    {
        TrajectoryGrpcServer::new(
            MockPositionServiceApi::new(),
            MockTrajectoryServiceApi::new(),
            passes,
        )
    }

    fn test_next_passes_server(
        passes: MockPassesServiceApi,
    ) -> TrajectoryGrpcServer<MockPositionServiceApi, MockTrajectoryServiceApi, MockPassesServiceApi>
    {
        TrajectoryGrpcServer::new(
            MockPositionServiceApi::new(),
            MockTrajectoryServiceApi::new(),
            passes,
        )
    }

    #[tokio::test]
    async fn get_passes_rejects_missing_range() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_get_passes().never();

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.range = None;

        let error = server
            .handle_get_passes(Request::new(request))
            .await
            .expect_err("expected missing range to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing range");
    }

    #[tokio::test]
    async fn get_passes_rejects_missing_observer() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_get_passes().never();

        let server = test_get_passes_server(passes);

        let request = valid_pass_prediction_request();

        let error = server
            .handle_get_passes(Request::new(request))
            .await
            .expect_err("expected missing observer to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing observer");
    }

    #[tokio::test]
    async fn get_passes_rejects_missing_units() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_get_passes().never();

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());
        request.units = None;

        let error = server
            .handle_get_passes(Request::new(request))
            .await
            .expect_err("expected missing units to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Measurement units must be specified");
    }

    #[tokio::test]
    async fn get_passes_rejects_unspecified_distance_unit() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_get_passes().never();

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());

        request.units = Some(UnitSettings {
            distance_unit: DistanceUnit::Unspecified as i32,
            angle_unit: AngleUnit::Degrees as i32,
        });

        let error = server
            .handle_get_passes(Request::new(request))
            .await
            .expect_err("expected unspecified distance unit to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_passes_rejects_unspecified_angle_unit() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_get_passes().never();

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());

        request.units = Some(UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Unspecified as i32,
        });

        let error = server
            .handle_get_passes(Request::new(request))
            .await
            .expect_err("expected unspecified angle unit to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn get_passes_uses_default_elevations() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_get_passes()
            .withf(|_satellites, options| {
                options.min_elevation.get::<degree>() == 0.0
                    && options.min_peak_elevation.get::<degree>() == 0.0
            })
            .times(1)
            .returning(|_, _| Ok((vec![], valid_metadata())));

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());

        server
            .handle_get_passes(Request::new(request))
            .await
            .expect("get_passes failed");
    }

    #[tokio::test]
    async fn get_passes_uses_radian_elevations() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_get_passes()
            .withf(|_, options| {
                (options.min_elevation.get::<radian>() - 0.1).abs() < f64::EPSILON
                    && (options.min_peak_elevation.get::<radian>() - 0.5).abs() < f64::EPSILON
            })
            .times(1)
            .returning(|_, _| Ok((vec![], valid_metadata())));

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());
        request.units = Some(valid_radian_units());

        request.min_elevation = Some(pass_prediction_request::MinElevation::MinElevationRad(0.1));
        request.min_peak_elevation =
            Some(pass_prediction_request::MinPeakElevation::MinPeakElevationRad(0.5));

        server
            .handle_get_passes(Request::new(request))
            .await
            .expect("get_passes failed");
    }

    #[tokio::test]
    async fn get_passes_uses_max_results() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_get_passes()
            .withf(|_, options| options.max_results == Some(10))
            .times(1)
            .returning(|_, _| Ok((vec![], valid_metadata())));

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());
        request.max_results = Some(10);

        server
            .handle_get_passes(Request::new(request))
            .await
            .expect("get_passes failed");
    }

    #[tokio::test]
    async fn get_passes_propagates_service_error() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_get_passes()
            .times(1)
            .returning(|_, _| Err(Status::internal("pass calculation failed").into()));

        let server = test_get_passes_server(passes);

        let mut request = valid_pass_prediction_request();
        request.observer = Some(valid_observer());

        let error = server
            .handle_get_passes(Request::new(request))
            .await
            .expect_err("expected get_passes to fail");

        assert_eq!(error.code(), Code::Internal);
        assert_eq!(error.message(), "pass calculation failed");
    }

    #[tokio::test]
    async fn next_passes_rejects_missing_observer() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_next_passes().never();

        let server = test_next_passes_server(passes);

        let request = valid_next_passes_request();

        let error = server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect_err("expected missing observer to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Missing observer");
    }

    #[tokio::test]
    async fn next_passes_rejects_missing_units() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_next_passes().never();

        let server = test_next_passes_server(passes);

        let mut request = valid_next_passes_request();
        request.observer = Some(valid_observer());
        request.units = None;

        let error = server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect_err("expected missing units to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
        assert_eq!(error.message(), "Measurement units must be specified");
    }

    #[tokio::test]
    async fn next_passes_rejects_unspecified_distance_unit() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_next_passes().never();

        let server = test_next_passes_server(passes);

        let mut request = valid_next_passes_request();
        request.observer = Some(valid_observer());
        request.units = Some(UnitSettings {
            distance_unit: DistanceUnit::Unspecified as i32,
            angle_unit: AngleUnit::Degrees as i32,
        });

        let error = server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect_err("expected unspecified distance unit to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn next_passes_rejects_unspecified_angle_unit() {
        let mut passes = MockPassesServiceApi::new();

        passes.expect_next_passes().never();

        let server = test_next_passes_server(passes);

        let mut request = valid_next_passes_request();
        request.observer = Some(valid_observer());
        request.units = Some(UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Unspecified as i32,
        });

        let error = server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect_err("expected unspecified angle unit to fail");

        assert_eq!(error.code(), Code::InvalidArgument);
    }

    #[tokio::test]
    async fn next_passes_uses_default_elevations() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_next_passes()
            .withf(|satellites, options| {
                satellites == &[DomainSatelliteIdentifier::NoradId(25544)]
                    && options.min_elevation.get::<degree>() == 0.0
                    && options.min_peak_elevation.get::<degree>() == 0.0
                    && options.passes_count == 5
            })
            .times(1)
            .returning(|_, _| Ok((vec![], valid_metadata())));

        let server = test_next_passes_server(passes);

        let mut request = valid_next_passes_request();
        request.observer = Some(valid_observer());

        server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect("next_passes failed");
    }

    #[tokio::test]
    async fn next_passes_uses_radian_elevations() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_next_passes()
            .withf(|_, options| {
                (options.min_elevation.get::<radian>() - 0.1).abs() < f64::EPSILON
                    && (options.min_peak_elevation.get::<radian>() - 0.5).abs() < f64::EPSILON
            })
            .times(1)
            .returning(|_, _| Ok((vec![], valid_metadata())));

        let server = test_next_passes_server(passes);

        let mut request = valid_next_passes_request();
        request.observer = Some(valid_observer());
        request.units = Some(valid_radian_units());

        request.min_elevation = Some(next_passes_request::MinElevation::MinElevationRad(0.1));
        request.min_peak_elevation = Some(
            next_passes_request::MinPeakElevation::MinPeakElevationRad(0.5),
        );

        server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect("next_passes failed");
    }

    #[tokio::test]
    async fn next_passes_propagates_service_error() {
        let mut passes = MockPassesServiceApi::new();

        passes
            .expect_next_passes()
            .times(1)
            .returning(|_, _| Err(Status::internal("next pass calculation failed").into()));

        let server = test_next_passes_server(passes);

        let mut request = valid_next_passes_request();
        request.observer = Some(valid_observer());

        let error = server
            .handle_get_next_passes(Request::new(request))
            .await
            .expect_err("expected next_passes to fail");

        assert_eq!(error.code(), Code::Internal);
        assert_eq!(error.message(), "next pass calculation failed");
    }

    fn valid_observer() -> crate::transport::grpc::server::trajectory_grpc::GeodeticInput {
        use crate::transport::grpc::server::trajectory_grpc::GeodeticInput;
        use crate::transport::grpc::server::trajectory_grpc::geodetic_input::{Alt, Lat, Lon};

        GeodeticInput {
            lat: Some(Lat::LatDeg(56.9496)),
            lon: Some(Lon::LonDeg(24.1052)),
            alt: Some(Alt::AltM(0.0)),
        }
    }
}
