use chrono::{DateTime, Utc};
use std::fmt;
use std::fmt::Formatter;

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

pub struct ComputationMetadata {
    pub propagation_model: String,
    pub computation_time: DateTime<Utc>,
    pub norad_id: u32,
    pub satellite_name: String,
    pub tle_epoch: DateTime<Utc>,
}

