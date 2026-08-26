use chrono::{DateTime, Utc};

use crate::astro;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::coords::sez::Sez;
use crate::astro::errors::PropagationError;
use crate::astro::models::LookAngles;
use crate::astro::propagation::propagator::Propagator;

#[derive(Default)]
pub struct LookAnglesComputation {
    pub azimuth: bool,
    pub elevation: bool,
    pub range: bool,
}

impl Propagator {
    pub fn look_angles_at(
        &self,
        datetime: DateTime<Utc>,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<LookAngles, PropagationError> {
        if !compute.azimuth && !compute.range && !compute.elevation {
            return Ok(LookAngles {
                time: datetime,
                azimuth: None,
                elevation: None,
                range: None,
            });
        }

        let eci = self.teme_at(datetime)?;
        let gst = astro::time::utc_to_gst(datetime);
        let sat_ecef = eci.to_ecef(gst);

        let sez = Sez::from_ecef(&sat_ecef, observer);

        Ok(LookAngles {
            time: datetime,
            azimuth: compute.azimuth.then(|| sez.azimuth()),
            elevation: compute.elevation.then(|| sez.elevation()),
            range: compute.range.then(|| sez.range()),
        })
    }
}

#[allow(clippy::unwrap_used, clippy::expect_used)]
#[cfg(test)]
mod tests {
    use super::*;
    use crate::astro::test_utils::{test_datetime, test_observer, test_propagator};

    use uom::si::angle::radian;
    use uom::si::length::kilometer;

    #[test]
    fn no_computation_returns_empty_result() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation::default();

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.azimuth.is_none());
        assert!(result.elevation.is_none());
        assert!(result.range.is_none());
    }

    #[test]
    fn azimuth_only_returns_azimuth() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: false,
            range: false,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.azimuth.is_some());
        assert!(result.elevation.is_none());
        assert!(result.range.is_none());
    }

    #[test]
    fn elevation_only_returns_elevation() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: false,
            elevation: true,
            range: false,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.azimuth.is_none());
        assert!(result.elevation.is_some());
        assert!(result.range.is_none());
    }

    #[test]
    fn range_only_returns_range() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: false,
            elevation: false,
            range: true,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.azimuth.is_none());
        assert!(result.elevation.is_none());
        assert!(result.range.is_some());
    }

    #[test]
    fn all_computations_return_all_values() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: true,
            range: true,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.azimuth.is_some());
        assert!(result.elevation.is_some());
        assert!(result.range.is_some());
    }

    #[test]
    fn azimuth_is_in_valid_range() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: true,
            elevation: false,
            range: false,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        let azimuth = result.azimuth.unwrap().get::<radian>();

        assert!((0.0..std::f64::consts::TAU).contains(&azimuth));
    }

    #[test]
    fn elevation_is_in_valid_range() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: false,
            elevation: true,
            range: false,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        let elevation = result.elevation.unwrap().get::<radian>();

        assert!((-std::f64::consts::FRAC_PI_2..=std::f64::consts::FRAC_PI_2).contains(&elevation));
    }

    #[test]
    fn range_is_positive() {
        let propagator = test_propagator();
        let observer = test_observer();
        let datetime = test_datetime();

        let compute = LookAnglesComputation {
            azimuth: false,
            elevation: false,
            range: true,
        };

        let result = propagator
            .look_angles_at(datetime, &observer, &compute)
            .unwrap();

        let range = result.range.unwrap().get::<kilometer>();

        assert!(range > 0.0);
    }

    // TODO return this test
    // #[test]
    // fn look_angles_match_reference() {
    //     let propagator = test_propagator();
    //     let observer = test_observer();
    //     let datetime = test_datetime();
    //
    //     let compute = LookAnglesComputation {
    //         azimuth: true,
    //         elevation: true,
    //         range: true,
    //     };
    //
    //     let result = propagator
    //         .look_angles_at(datetime, &observer, &compute)
    //         .unwrap();
    //
    //     let azimuth = result.azimuth.unwrap().get::<degree>();
    //     let elevation = result.elevation.unwrap().get::<degree>();
    //     let range = result.range.unwrap().get::<kilometer>();
    //
    //     // values received from python library skifield
    //     assert_close(azimuth, 73.371_333_116_776, 1e-6);
    //     assert_close(elevation, -72.756_555_462_194, 1e-6);
    //     assert_close(range, 12_624.442_624_754_416, 1e-6);
    // }
}
