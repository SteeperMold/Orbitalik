use chrono::{DateTime, Utc};
use uom::si::angle::{degree, radian};
use uom::si::f64::{Angle, Length};
use uom::si::length::{kilometer, meter, mile};

use crate::domain::errors::UnitConversionError;
use crate::transport::grpc::server::trajectory_grpc;
use crate::transport::grpc::server::trajectory_grpc::unit_settings::{AngleUnit, DistanceUnit};

pub struct TrajectoryComputationMetadata {
    pub propagation_model: String,

    pub computation_time: DateTime<Utc>,

    pub norad_id: u32,
    pub satellite_name: String,

    pub tle_epoch: DateTime<Utc>,
}

pub struct PassesComputationMetadata {
    pub propagation_model: String,

    pub computation_time: DateTime<Utc>,

    pub norad_ids: Vec<u32>,
    pub satellite_names: Vec<String>,

    pub tle_epoch: DateTime<Utc>,

    pub satellites_evaluated: u32,
    pub passes_found: u32,

    pub computation_ms: u32,
}

pub struct UnitContext {
    pub settings: trajectory_grpc::UnitSettings,
    pub distance: DistanceUnit,
    pub angle: AngleUnit,
}

impl UnitContext {
    pub fn angle(&self, value: Angle) -> Result<f64, UnitConversionError> {
        match self.angle {
            AngleUnit::Degrees => Ok(value.get::<degree>()),
            AngleUnit::Radians => Ok(value.get::<radian>()),
            AngleUnit::Unspecified => Err(UnitConversionError::UnitsUnspecified),
        }
    }

    pub fn distance(&self, value: Length) -> Result<f64, UnitConversionError> {
        match self.distance {
            DistanceUnit::Meters => Ok(value.get::<meter>()),
            DistanceUnit::Kilometers => Ok(value.get::<kilometer>()),
            DistanceUnit::Miles => Ok(value.get::<mile>()),
            DistanceUnit::Unspecified => Err(UnitConversionError::UnitsUnspecified),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::domain::models::UnitContext;
    use crate::transport::grpc::server::trajectory_grpc::UnitSettings;

    fn valid_units() -> UnitSettings {
        UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Degrees as i32,
        }
    }

    #[test]
    fn unit_context_accepts_valid_units() {
        let context = UnitContext::try_from(Some(valid_units())).unwrap();

        assert_eq!(context.distance, DistanceUnit::Meters);
        assert_eq!(context.angle, AngleUnit::Degrees);
    }

    #[test]
    fn unit_context_rejects_missing_units() {
        let result = UnitContext::try_from(None);

        assert!(result.is_err());
    }

    #[test]
    fn unit_context_rejects_unspecified_distance() {
        let units = UnitSettings {
            distance_unit: DistanceUnit::Unspecified as i32,
            angle_unit: AngleUnit::Degrees as i32,
        };

        assert!(UnitContext::try_from(Some(units)).is_err());
    }

    #[test]
    fn unit_context_rejects_unspecified_angle() {
        let units = UnitSettings {
            distance_unit: DistanceUnit::Meters as i32,
            angle_unit: AngleUnit::Unspecified as i32,
        };

        assert!(UnitContext::try_from(Some(units)).is_err());
    }
}
