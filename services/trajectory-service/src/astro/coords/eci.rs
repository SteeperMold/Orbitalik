use uom::si::f64::{Angle, Length};
use uom::si::length::kilometer;

use crate::astro::coords::ecef::Ecef;

/// Earth-Centered Inertial coordinates
pub struct Eci {
    pub x: Length,
    pub y: Length,
    pub z: Length,
}

impl From<[f64; 3]> for Eci {
    fn from(v: [f64; 3]) -> Self {
        Self {
            x: Length::new::<kilometer>(v[0]),
            y: Length::new::<kilometer>(v[1]),
            z: Length::new::<kilometer>(v[2]),
        }
    }
}

impl Eci {
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
