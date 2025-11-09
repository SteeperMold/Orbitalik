use chrono::{DateTime, Utc};

use crate::astro;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::SatellitePosition;
use crate::astro::propagator::Propagator;
use crate::domain::errors::PropagationError;

#[derive(Default)]
pub struct PositionComputation {
    pub eci: bool,
    pub ecef: bool,
    pub geodetic: bool,
}

impl Propagator {
    pub fn position_at(
        &self,
        datetime: DateTime<Utc>,
        compute: &PositionComputation,
    ) -> Result<SatellitePosition, PropagationError> {
        // compute eci only if any dependent coordinate is needed
        let eci = if compute.eci || compute.ecef || compute.geodetic {
            Some(self.eci_at(datetime)?)
        } else {
            None
        };

        // only if ecef is requested and eci is available
        let ecef = match (compute.ecef || compute.geodetic, &eci) {
            (true, Some(eci_val)) => {
                let gst = astro::time::utc_to_gst(datetime);
                Some(eci_val.to_ecef(gst))
            }
            _ => None,
        };

        // only if geodetic is requested and ecef is available
        let geodetic = match (compute.geodetic, &ecef) {
            (true, Some(ecef_val)) => Some(Geodetic::from(ecef_val)),
            _ => None,
        };

        Ok(SatellitePosition {
            eci: compute.eci.then_some(eci).flatten(),
            ecef: compute.ecef.then_some(ecef).flatten(),
            geodetic,
        })
    }
}
