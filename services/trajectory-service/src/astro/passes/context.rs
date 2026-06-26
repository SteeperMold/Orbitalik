use crate::astro::errors::PropagationError;
use crate::astro::models::{SatelliteIdentifier, Tle};
use crate::astro::propagation::propagator::Propagator;

pub struct SatelliteContext {
    pub identifier: SatelliteIdentifier,
    pub propagator: Propagator,
}

impl SatelliteContext {
    pub fn from_tle(tle: &Tle) -> Result<Self, PropagationError> {
        Ok(Self {
            identifier: SatelliteIdentifier::NoradId(tle.norad_id),
            propagator: Propagator::from_tle(tle)?,
        })
    }
}
