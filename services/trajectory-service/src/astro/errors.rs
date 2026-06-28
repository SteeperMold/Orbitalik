use thiserror::Error;

#[derive(Debug, Error)]
pub enum PropagationError {
    #[error("TLE parsing failed: {0}")]
    TleParse(#[from] sgp4::TleError),
    #[error("Failed to create SGP4 elements: {0}")]
    ElementsCreation(#[from] sgp4::ElementsError),
    #[error("Failed to convert datetime into minutes since epoch: {0}")]
    DatetimeToMinutesSinceEpochFailed(#[from] sgp4::DatetimeToMinutesSinceEpochError),
    #[error("SGP4 propagation failed: {0}")]
    PropagationFailed(#[from] sgp4::Error),
    #[error("Failed to convert numeric types")]
    IntConversion(#[from] std::num::TryFromIntError),
}

#[derive(Debug, Error)]
pub enum TimeRangeError {
    #[error("Start must be <= end")]
    InvalidOrder,
}

#[derive(Debug, Error)]
pub enum SamplingError {
    #[error("Step must be > 0")]
    InvalidStep,
}
