//! Project-wide lint policy
#![deny(
    warnings,
    unused_must_use,
    clippy::unwrap_used,
    clippy::expect_used,
    clippy::dbg_macro,
    clippy::print_stdout,
    clippy::todo,
    clippy::unimplemented,
    clippy::panic
)]
#![warn(missing_docs, clippy::all, clippy::pedantic, clippy::nursery)]
#![allow(clippy::suboptimal_flops)]

use std::sync::Arc;

use crate::domain::errors::StartupError;
use crate::service::look_angles::LookAnglesService;
use crate::service::position::PositionService;
use crate::transport::adapter::tle_client::TleGrpcClient;

mod astro;
mod domain;
mod infrastructure;
mod service;
mod transport;

#[actix_web::main]
async fn main() -> Result<(), StartupError> {
    let config = infrastructure::config::AppConfig::from_dotenv();

    infrastructure::logger::init_logger(&config.app_env)?;

    let tle_grpc_client = Arc::new(
        TleGrpcClient::new(config.tle_service_address)
            .await
            .map_err(StartupError::from)?,
    );

    let position_service = PositionService::new(tle_grpc_client.clone());
    let look_angles_service = LookAnglesService::new(tle_grpc_client.clone());

    let http_server = transport::http::server::run(config.http_port)?;
    let grpc_server =
        transport::grpc::server::run(config.grpc_port, position_service, look_angles_service);

    tokio::try_join!(
        async { http_server.await.map_err(StartupError::from) },
        async { grpc_server.await.map_err(StartupError::from) },
    )?;

    Ok(())
}
