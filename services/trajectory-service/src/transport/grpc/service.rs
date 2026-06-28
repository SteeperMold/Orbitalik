use tonic::{Request, Response, Status};

use crate::service::passes::PassesService;
use crate::service::position::PositionService;
use crate::service::trajectory::TrajectoryService;
use crate::transport::grpc::service::trajectory_grpc::{GetPassesResponse, NextPassesRequest, NextPassesResponse, ObserverTrajectoryRequest, ObserverTrajectoryResponse, PassPredictionRequest, TrajectoryRequest, TrajectoryResponse};
use trajectory_grpc::{
    LookAnglesRequest, LookAnglesResponse, PositionRequest, PositionResponse,
    trajectory_service_server::TrajectoryService as TonicTrajectoryService,
};

pub mod trajectory_grpc {
    tonic::include_proto!("trajectory");
}

pub struct TrajectoryGrpcServer {
    pub position: PositionService,
    pub trajectory: TrajectoryService,
    pub passes: PassesService,
}

impl TrajectoryGrpcServer {
    pub const fn new(
        position_service: PositionService,
        trajectory_service: TrajectoryService,
        passes_service: PassesService,
    ) -> Self {
        Self {
            position: position_service,
            trajectory: trajectory_service,
            passes: passes_service,
        }
    }
}

#[tonic::async_trait]
impl TonicTrajectoryService for TrajectoryGrpcServer {
    async fn get_position(
        &self,
        request: Request<PositionRequest>,
    ) -> Result<Response<PositionResponse>, Status> {
        self.handle_get_position(request).await
    }

    async fn get_look_angles(
        &self,
        request: Request<LookAnglesRequest>,
    ) -> Result<Response<LookAnglesResponse>, Status> {
        self.handle_look_angles(request).await
    }

    async fn get_trajectory(
        &self,
        request: Request<TrajectoryRequest>,
    ) -> Result<Response<TrajectoryResponse>, Status> {
        self.handle_get_trajectory(request).await
    }

    async fn get_observer_trajectory(
        &self,
        request: Request<ObserverTrajectoryRequest>,
    ) -> Result<Response<ObserverTrajectoryResponse>, Status> {
        self.handle_get_observer_trajectory(request).await
    }

    async fn get_passes(
        &self,
        request: Request<PassPredictionRequest>,
    ) -> Result<Response<GetPassesResponse>, Status> {
        self.handle_get_passes(request).await
    }

    async fn get_next_passes(
        &self,
        request: Request<NextPassesRequest>,
    ) -> Result<Response<NextPassesResponse>, Status> {
        self.handle_get_next_passes(request).await
    }
}
