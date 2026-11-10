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

        if r < 1e-12 {
            let lat = if z_km >= 0.0 {
                std::f64::consts::FRAC_PI_2
            } else {
                -std::f64::consts::FRAC_PI_2
            };

            let polar_radius = A * (1.0 - E2).sqrt();

            let h = z_km.abs() - polar_radius;

            return Self {
                lat: Angle::new::<radian>(lat),
                lon: Angle::new::<radian>(0.0),
                alt: Length::new::<kilometer>(h),
            };
        }

        let mut lat = z_km.atan2(r * (1.0 - E2));
        let mut h;

        loop {
            let sin_lat = lat.sin();
            let cos_lat = lat.cos();

            let n = A / (1.0 - E2 * sin_lat * sin_lat).sqrt();

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

#[cfg(test)]
mod tests {
    use super::*;
    use crate::astro::test_utils::assert_close;

    use uom::si::angle::degree;
    use uom::si::length::kilometer;

    use crate::astro::coords::ecef::Ecef;

    #[test]
    fn converts_45_north_45_east() {
        let ecef = Ecef {
            x: Length::new::<kilometer>(3_194.419_145_061),
            y: Length::new::<kilometer>(3_194.419_145_061),
            z: Length::new::<kilometer>(4_487.348_408_866),
        };

        let geodetic = Geodetic::from(&ecef);

        assert_close(geodetic.lat.get::<degree>(), 45.0, 1e-9);
        assert_close(geodetic.lon.get::<degree>(), 45.0, 1e-9);
        assert_close(geodetic.alt.get::<kilometer>(), 0.0, 1e-6);
    }

    #[test]
    fn converts_equator_prime_meridian() {
        let ecef = Ecef {
            x: Length::new::<kilometer>(A),
            y: Length::new::<kilometer>(0.0),
            z: Length::new::<kilometer>(0.0),
        };

        let geodetic = Geodetic::from(&ecef);

        assert_close(geodetic.lat.get::<degree>(), 0.0, 1e-10);
        assert_close(geodetic.lon.get::<degree>(), 0.0, 1e-10);
        assert_close(geodetic.alt.get::<kilometer>(), 0.0, 1e-9);
    }

    #[test]
    fn converts_equator_90_east() {
        let ecef = Ecef {
            x: Length::new::<kilometer>(0.0),
            y: Length::new::<kilometer>(A),
            z: Length::new::<kilometer>(0.0),
        };

        let geodetic = Geodetic::from(&ecef);

        assert_close(geodetic.lat.get::<degree>(), 0.0, 1e-10);
        assert_close(geodetic.lon.get::<degree>(), 90.0, 1e-10);
        assert_close(geodetic.alt.get::<kilometer>(), 0.0, 1e-9);
    }

    #[test]
    fn converts_north_pole() {
        let z = A * (1.0 - E2).sqrt();

        let ecef = Ecef {
            x: Length::new::<kilometer>(0.0),
            y: Length::new::<kilometer>(0.0),
            z: Length::new::<kilometer>(z),
        };

        let geodetic = Geodetic::from(&ecef);

        assert_close(geodetic.lat.get::<degree>(), 90.0, 1e-9);
        assert_close(geodetic.alt.get::<kilometer>(), 0.0, 1e-6);
    }

    #[test]
    fn geodetic_ecef_round_trip() {
        let original = Geodetic {
            lat: Angle::new::<degree>(56.9496),
            lon: Angle::new::<degree>(24.1052),
            alt: Length::new::<kilometer>(0.123),
        };

        let ecef = Ecef::from(&original);
        let result = Geodetic::from(&ecef);

        assert_close(
            result.lat.get::<degree>(),
            original.lat.get::<degree>(),
            1e-8,
        );

        assert_close(
            result.lon.get::<degree>(),
            original.lon.get::<degree>(),
            1e-8,
        );

        assert_close(
            result.alt.get::<kilometer>(),
            original.alt.get::<kilometer>(),
            1e-8,
        );
    }

    #[test]
    fn known_geodetic_point() {
        let ecef = Ecef {
            x: Length::new::<kilometer>(3_182.645_531_427_102),
            y: Length::new::<kilometer>(1_424.012_805_324_485),
            z: Length::new::<kilometer>(5_322.841_235_217_519),
        };

        let geodetic = Geodetic::from(&ecef);

        assert_close(geodetic.lat.get::<degree>(), 56.9496, 1e-10);
        assert_close(geodetic.lon.get::<degree>(), 24.1052, 1e-10);
        assert_close(geodetic.alt.get::<kilometer>(), 0.0, 1e-9);
    }
}
