use chrono::{DateTime, Utc};
use uom::si::f64::{Angle, Length};

use crate::astro::coords::{ecef::Ecef, eci::Eci, geodetic::Geodetic};

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
