use tonic::Status;
use uom::si::angle::{degree, radian};
use uom::si::f64::{Angle, Length};
use uom::si::length::meter;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{Sampling, SatelliteIdentifier, TimeRange};
use crate::domain::models::UnitContext;
use crate::transport::grpc::conversions::timestamp::ToChrono;
use crate::transport::grpc::server::trajectory_grpc;
use crate::transport::grpc::server::trajectory_grpc::geodetic_input;
use crate::transport::grpc::server::trajectory_grpc::unit_settings::{AngleUnit, DistanceUnit};
use trajectory_grpc::satellite_identifier;

impl TryFrom<trajectory_grpc::SatelliteIdentifier> for SatelliteIdentifier {
    type Error = Status;

    fn try_from(value: trajectory_grpc::SatelliteIdentifier) -> Result<Self, Self::Error> {
        match value.kind {
            Some(satellite_identifier::Kind::NoradId(id)) => Ok(Self::NoradId(id)),
            Some(satellite_identifier::Kind::SatelliteName(name)) => Ok(Self::Name(name)),
            None => Err(Status::invalid_argument("Missing satellite identifier")),
        }
    }
}

impl TryFrom<trajectory_grpc::TimeRange> for TimeRange {
    type Error = Status;

    fn try_from(value: trajectory_grpc::TimeRange) -> Result<Self, Self::Error> {
        let start = value
            .start
            .ok_or_else(|| Status::invalid_argument("Missing start timestamp"))?
            .to_chrono()
            .map_err(|e| Status::invalid_argument(format!("Invalid start timestamp: {e}")))?;

        let end = value
            .end
            .ok_or_else(|| Status::invalid_argument("Missing end timestamp"))?
            .to_chrono()
            .map_err(|e| Status::invalid_argument(format!("Invalid end timestamp: {e}")))?;

        Self::new(start, end).map_err(|_| Status::invalid_argument("Start must be <= end"))
    }
}

impl TryFrom<trajectory_grpc::GeodeticInput> for Geodetic {
    type Error = Status;

    fn try_from(value: trajectory_grpc::GeodeticInput) -> Result<Self, Self::Error> {
        let lat = value
            .lat
            .ok_or_else(|| Status::invalid_argument("Missing latitude"))?;

        let lon = value
            .lon
            .ok_or_else(|| Status::invalid_argument("Missing longitude"))?;

        let alt = value
            .alt
            .ok_or_else(|| Status::invalid_argument("Missing alt"))?;

        let lat_rad = match lat {
            geodetic_input::Lat::LatRad(value) => value,
            geodetic_input::Lat::LatDeg(value) => value.to_radians(),
        };

        let lon_rad = match lon {
            geodetic_input::Lon::LonRad(value) => value,
            geodetic_input::Lon::LonDeg(value) => value.to_radians(),
        };

        let alt_m = match alt {
            geodetic_input::Alt::AltM(value) => value,
            geodetic_input::Alt::AltKm(value) => value * 1000.0,
        };

        Ok(Self {
            lat: Angle::new::<radian>(lat_rad),
            lon: Angle::new::<radian>(lon_rad),
            alt: Length::new::<meter>(alt_m),
        })
    }
}

impl TryFrom<Option<trajectory_grpc::UnitSettings>> for UnitContext {
    type Error = Status;

    fn try_from(value: Option<trajectory_grpc::UnitSettings>) -> Result<Self, Status> {
        let settings =
            value.ok_or_else(|| Status::invalid_argument("Measurement units must be specified"))?;

        let distance = DistanceUnit::try_from(settings.distance_unit)
            .map_err(|_| Status::invalid_argument("Invalid distance unit"))?;

        let angle = AngleUnit::try_from(settings.angle_unit)
            .map_err(|_| Status::invalid_argument("Invalid angle unit"))?;

        if distance == DistanceUnit::Unspecified || angle == AngleUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Measurement units must be specified",
            ));
        }

        Ok(Self {
            settings,
            distance,
            angle,
        })
    }
}

impl TryFrom<trajectory_grpc::SamplingOptions> for Sampling {
    type Error = Status;

    fn try_from(value: trajectory_grpc::SamplingOptions) -> Result<Self, Status> {
        Self::new(value.step_seconds).map_err(|_| Status::invalid_argument("Step must be > 0"))
    }
}

pub trait IntoAngle {
    fn into_angle(self) -> Angle;
}

impl IntoAngle for trajectory_grpc::pass_prediction_request::MinElevation {
    fn into_angle(self) -> Angle {
        match self {
            Self::MinElevationDeg(value) => Angle::new::<degree>(value),
            Self::MinElevationRad(value) => Angle::new::<radian>(value),
        }
    }
}

impl IntoAngle for trajectory_grpc::pass_prediction_request::MinPeakElevation {
    fn into_angle(self) -> Angle {
        match self {
            Self::MinPeakElevationDeg(value) => Angle::new::<degree>(value),
            Self::MinPeakElevationRad(value) => Angle::new::<radian>(value),
        }
    }
}

impl IntoAngle for trajectory_grpc::next_passes_request::MinElevation {
    fn into_angle(self) -> Angle {
        match self {
            Self::MinElevationDeg(value) => Angle::new::<degree>(value),
            Self::MinElevationRad(value) => Angle::new::<radian>(value),
        }
    }
}

impl IntoAngle for trajectory_grpc::next_passes_request::MinPeakElevation {
    fn into_angle(self) -> Angle {
        match self {
            Self::MinPeakElevationDeg(value) => Angle::new::<degree>(value),
            Self::MinPeakElevationRad(value) => Angle::new::<radian>(value),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    use uom::si::angle::{degree, radian};

    use crate::transport::grpc::server::trajectory_grpc::{
        next_passes_request, pass_prediction_request,
    };

    #[test]
    fn into_angle_converts_degrees_and_radians() {
        let degrees = pass_prediction_request::MinElevation::MinElevationDeg(45.0).into_angle();

        let radians = pass_prediction_request::MinElevation::MinElevationRad(std::f64::consts::PI)
            .into_angle();

        assert_eq!(degrees.get::<degree>(), 45.0);
        assert_eq!(radians.get::<radian>(), std::f64::consts::PI);
    }

    #[test]
    fn all_angle_oneofs_convert() {
        let values = [
            pass_prediction_request::MinElevation::MinElevationDeg(10.0).into_angle(),
            pass_prediction_request::MinPeakElevation::MinPeakElevationRad(0.5).into_angle(),
            next_passes_request::MinElevation::MinElevationDeg(20.0).into_angle(),
            next_passes_request::MinPeakElevation::MinPeakElevationRad(1.0).into_angle(),
        ];

        assert_eq!(values[0].get::<degree>(), 10.0);
        assert_eq!(values[1].get::<radian>(), 0.5);
        assert_eq!(values[2].get::<degree>(), 20.0);
        assert_eq!(values[3].get::<radian>(), 1.0);
    }
}
