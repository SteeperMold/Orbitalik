use std::fmt;
use std::fmt::Formatter;
use chrono::{DateTime, Utc};
use uom::si::f64::{Angle, Length};

use crate::astro::coords::{ecef::Ecef, eci::Eci, geodetic::Geodetic};
use crate::astro::errors::{SamplingError, TimeRangeError};

pub struct Tle {
    pub norad_id: u32,
    pub satellite_name: String,
    pub line1: String,
    pub line2: String,
    pub epoch: DateTime<Utc>,
}

pub struct SatellitePosition {
    pub time: DateTime<Utc>,

    pub eci: Option<Eci>,
    pub ecef: Option<Ecef>,
    pub geodetic: Option<Geodetic>,
}

pub struct LookAngles {
    pub time: DateTime<Utc>,

    pub azimuth: Option<Angle>,
    pub elevation: Option<Angle>,
    pub range: Option<Length>,
}

pub struct Trajectory {
    pub start: DateTime<Utc>,
    pub end: DateTime<Utc>,

    pub step_seconds: u32,

    pub samples: Vec<SatellitePosition>,
}

pub struct ObserverTrajectory {
    pub start: DateTime<Utc>,
    pub end: DateTime<Utc>,

    pub step_seconds: u32,

    pub samples: Vec<LookAngles>,
}

pub struct Pass {
    pub satellite: SatelliteIdentifier,

    pub aos: DateTime<Utc>,
    pub aos_azimuth: f64,

    pub max_elevation_time: DateTime<Utc>,
    pub max_elevation: f64,
    pub max_elevation_azimuth: f64,

    pub los: DateTime<Utc>,
    pub los_azimuth: f64,

    pub duration_seconds: u32,
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