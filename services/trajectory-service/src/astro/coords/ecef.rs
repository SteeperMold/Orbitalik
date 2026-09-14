use std::ops::Sub;
use uom::si::angle::radian;
use uom::si::f64::Length;
use uom::si::length::kilometer;

use crate::astro::consts::{A, E2};
use crate::astro::coords::geodetic::Geodetic;

/// Earth-Centered Earth-Fixed coordinates
#[derive(Clone, Copy)]
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

#[cfg(test)]
mod tests {
    use super::*;
    use crate::astro::test_utils::{assert_close, test_observer};

    use uom::si::angle::degree;
    use uom::si::f64::{Angle, Length};
    use uom::si::length::kilometer;

    use crate::astro::coords::geodetic::Geodetic;

    fn geodetic(lat_deg: f64, lon_deg: f64, alt_km: f64) -> Geodetic {
        Geodetic {
            lat: Angle::new::<degree>(lat_deg),
            lon: Angle::new::<degree>(lon_deg),
            alt: Length::new::<kilometer>(alt_km),
        }
    }

    #[test]
    fn equator_and_prime_meridian() {
        let geodetic = geodetic(0.0, 0.0, 0.0);

        let ecef = Ecef::from(&geodetic);

        assert_close(ecef.x.get::<kilometer>(), A, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), 0.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 0.0, 1e-12);
    }

    #[test]
    fn equator_90_degrees_east() {
        let geodetic = geodetic(0.0, 90.0, 0.0);

        let ecef = Ecef::from(&geodetic);

        assert_close(ecef.x.get::<kilometer>(), 0.0, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), A, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 0.0, 1e-12);
    }

    #[test]
    fn equator_180_degrees() {
        let geodetic = geodetic(0.0, 180.0, 0.0);

        let ecef = Ecef::from(&geodetic);

        assert_close(ecef.x.get::<kilometer>(), -A, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), 0.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 0.0, 1e-12);
    }

    #[test]
    fn equator_with_altitude() {
        let altitude = 1000.0;
        let geodetic = geodetic(0.0, 0.0, altitude);

        let ecef = Ecef::from(&geodetic);

        assert_close(ecef.x.get::<kilometer>(), A + altitude, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), 0.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 0.0, 1e-12);
    }

    #[test]
    fn north_pole() {
        let geodetic = geodetic(90.0, 0.0, 0.0);

        let ecef = Ecef::from(&geodetic);

        let expected_z = A * (1.0 - E2).sqrt();

        assert_close(ecef.x.get::<kilometer>(), 0.0, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), 0.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), expected_z, 1e-12);
    }

    #[test]
    fn southern_hemisphere() {
        let geodetic = geodetic(-45.0, 30.0, 0.0);

        let ecef = Ecef::from(&geodetic);

        assert!(ecef.x.get::<kilometer>() > 0.0);
        assert!(ecef.y.get::<kilometer>() > 0.0);
        assert!(ecef.z.get::<kilometer>() < 0.0);
    }

    #[test]
    fn known_wgs84_point() {
        let geodetic = test_observer();

        let ecef = Ecef::from(&geodetic);

        assert_close(ecef.x.get::<kilometer>(), 3_182.645_531_427_102, 1e-6);
        assert_close(ecef.y.get::<kilometer>(), 1_424.012_805_324_485, 1e-6);
        assert_close(ecef.z.get::<kilometer>(), 5_322.841_235_217_519, 1e-6);
    }

    #[test]
    fn geodetic_ecef_round_trip() {
        let original = test_observer();

        let ecef = Ecef::from(&original);
        let result = Geodetic::from(&ecef);

        assert_close(
            result.lat.get::<degree>(),
            original.lat.get::<degree>(),
            1e-10,
        );

        assert_close(
            result.lon.get::<degree>(),
            original.lon.get::<degree>(),
            1e-10,
        );

        assert_close(
            result.alt.get::<kilometer>(),
            original.alt.get::<kilometer>(),
            1e-9,
        );
    }

    #[allow(clippy::float_cmp)]
    #[test]
    fn subtracts_ecef_coordinates() {
        let a = Ecef {
            x: Length::new::<kilometer>(10.0),
            y: Length::new::<kilometer>(20.0),
            z: Length::new::<kilometer>(30.0),
        };

        let b = Ecef {
            x: Length::new::<kilometer>(1.0),
            y: Length::new::<kilometer>(2.0),
            z: Length::new::<kilometer>(3.0),
        };

        let result = a - b;

        assert_eq!(result.x.get::<kilometer>(), 9.0);
        assert_eq!(result.y.get::<kilometer>(), 18.0);
        assert_eq!(result.z.get::<kilometer>(), 27.0);
    }
}
