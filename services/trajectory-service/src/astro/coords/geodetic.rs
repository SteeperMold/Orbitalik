use uom::si::angle::radian;
use uom::si::f64::{Angle, Length};
use uom::si::length::kilometer;

use crate::astro::consts::{A, E2};
use crate::astro::coords::ecef::Ecef;

/// Geodetic coordinates
pub struct Geodetic {
    pub lat: Angle,
    pub lon: Angle,
    pub alt: Length,
}

impl From<&Ecef> for Geodetic {
    fn from(ecef: &Ecef) -> Self {
        let x_km = ecef.x.get::<kilometer>();
        let y_km = ecef.y.get::<kilometer>();
        let z_km = ecef.z.get::<kilometer>();

        let lon = Angle::new::<radian>(y_km.atan2(x_km));

        let r = x_km.hypot(y_km);
        let mut lat = z_km.atan2(r * (1.0 - E2));
        let mut h;

        loop {
            let sin_lat = lat.sin();
            let cos_lat = lat.cos();

            let n = A / (E2 * sin_lat).mul_add(-sin_lat, 1.0).sqrt();
            h = r / cos_lat - n;

            let new_lat = (E2 * n).mul_add(sin_lat, z_km).atan2(r);

            if (new_lat - lat).abs() < 1e-12 {
                lat = new_lat;
                break;
            }

            lat = new_lat;
        }

        Self {
            lat: Angle::new::<radian>(lat),
            lon,
            alt: Length::new::<kilometer>(h),
        }
    }
}
