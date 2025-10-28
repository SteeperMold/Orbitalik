use thiserror::Error;

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
pub enum PropagationError {
    #[error("Failed to fetch TLE: {0}")]
    TleFetch(#[from] tonic::Status),
    #[error("TLE parsing failed: {0}")]
    TleParse(#[from] sgp4::TleError),
    #[error("Failed to create SGP4 elements: {0}")]
    ElementsCreation(#[from] sgp4::ElementsError),
    #[error("Failed to convert datetime into minutes since epoch: {0}")]
    DatetimeToMinutesSinceEpochFailed(#[from] sgp4::DatetimeToMinutesSinceEpochError),
    #[error("SGP4 propagation failed: {0}")]
    PropagationFailed(#[from] sgp4::Error),
}

impl From<PropagationError> for tonic::Status {
    fn from(value: PropagationError) -> Self {
        match value {
            PropagationError::TleFetch(status) => status,
            PropagationError::TleParse(_)
            | PropagationError::ElementsCreation(_)
            | PropagationError::DatetimeToMinutesSinceEpochFailed(_)
            | PropagationError::PropagationFailed(_) => {
                tracing::error!("propagation error: {:?}", value);
                Self::internal("Internal server error")
            }
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
