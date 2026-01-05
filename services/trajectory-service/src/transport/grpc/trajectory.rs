use tonic::{Request, Response, Status};

use crate::astro::look_angles::LookAnglesComputation;
use crate::astro::position::PositionComputation;
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

pub struct TrajectoryGrpcServer {
    position_service: PositionService,
    look_angles_service: LookAnglesService,
}

impl TrajectoryGrpcServer {
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
impl TrajectoryService for TrajectoryGrpcServer {
    async fn get_position(
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
            .get_position_with_metadata(identifier, datetime, &compute)
            .await?;

        let response = PositionResponse::from_position(&position, metadata, req.units)?;
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

        let mask = req.output_mask.as_ref();

        let compute = mask.map_or_else(LookAnglesComputation::default, LookAnglesComputation::from);

        let (look_angles, metadata) = self
            .look_angles_service
            .get_look_angles_with_metadata(identifier, datetime, &observer, &compute)
            .await?;

        let response = LookAnglesResponse::from_look_angles(&look_angles, metadata, req.units)?;
        Ok(Response::new(response))
    }
}
