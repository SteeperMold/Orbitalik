use std::ops::Sub;
use uom::si::angle::radian;
use uom::si::f64::Length;
use uom::si::length::kilometer;

use crate::astro::consts::{A, E2};
use crate::astro::coords::geodetic::Geodetic;

/// Earth-Centered Earth-Fixed coordinates
pub struct Ecef {
    pub x: Length,
    pub y: Length,
    pub z: Length,
}

impl From<&Geodetic> for Ecef {
    fn from(geodetic: &Geodetic) -> Self {
        let alt_km = geodetic.alt.get::<kilometer>();

        let sin_lat = geodetic.lat.get::<radian>().sin();
        let cos_lat = geodetic.lat.get::<radian>().cos();
        let sin_lon = geodetic.lon.get::<radian>().sin();
        let cos_lon = geodetic.lon.get::<radian>().cos();

        let n = A / (1.0 - E2 * sin_lat * sin_lat).sqrt();

        let x_km = (n + alt_km) * cos_lat * cos_lon;
        let y_km = (n + alt_km) * cos_lat * sin_lon;
        let z_km = (n * (1.0 - E2) + alt_km) * sin_lat;

        Self {
            x: Length::new::<kilometer>(x_km),
            y: Length::new::<kilometer>(y_km),
            z: Length::new::<kilometer>(z_km),
        }
    }
}

impl Sub for Ecef {
    type Output = Self;

    fn sub(self, rhs: Self) -> Self::Output {
        Self {
            x: self.x - rhs.x,
            y: self.y - rhs.y,
            z: self.z - rhs.z,
        }
    }
}
