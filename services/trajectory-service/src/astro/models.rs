use chrono::{DateTime, Utc};
use std::fmt;
use std::fmt::Formatter;
use uom::si::f64::{Angle, Length};

use crate::astro::coords::{ecef::Ecef, geodetic::Geodetic, teme::Teme};
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

    pub teme: Option<Teme>,
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
    pub aos_azimuth: Angle,

    pub max_elevation_time: DateTime<Utc>,
    pub max_elevation: Angle,
    pub max_elevation_azimuth: Angle,

    pub los: DateTime<Utc>,
    pub los_azimuth: Angle,

    pub duration_seconds: u32,
}

#[derive(Debug, Clone, PartialEq, Eq)]
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

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::{DateTime, Utc};

    #[allow(clippy::unwrap_used)]
    fn dt(value: &str) -> DateTime<Utc> {
        value.parse().unwrap()
    }

    #[allow(clippy::unwrap_used)]
    #[test]
    fn time_range_accepts_start_before_end() {
        let start = dt("2026-01-01T00:00:00Z");
        let end = dt("2026-01-02T00:00:00Z");

        let range = TimeRange::new(start, end).unwrap();

        assert_eq!(range.start, start);
        assert_eq!(range.end, end);
    }

    #[allow(clippy::unwrap_used)]
    #[test]
    fn time_range_accepts_equal_start_and_end() {
        let time = dt("2026-01-01T00:00:00Z");

        let range = TimeRange::new(time, time).unwrap();

        assert_eq!(range.start, time);
        assert_eq!(range.end, time);
    }

    #[test]
    fn time_range_rejects_start_after_end() {
        let start = dt("2026-01-02T00:00:00Z");
        let end = dt("2026-01-01T00:00:00Z");

        let result = TimeRange::new(start, end);

        assert!(matches!(result, Err(TimeRangeError::InvalidOrder)));
    }

    #[allow(clippy::unwrap_used)]
    #[test]
    fn sampling_accepts_positive_step() {
        let sampling = Sampling::new(60).unwrap();

        assert_eq!(sampling.step_seconds, 60);
    }

    #[test]
    fn sampling_rejects_zero_step() {
        let result = Sampling::new(0);

        assert!(matches!(result, Err(SamplingError::InvalidStep)));
    }
}
