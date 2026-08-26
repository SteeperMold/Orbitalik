use chrono::{DateTime, Utc};

use crate::astro::coords::teme::Teme;
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

        let constants = sgp4::Constants::from_elements_afspc_compatibility_mode(&elements)?;

        Ok(Self {
            _satellite_name: tle.satellite_name.clone(),
            norad_id: tle.norad_id,
            elements,
            constants,
        })
    }

    pub fn teme_at(&self, datetime: DateTime<Utc>) -> Result<Teme, PropagationError> {
        let minutes_since_epoch = self
            .elements
            .datetime_to_minutes_since_epoch(&datetime.naive_utc())?;

        let prediction = self.constants.propagate(minutes_since_epoch)?;

        Ok(Teme::from(prediction.position))
    }
}

#[allow(clippy::unwrap_used)]
#[cfg(test)]
mod tests {
    use crate::astro::test_utils::{assert_close, test_datetime, test_propagator};

    use uom::si::length::kilometer;

    #[test]
    fn from_tle_creates_propagator() {
        let propagator = test_propagator();

        assert!(propagator.norad_id > 0);
    }

    #[test]
    fn teme_at_returns_position() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let result = propagator.teme_at(datetime);

        assert!(result.is_ok());

        let teme = result.unwrap();

        let x = teme.x.get::<kilometer>();
        let y = teme.y.get::<kilometer>();
        let z = teme.z.get::<kilometer>();

        assert!(x != 0.0 || y != 0.0 || z != 0.0);
    }

    #[test]
    fn teme_at_matches_reference() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let teme = propagator.teme_at(datetime).unwrap();

        assert_close(teme.x.get::<kilometer>(), -672.405_384_347_387, 1e-6);
        assert_close(teme.y.get::<kilometer>(), -5_319.625_479_765_174, 1e-6);
        assert_close(teme.z.get::<kilometer>(), -4_199.191_899_815_303, 1e-6);
    }

    #[test]
    fn teme_at_changes_with_time() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let first = propagator.teme_at(datetime).unwrap();

        let later = datetime + chrono::Duration::minutes(10);
        let second = propagator.teme_at(later).unwrap();

        assert!(
            (first.x.get::<kilometer>() - second.x.get::<kilometer>()).abs() > 1e-6
                || (first.y.get::<kilometer>() - second.y.get::<kilometer>()).abs() > 1e-6
                || (first.z.get::<kilometer>() - second.z.get::<kilometer>()).abs() > 1e-6
        );
    }

    #[test]
    fn teme_at_is_deterministic() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let first = propagator.teme_at(datetime).unwrap();
        let second = propagator.teme_at(datetime).unwrap();

        assert_close(
            first.x.get::<kilometer>(),
            second.x.get::<kilometer>(),
            1e-12,
        );
        assert_close(
            first.y.get::<kilometer>(),
            second.y.get::<kilometer>(),
            1e-12,
        );
        assert_close(
            first.z.get::<kilometer>(),
            second.z.get::<kilometer>(),
            1e-12,
        );
    }

    #[test]
    fn teme_at_position_is_reasonable() {
        let propagator = test_propagator();
        let datetime = test_datetime();

        let teme = propagator.teme_at(datetime).unwrap();

        let x = teme.x.get::<kilometer>();
        let y = teme.y.get::<kilometer>();
        let z = teme.z.get::<kilometer>();

        let radius = (x * x + y * y + z * z).sqrt();

        assert!(radius > 6_000.0);
        assert!(radius < 50_000.0);
    }
}
