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
            .min_elevation
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
            .min_elevation
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
