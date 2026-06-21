use chrono::{DateTime, Utc};
use std::fmt;
use std::fmt::Formatter;

use crate::domain::errors::{SamplingError, TimeRangeError};

pub struct ComputationMetadata {
    pub propagation_model: String,
    pub computation_time: DateTime<Utc>,
    pub norad_id: u32,
    pub satellite_name: String,
    pub tle_epoch: DateTime<Utc>,
}

#[derive(Debug, Clone)]
pub enum SatelliteIdentifier {
    NoradId(u32),
    Name(String),
}

impl fmt::Display for SatelliteIdentifier {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match self {
            Self::NoradId(id) => write!(f, "NORAD ID {id}"),
            Self::Name(name) => write!(f, "Satellite '{name}'"),
        }
    }
}

#[derive(Debug, Clone, Copy)]
pub struct TimeRange {
    pub start: DateTime<Utc>,
    pub end: DateTime<Utc>,
}

impl TimeRange {
    pub fn new(start: DateTime<Utc>, end: DateTime<Utc>) -> Result<Self, TimeRangeError> {
        if start > end {
            return Err(TimeRangeError::InvalidOrder);
        }

        Ok(Self { start, end })
    }
}

#[derive(Debug, Clone, Copy)]
pub struct Sampling {
    pub step_seconds: u32,
}

impl Sampling {
    pub const fn new(step_seconds: u32) -> Result<Self, SamplingError> {
        if step_seconds == 0 {
            return Err(SamplingError::InvalidStep);
        }

        Ok(Self { step_seconds })
    }
}
