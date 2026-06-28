use chrono::{DateTime, Utc};

use crate::astro::coords::eci::Eci;
use crate::astro::errors::PropagationError;
use crate::astro::models::Tle;

pub struct Propagator {
    pub _satellite_name: String,
    pub norad_id: u32,

    elements: sgp4::Elements,
    constants: sgp4::Constants,
}

impl Propagator {
    pub fn from_tle(tle: &Tle) -> Result<Self, PropagationError> {
        let elements = sgp4::Elements::from_tle(
            Some(tle.satellite_name.clone()),
            tle.line1.as_bytes(),
            tle.line2.as_bytes(),
        )?;

        let constants = sgp4::Constants::from_elements(&elements)?;

        Ok(Self {
            _satellite_name: tle.satellite_name.clone(),
            norad_id: tle.norad_id,
            elements,
            constants,
        })
    }

    pub fn eci_at(&self, datetime: DateTime<Utc>) -> Result<Eci, PropagationError> {
        let minutes_since_epoch = self
            .elements
            .datetime_to_minutes_since_epoch(&datetime.naive_utc())?;

        let prediction = self.constants.propagate(minutes_since_epoch)?;

        Ok(Eci::from(prediction.position))
    }
}
