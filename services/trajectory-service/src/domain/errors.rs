use thiserror::Error;

use crate::astro;

#[derive(Debug, Error)]
pub enum StartupError {
    #[error("Invalid APP_ENV: {0}")]
    InvalidAppEnv(String),
    #[error("Failed to initialize logger: {0}")]
    LoggerInit(String),
    #[error("HTTP server error: {0}")]
    Http(#[from] HttpServerError),
    #[error("gRPC server error: {0}")]
    Grpc(#[from] GrpcServerError),
    #[error("Failed to initialize TLE gRPC client: {0}")]
    GrpcClientInit(#[from] tonic::transport::Error),
}

impl From<std::io::Error> for StartupError {
    fn from(e: std::io::Error) -> Self {
        Self::Http(HttpServerError::from(e))
    }
}

#[derive(Debug, Error)]
pub enum HttpServerError {
    #[error("Failed to build Prometheus middleware: {0}")]
    Prometheus(String),
    #[error("I/O error while starting HTTP server: {0}")]
    Io(#[from] std::io::Error),
}

#[derive(Debug, Error)]
pub enum GrpcServerError {
    #[error("I/O error while starting gRPC server: {0}")]
    Io(#[from] std::io::Error),
    #[error("transport error: {0}")]
    Transport(#[from] tonic::transport::Error),
}

#[derive(Debug, Error)]
pub enum ServiceError {
    #[error("Failed to fetch TLE: {0}")]
    TleFetch(#[from] tonic::Status),
    #[error("No satellites provided")]
    NoSatellites,
    #[error("Too many satellites: {provided} (max {max})")]
    TooManySatellites { provided: usize, max: usize },
    #[error(transparent)]
    Propagation(#[from] astro::errors::PropagationError),
}

impl From<ServiceError> for tonic::Status {
    fn from(value: ServiceError) -> Self {
        match value {
            ServiceError::TleFetch(status) => status,
            ServiceError::Propagation(e) => {
                tracing::error!("{:?}", e);
                Self::internal("Internal server error")
            }
            other => Self::invalid_argument(other.to_string()),
        }
    }
}

#[derive(Debug, Error)]
pub enum TimestampConversionError {
    #[error("Failed to convert nanos: {0}")]
    NanosOutOfRange(#[from] std::num::TryFromIntError),
    #[error("Failed to convert prost_types::Timestamp to chrono::DateTime<Utc>")]
    InvalidTimestamp,
}

impl From<TimestampConversionError> for tonic::Status {
    fn from(value: TimestampConversionError) -> Self {
        tracing::error!("timestamp conversion failed: {:?}", value);
        Self::internal("Internal server error")
    }
}
