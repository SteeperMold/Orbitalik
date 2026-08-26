use tonic::{Request, Response, Status, transport::Server};

use crate::{
    domain::{errors::GrpcServerError, position::PositionServiceApi},
    service::passes::PassesService,
    service::position::PositionService,
    service::trajectory::TrajectoryService,
    transport::grpc::{
        interceptors::LoggingMiddlewareLayer, server::trajectory_grpc::GetPassesResponse,
        server::trajectory_grpc::NextPassesRequest, server::trajectory_grpc::NextPassesResponse,
        server::trajectory_grpc::ObserverTrajectoryRequest,
        server::trajectory_grpc::ObserverTrajectoryResponse,
        server::trajectory_grpc::PassPredictionRequest, server::trajectory_grpc::TrajectoryRequest,
        server::trajectory_grpc::TrajectoryResponse,
    },
};
use trajectory_grpc::{
    LookAnglesRequest, LookAnglesResponse, PositionRequest, PositionResponse,
    trajectory_service_server::TrajectoryService as TonicTrajectoryService,
    trajectory_service_server::TrajectoryServiceServer,
};
use crate::domain::passes::PassesServiceApi;
use crate::domain::trajectory::TrajectoryServiceApi;

pub mod trajectory_grpc {
    tonic::include_proto!("trajectory");
}

pub async fn run(
    port: u16,
    position_service: PositionService,
    trajectory_service: TrajectoryService,
    passes_service: PassesService,
) -> Result<(), GrpcServerError> {
    let service = TrajectoryGrpcServer::new(position_service, trajectory_service, passes_service);

    let layer = tower::ServiceBuilder::new()
        .layer(LoggingMiddlewareLayer::default())
        .into_inner();

    Server::builder()
        .layer(layer)
        .add_service(TrajectoryServiceServer::new(service))
        .serve(([0, 0, 0, 0], port).into())
        .await
        .map_err(GrpcServerError::from)
}

pub struct TrajectoryGrpcServer<P, T, Pa> {
    pub position: P,
    pub trajectory: T,
    pub passes: Pa,
}

impl<P, T, Pa> TrajectoryGrpcServer<P, T, Pa> {
    pub const fn new(
        position_service: P,
        trajectory_service: T,
        passes_service: Pa,
    ) -> Self {
        Self {
            position: position_service,
            trajectory: trajectory_service,
            passes: passes_service,
        }
    }
}

#[tonic::async_trait]
impl<P, T, Pa> TonicTrajectoryService for TrajectoryGrpcServer<P, T, Pa>
where
    P: PositionServiceApi + 'static,
    T: TrajectoryServiceApi + 'static,
    Pa: PassesServiceApi + 'static,
{
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

