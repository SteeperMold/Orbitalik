use tonic::{Request, Response, Status};
use uom::si::angle::{degree, radian};
use uom::si::f64::Angle;

use crate::astro::models::SatelliteIdentifier;
use crate::domain::passes::PassesServiceApi;
use crate::service::passes::{GetPassesOptions, NextPassesOptions};
use crate::transport::grpc::server::TrajectoryGrpcServer;
use crate::transport::grpc::server::trajectory_grpc::{
    GetPassesResponse, NextPassesRequest, NextPassesResponse, PassPredictionRequest,
    next_passes_request, pass_prediction_request,
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
            .map(std::convert::TryInto::try_into)
            .collect::<Result<_, _>>()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let min_elevation = match req.min_elevation {
            Some(pass_prediction_request::MinElevation::MinElevationDeg(v)) => {
                Angle::new::<degree>(v)
            }
            Some(pass_prediction_request::MinElevation::MinElevationRad(v)) => {
                Angle::new::<radian>(v)
            }
            None => Angle::new::<degree>(0.0),
        };
        let min_peak_elevation = match req.min_peak_elevation {
            Some(pass_prediction_request::MinPeakElevation::MinPeakElevationDeg(v)) => {
                Angle::new::<degree>(v)
            }
            Some(pass_prediction_request::MinPeakElevation::MinPeakElevationRad(v)) => {
                Angle::new::<radian>(v)
            }
            None => Angle::new::<degree>(0.0),
        };

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

        let response = GetPassesResponse::from_passes(&passes, metadata, range, req.units)?;

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
            .map(std::convert::TryInto::try_into)
            .collect::<Result<_, _>>()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let min_elevation = match req.min_elevation {
            Some(next_passes_request::MinElevation::MinElevationDeg(v)) => Angle::new::<degree>(v),
            Some(next_passes_request::MinElevation::MinElevationRad(v)) => Angle::new::<radian>(v),
            None => Angle::new::<degree>(0.0),
        };
        let min_peak_elevation = match req.min_peak_elevation {
            Some(next_passes_request::MinPeakElevation::MinPeakElevationDeg(v)) => {
                Angle::new::<degree>(v)
            }
            Some(next_passes_request::MinPeakElevation::MinPeakElevationRad(v)) => {
                Angle::new::<radian>(v)
            }
            None => Angle::new::<degree>(0.0),
        };

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

        let response = NextPassesResponse::from_passes(&passes, metadata, req.units)?;

        Ok(Response::new(response))
    }
}
