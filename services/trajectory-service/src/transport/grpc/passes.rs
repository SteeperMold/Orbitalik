use tonic::{Request, Response, Status};
use uom::si::angle::{degree, radian};
use uom::si::f64::Angle;

use crate::astro::models::SatelliteIdentifier;
use crate::astro::passes::detector::PassPredictionOptions;
use crate::transport::grpc::service::TrajectoryGrpcServer;
use crate::transport::grpc::service::trajectory_grpc::{
    NextPassesRequest, PassPredictionRequest, PassPredictionResponse, pass_prediction_request,
};

impl TrajectoryGrpcServer {
    pub async fn handle_get_passes(
        &self,
        request: Request<PassPredictionRequest>,
    ) -> Result<Response<PassPredictionResponse>, Status> {
        let req = request.into_inner();

        let satellites: Vec<SatelliteIdentifier> = req
            .satellites
            .into_iter()
            .map(|s| s.try_into())
            .collect::<Result<_, _>>()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let observer = &req
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

        let prediction_options = PassPredictionOptions {
            range,
            observer,
            min_elevation,
            min_peak_elevation,
            max_results: req.max_results.map(|n| n as usize),
        };

        let (passes, metadata) = self
            .passes_service
            .get_passes(satellites, &prediction_options)
            .await?;

        let response = PassPredictionResponse::from_passes(&passes, metadata, range, req.units)?;

        Ok(Response::new(response))
    }

    pub async fn handle_get_next_passes(
        &self,
        _request: Request<NextPassesRequest>,
    ) -> Result<Response<PassPredictionResponse>, Status> {
        todo!()
    }
}
