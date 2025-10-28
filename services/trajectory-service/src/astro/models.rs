use chrono::{DateTime, Utc};
use uom::si::f64::{Angle, Length};

use crate::astro::coords::{ecef::Ecef, eci::Eci, geodetic::Geodetic};

pub struct Tle {
    pub _norad_id: u32,
    pub satellite_name: String,
    pub line1: String,
    pub line2: String,
    pub _epoch: DateTime<Utc>,
}

pub struct SatellitePosition {
    pub eci: Eci,
    pub ecef: Ecef,
    pub geodetic: Geodetic,
}

pub struct LookAngles {
    pub azimuth: Angle,
    pub elevation: Angle,
    pub range: Length,
}
