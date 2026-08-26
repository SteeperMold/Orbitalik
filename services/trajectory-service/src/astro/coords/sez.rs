use uom::si::angle::radian;
use uom::si::f64::{Angle, Length};
use uom::si::length::kilometer;

use crate::astro::consts::TWO_PI;
use crate::astro::coords::ecef::Ecef;
use crate::astro::coords::geodetic::Geodetic;

#[derive(Debug, Clone, Copy)]
pub struct Sez {
    pub south: Length,
    pub east: Length,
    pub zenith: Length,
}

impl Sez {
    pub fn from_ecef(satellite: &Ecef, observer: &Geodetic) -> Self {
        let observer_ecef = Ecef::from(observer);
        let rho = *satellite - observer_ecef;

        let sin_lat = observer.lat.sin();
        let cos_lat = observer.lat.cos();
        let sin_lon = observer.lon.sin();
        let cos_lon = observer.lon.cos();

        let south = -sin_lat * cos_lon * rho.x - sin_lat * sin_lon * rho.y + cos_lat * rho.z;

        let east = -sin_lon * rho.x + cos_lon * rho.y;

        let zenith = cos_lat * cos_lon * rho.x + cos_lat * sin_lon * rho.y + sin_lat * rho.z;

        Self {
            south,
            east,
            zenith,
        }
    }

    pub fn range(&self) -> Length {
        let south_km = self.south.get::<kilometer>();
        let east_km = self.east.get::<kilometer>();
        let zenith_km = self.zenith.get::<kilometer>();

        let range_km = (south_km * south_km + east_km * east_km + zenith_km * zenith_km).sqrt();

        Length::new::<kilometer>(range_km)
    }

    pub fn azimuth(&self) -> Angle {
        let south_km = self.south.get::<kilometer>();
        let east_km = self.east.get::<kilometer>();

        let az_rad = east_km.atan2(south_km).rem_euclid(TWO_PI);

        Angle::new::<radian>(az_rad)
    }

    pub fn elevation(&self) -> Angle {
        let range_km = self.range().get::<kilometer>();

        // prevent division by zero for the case where the
        // satellite and observer occupy exactly the same position
        let range_km = if range_km == 0.0 {
            f64::EPSILON
        } else {
            range_km
        };

        let zenith_km = self.zenith.get::<kilometer>();
        let elevation_rad = (zenith_km / range_km).asin();

        Angle::new::<radian>(elevation_rad)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    use uom::si::angle::degree;
    use uom::si::length::kilometer;

    use crate::astro::coords::ecef::Ecef;
    use crate::astro::coords::geodetic::Geodetic;
    use crate::astro::test_utils::assert_close;

    fn observer() -> Geodetic {
        Geodetic {
            lat: Angle::new::<degree>(0.0),
            lon: Angle::new::<degree>(0.0),
            alt: Length::new::<kilometer>(0.0),
        }
    }

    fn observer_ecef() -> Ecef {
        Ecef::from(&observer())
    }

    fn satellite_offset(south: f64, east: f64, zenith: f64) -> Ecef {
        let observer = observer_ecef();

        Ecef {
            x: observer.x + Length::new::<kilometer>(zenith),
            y: observer.y + Length::new::<kilometer>(east),
            z: observer.z + Length::new::<kilometer>(south),
        }
    }

    #[test]
    fn straight_up() {
        let satellite = satellite_offset(0.0, 0.0, 1000.0);

        let sez = Sez::from_ecef(&satellite, &observer());

        assert_close(0.0, sez.south.get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.east.get::<kilometer>(), 1e-9);
        assert_close(1000.0, sez.zenith.get::<kilometer>(), 1e-9);

        assert_close(1000.0, sez.range().get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.azimuth().get::<degree>(), 1e-9);
        assert_close(90.0, sez.elevation().get::<degree>(), 1e-9);
    }

    #[test]
    fn directly_south() {
        let satellite = satellite_offset(1000.0, 0.0, 0.0);

        let sez = Sez::from_ecef(&satellite, &observer());

        assert_close(1000.0, sez.south.get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.east.get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.zenith.get::<kilometer>(), 1e-9);

        assert_close(1000.0, sez.range().get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.azimuth().get::<degree>(), 1e-9);
        assert_close(0.0, sez.elevation().get::<degree>(), 1e-9);
    }

    #[test]
    fn directly_east() {
        let satellite = satellite_offset(0.0, 1000.0, 0.0);

        let sez = Sez::from_ecef(&satellite, &observer());

        assert_close(0.0, sez.south.get::<kilometer>(), 1e-9);
        assert_close(1000.0, sez.east.get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.zenith.get::<kilometer>(), 1e-9);

        assert_close(1000.0, sez.range().get::<kilometer>(), 1e-9);
        assert_close(90.0, sez.azimuth().get::<degree>(), 1e-9);
        assert_close(0.0, sez.elevation().get::<degree>(), 1e-9);
    }

    #[test]
    fn directly_north() {
        let satellite = satellite_offset(-1000.0, 0.0, 0.0);

        let sez = Sez::from_ecef(&satellite, &observer());

        assert_close(-1000.0, sez.south.get::<kilometer>(), 1e-9);

        assert_close(180.0, sez.azimuth().get::<degree>(), 1e-9);
        assert_close(0.0, sez.elevation().get::<degree>(), 1e-9);
        assert_close(1000.0, sez.range().get::<kilometer>(), 1e-9);
    }

    #[test]
    fn directly_west() {
        let satellite = satellite_offset(0.0, -1000.0, 0.0);

        let sez = Sez::from_ecef(&satellite, &observer());

        assert_close(0.0, sez.south.get::<kilometer>(), 1e-9);
        assert_close(-1000.0, sez.east.get::<kilometer>(), 1e-9);
        assert_close(0.0, sez.zenith.get::<kilometer>(), 1e-9);

        assert_close(1000.0, sez.range().get::<kilometer>(), 1e-9);
        assert_close(270.0, sez.azimuth().get::<degree>(), 1e-9);
        assert_close(0.0, sez.elevation().get::<degree>(), 1e-9);
    }

    #[test]
    fn diagonal() {
        let satellite = satellite_offset(1000.0, 1000.0, 1000.0);

        let sez = Sez::from_ecef(&satellite, &observer());

        let expected_range = 3.0_f64.sqrt() * 1000.0;

        assert_close(1000.0, sez.south.get::<kilometer>(), 1e-9);
        assert_close(1000.0, sez.east.get::<kilometer>(), 1e-9);
        assert_close(1000.0, sez.zenith.get::<kilometer>(), 1e-9);

        assert_close(expected_range, sez.range().get::<kilometer>(), 1e-9);
        assert_close(45.0, sez.azimuth().get::<degree>(), 1e-9);
        assert_close(
            35.264_389_682_754_654,
            sez.elevation().get::<degree>(),
            1e-9,
        );
    }
}
