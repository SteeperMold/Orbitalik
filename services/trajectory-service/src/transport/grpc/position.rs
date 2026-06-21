use tonic::{Request, Response, Status};

use crate::astro::look_angles::LookAnglesComputation;
use crate::astro::position::PositionComputation;
use crate::transport::grpc::converters::ToChrono;
use crate::transport::grpc::service::TrajectoryGrpcServer;
use crate::transport::grpc::service::trajectory_grpc::{
    LookAnglesRequest, LookAnglesResponse, PositionRequest, PositionResponse,
};

impl TrajectoryGrpcServer {
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
            .to_chrono()?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(PositionComputation::default, PositionComputation::from);

        let (position, metadata) = self
            .position_service
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
            .position_service
            .get_look_angles(identifier, datetime, &observer, &compute)
            .await?;

        let response = LookAnglesResponse::from_look_angles(&look_angles, metadata, req.units)?;
        Ok(Response::new(response))
    }
}
