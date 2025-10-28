use tonic::{Request, Response, Status};

use crate::service::look_angles::LookAnglesService;
use crate::service::position::PositionService;
use crate::transport::grpc::converters::ToChrono;
use trajectory_grpc::{
    LookAnglesRequest, LookAnglesResponse, PositionRequest, PositionResponse,
    trajectory_service_server::TrajectoryService,
};

pub mod trajectory_grpc {
    tonic::include_proto!("trajectory");
}

pub struct TrajectoryGrpc {
    position_service: PositionService,
    look_angles_service: LookAnglesService,
}

impl TrajectoryGrpc {
    pub const fn new(
        position_service: PositionService,
        look_angles_service: LookAnglesService,
    ) -> Self {
        Self {
            position_service,
            look_angles_service,
        }
    }
}

#[tonic::async_trait]
impl TrajectoryService for TrajectoryGrpc {
    async fn get_position(
        &self,
        request: Request<PositionRequest>,
    ) -> Result<Response<PositionResponse>, Status> {
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

        let sat_pos = self
            .position_service
            .calculate_position(identifier, datetime)
            .await?;

        let response = PositionResponse::from(sat_pos);
        Ok(Response::new(response))
    }

    async fn get_look_angles(
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

        let look_angles = self
            .look_angles_service
            .calculate_look_angles(identifier, datetime, &observer)
            .await?;

        let response = LookAnglesResponse::from(look_angles);
        Ok(Response::new(response))
    }
}
