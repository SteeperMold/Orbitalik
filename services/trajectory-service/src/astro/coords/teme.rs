use uom::si::f64::{Angle, Length};
use uom::si::length::kilometer;

use crate::astro::coords::ecef::Ecef;

/// True Equator, Mean Equinox coordinates
#[derive(Clone)]
pub struct Teme {
    pub x: Length,
    pub y: Length,
    pub z: Length,
}

impl From<[f64; 3]> for Teme {
    fn from(v: [f64; 3]) -> Self {
        Self {
            x: Length::new::<kilometer>(v[0]),
            y: Length::new::<kilometer>(v[1]),
            z: Length::new::<kilometer>(v[2]),
        }
    }
}

impl Teme {
    pub fn to_ecef(&self, gst: Angle) -> Ecef {
        let sin_gst = gst.sin();
        let cos_gst = gst.cos();

        let x_ecef = cos_gst.mul_add(self.x, sin_gst * self.y);
        let y_ecef = (-sin_gst).mul_add(self.x, cos_gst * self.y);
        let z_ecef = self.z; // z axis unchanged

        Ecef {
            x: x_ecef,
            y: y_ecef,
            z: z_ecef,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::astro::test_utils::assert_close;

    use uom::si::angle::degree;

    #[test]
    fn teme_to_ecef_with_zero_gst() {
        let teme = Teme::from([1000.0, 2000.0, 3000.0]);

        let gst = Angle::new::<degree>(0.0);

        let ecef = teme.to_ecef(gst);

        assert_close(ecef.x.get::<kilometer>(), 1000.0, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), 2000.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 3000.0, 1e-12);
    }

    #[test]
    fn teme_to_ecef_with_90_degree_gst() {
        let teme = Teme::from([1000.0, 2000.0, 3000.0]);

        let gst = Angle::new::<degree>(90.0);

        let ecef = teme.to_ecef(gst);

        assert_close(ecef.x.get::<kilometer>(), 2000.0, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), -1000.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 3000.0, 1e-12);
    }

    #[test]
    fn teme_to_ecef_with_180_degree_gst() {
        let teme = Teme::from([1000.0, 2000.0, 3000.0]);

        let gst = Angle::new::<degree>(180.0);

        let ecef = teme.to_ecef(gst);

        assert_close(ecef.x.get::<kilometer>(), -1000.0, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), -2000.0, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 3000.0, 1e-12);
    }

    #[test]
    fn teme_to_ecef_with_45_degree_gst() {
        let teme = Teme::from([1000.0, 0.0, 3000.0]);
        let gst = Angle::new::<degree>(45.0);
        let ecef = teme.to_ecef(gst);

        let expected = 1000.0 / 2.0_f64.sqrt();

        assert_close(ecef.x.get::<kilometer>(), expected, 1e-12);
        assert_close(ecef.y.get::<kilometer>(), -expected, 1e-12);
        assert_close(ecef.z.get::<kilometer>(), 3000.0, 1e-12);
    }

    #[test]
    fn teme_to_ecef_preserves_distance_from_origin() {
        let teme = Teme::from([7000.0, 2000.0, -1000.0]);

        let gst = Angle::new::<degree>(37.5);

        let ecef = teme.to_ecef(gst);

        let teme_r = (teme.x.get::<kilometer>().powi(2)
            + teme.y.get::<kilometer>().powi(2)
            + teme.z.get::<kilometer>().powi(2))
        .sqrt();

        let ecef_r = (ecef.x.get::<kilometer>().powi(2)
            + ecef.y.get::<kilometer>().powi(2)
            + ecef.z.get::<kilometer>().powi(2))
        .sqrt();

        assert_close(ecef_r, teme_r, 1e-12);
    }
}
