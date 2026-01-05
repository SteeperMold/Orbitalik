use tonic::transport::Server;

use crate::domain::errors::GrpcServerError;
use crate::service::look_angles::LookAnglesService;
use crate::service::position::PositionService;
use crate::transport::grpc::interceptors::LoggingMiddlewareLayer;
use crate::transport::grpc::trajectory::{
    TrajectoryGrpcServer, trajectory_grpc::trajectory_service_server::TrajectoryServiceServer,
};

pub async fn run(
    port: u16,
    position_service: PositionService,
    look_angles_service: LookAnglesService,
) -> Result<(), GrpcServerError> {
    let trajectory_service = TrajectoryGrpcServer::new(position_service, look_angles_service);

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
