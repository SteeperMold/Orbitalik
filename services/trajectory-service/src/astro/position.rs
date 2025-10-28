use chrono::{DateTime, Utc};

use crate::astro;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::SatellitePosition;
use crate::astro::propagator::Propagator;
use crate::domain::errors::PropagationError;

impl Propagator {
    pub fn position_at(
        &self,
        datetime: DateTime<Utc>,
    ) -> Result<SatellitePosition, PropagationError> {
        let eci = self.eci_at(datetime)?;

        let gst = astro::time::utc_to_gst(datetime);
        let ecef = eci.to_ecef(gst);

        let geodetic = Geodetic::from(&ecef);

        Ok(SatellitePosition {
            eci,
            ecef,
            geodetic,
        })
    }
}
