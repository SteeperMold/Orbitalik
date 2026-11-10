use chrono::{DateTime, Utc};

use crate::astro;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::errors::PropagationError;
use crate::astro::models::SatellitePosition;
use crate::astro::propagation::propagator::Propagator;

#[derive(Default, PartialEq, Eq)]
pub struct PositionComputation {
    pub teme: bool,
    pub ecef: bool,
    pub geodetic: bool,
}

impl Propagator {
    pub fn position_at(
        &self,
        datetime: DateTime<Utc>,
        compute: &PositionComputation,
    ) -> Result<SatellitePosition, PropagationError> {
        // compute teme only if any dependent coordinate is needed
        let teme = if compute.teme || compute.ecef || compute.geodetic {
            Some(self.teme_at(datetime)?)
        } else {
            None
        };

        // only if ecef is requested and teme is available
        let ecef = match (compute.ecef || compute.geodetic, &teme) {
            (true, Some(teme_val)) => {
                let gst = astro::time::utc_to_gst(datetime);
                Some(teme_val.to_ecef(gst))
            }
            _ => None,
        };

        // only if geodetic is requested and ecef is available
        let geodetic = match (compute.geodetic, &ecef) {
            (true, Some(ecef_val)) => Some(Geodetic::from(ecef_val)),
            _ => None,
        };

        Ok(SatellitePosition {
            time: datetime,

            teme: compute.teme.then_some(teme).flatten(),
            ecef: compute.ecef.then_some(ecef).flatten(),
            geodetic,
        })
    }
}

#[allow(clippy::unwrap_used)]
#[cfg(test)]
mod tests {
    use super::*;

    use crate::astro::test_utils::{test_datetime, test_propagator};

    #[test]
    fn no_computation_returns_empty_result() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation::default();

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.teme.is_none());
        assert!(result.ecef.is_none());
        assert!(result.geodetic.is_none());
    }

    #[test]
    fn teme_only_returns_teme() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation {
            teme: true,
            ecef: false,
            geodetic: false,
        };

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.teme.is_some());
        assert!(result.ecef.is_none());
        assert!(result.geodetic.is_none());
    }

    #[test]
    fn ecef_only_returns_ecef() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation {
            teme: false,
            ecef: true,
            geodetic: false,
        };

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.teme.is_none());
        assert!(result.ecef.is_some());
        assert!(result.geodetic.is_none());
    }

    #[test]
    fn geodetic_only_returns_geodetic() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation {
            teme: false,
            ecef: false,
            geodetic: true,
        };

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.teme.is_none());
        assert!(result.ecef.is_none());
        assert!(result.geodetic.is_some());
    }

    #[test]
    fn all_coordinates_returns_all_values() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation {
            teme: true,
            ecef: true,
            geodetic: true,
        };

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert_eq!(result.time, datetime);
        assert!(result.teme.is_some());
        assert!(result.ecef.is_some());
        assert!(result.geodetic.is_some());
    }

    #[test]
    fn ecef_computation_includes_internal_teme() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation {
            teme: false,
            ecef: true,
            geodetic: false,
        };

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert!(result.ecef.is_some());
        assert!(result.teme.is_none());
    }

    #[test]
    fn geodetic_computation_includes_internal_ecef() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let compute = PositionComputation {
            teme: false,
            ecef: false,
            geodetic: true,
        };

        let result = propagator.position_at(datetime, &compute).unwrap();

        assert!(result.geodetic.is_some());
        assert!(result.ecef.is_none());
        assert!(result.teme.is_none());
    }

    // TODO return this test
    // #[test]
    // fn position_matches_reference() {
    //     let propagator = test_propagator();
    //     let datetime = test_datetime();
    //
    //     let compute = PositionComputation {
    //         teme: true,
    //         ecef: true,
    //         geodetic: true,
    //     };
    //
    //     let result = propagator.position_at(datetime, &compute).unwrap();
    //
    //     let teme = result.teme.unwrap();
    //     let ecef = result.ecef.unwrap();
    //     let geodetic = result.geodetic.unwrap();
    //
    //     assert_close(teme.x.get::<kilometer>(), -672.405_384_347_387, 1e-6);
    //     assert_close(teme.y.get::<kilometer>(), -5_319.625_479_765_174, 1e-6);
    //     assert_close(teme.z.get::<kilometer>(), -4_199.191_899_815_303, 1e-6);
    //
    //     assert_close(ecef.x.get::<kilometer>(), -5_103.404_856_315_442, 1e-6);
    //     assert_close(ecef.y.get::<kilometer>(), 1_644.932_557_402_271_4, 1e-6);
    //     assert_close(ecef.z.get::<kilometer>(), -4_199.191_899_815_276, 1e-6);
    //
    //     assert_close(geodetic.lat.get::<degree>(), -38.240_979_007_316_32, 1e-6);
    //     assert_close(geodetic.lon.get::<degree>(), 162.134_799_204_196_2, 1e-6);
    //     assert_close(
    //         geodetic.alt.get::<kilometer>(),
    //         440.578_369_455_746_04,
    //         1e-6,
    //     );
    // }
}
