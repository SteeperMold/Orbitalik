use tonic::transport::Server;

use crate::domain::errors::GrpcServerError;
use crate::service::position::PositionService;
use crate::service::trajectory::TrajectoryService;
use crate::transport::grpc::interceptors::LoggingMiddlewareLayer;
use crate::transport::grpc::service::{
    TrajectoryGrpcServer, trajectory_grpc::trajectory_service_server::TrajectoryServiceServer,
};

pub async fn run(
    port: u16,
    position_service: PositionService,
    trajectory_service: TrajectoryService,
) -> Result<(), GrpcServerError> {
    let trajectory_service = TrajectoryGrpcServer::new(position_service, trajectory_service);

    let layer = tower::ServiceBuilder::new()
        .layer(LoggingMiddlewareLayer::default())
        .into_inner();

    Server::builder()
        .layer(layer)
        .add_service(TrajectoryServiceServer::new(trajectory_service))
        .serve(([0, 0, 0, 0], port).into())
        .await
        .map_err(GrpcServerError::from)
}
