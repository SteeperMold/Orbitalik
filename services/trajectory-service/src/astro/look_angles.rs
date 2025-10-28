use chrono::{DateTime, Utc};
use uom::si::angle::radian;
use uom::si::f64::{Angle, Length};
use uom::si::length::kilometer;

use crate::astro;
use crate::astro::consts::TWO_PI;
use crate::astro::coords::ecef::Ecef;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::LookAngles;
use crate::astro::propagator::Propagator;
use crate::domain::errors::PropagationError;

impl Propagator {
    pub fn look_angles_at(
        &self,
        datetime: DateTime<Utc>,
        observer: &Geodetic,
    ) -> Result<LookAngles, PropagationError> {
        let eci = self.eci_at(datetime)?;
        let gst = astro::time::utc_to_gst(datetime);

        let sat_ecef = eci.to_ecef(gst);
        let obs_ecef = Ecef::from(observer);

        let rho = sat_ecef - obs_ecef;

        let sin_lat = observer.lat.sin();
        let cos_lat = observer.lat.cos();
        let sin_lon = observer.lon.sin();
        let cos_lon = observer.lon.cos();

        let s = -sin_lat * cos_lon * rho.x - sin_lat * sin_lon * rho.y + cos_lat * rho.z;
        let e = -sin_lon * rho.x + cos_lon * rho.y;
        let z = cos_lat * cos_lon * rho.x + cos_lat * sin_lon * rho.y + sin_lat * rho.z;

        let s_km = s.get::<kilometer>();
        let e_km = e.get::<kilometer>();
        let z_km = z.get::<kilometer>();

        let mut range_km = (s_km * s_km + e_km * e_km + z_km * z_km).sqrt();
        if range_km == 0.0 {
            // prevent division by zero (practically impossible)
            range_km = f64::EPSILON;
        }

        let az_rad = e_km.atan2(s_km).rem_euclid(TWO_PI);
        let el_rad = (z_km / range_km).asin();

        Ok(LookAngles {
            azimuth: Angle::new::<radian>(az_rad),
            elevation: Angle::new::<radian>(el_rad),
            range: Length::new::<kilometer>(range_km),
        })
    }
}
