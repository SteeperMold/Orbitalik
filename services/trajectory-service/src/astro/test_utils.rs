use chrono::{DateTime, Utc};
use uom::si::angle::degree;
use uom::si::f64::{Angle, Length};
use uom::si::length::kilometer;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::Tle;
use crate::astro::propagation::propagator::Propagator;

#[allow(clippy::unwrap_used)]
pub fn test_tle() -> Tle {
    Tle {
        norad_id: 25544,
        satellite_name: "ISS (ZARYA)".to_string(),
        line1: "1 25544U 98067A   26232.17880947  .00009753  00000+0  18154-3 0  9994".to_string(),
        line2: "2 25544  51.6332 343.3775 0007674  65.0551 295.1235 15.49524101581676".to_string(),
        epoch: "2026-01-01T00:00:00Z".parse().unwrap(),
    }
}

#[allow(clippy::expect_used)]
pub fn test_propagator() -> Propagator {
    Propagator::from_tle(&test_tle()).expect("test TLE should produce a valid propagator")
}

pub fn test_observer() -> Geodetic {
    Geodetic {
        // Riga
        lat: Angle::new::<degree>(56.9496),
        lon: Angle::new::<degree>(24.1052),
        alt: Length::new::<kilometer>(0.0),
    }
}

#[allow(clippy::unwrap_used)]
pub fn test_datetime() -> DateTime<Utc> {
    "2026-01-01T00:00:00Z".parse().unwrap()
}

pub fn assert_close(actual: f64, expected: f64, tolerance: f64) {
    let error = (actual - expected).abs();

    assert!(
        error <= tolerance,
        "expected {expected}, got {actual}, error {error} > tolerance {tolerance}"
    );
}
