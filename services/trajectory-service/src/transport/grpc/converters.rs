use chrono::{DateTime, TimeZone, Utc};
use prost_types::{FieldMask, Timestamp};
use tonic::Status;
use uom::si::angle::{degree, radian};
use uom::si::f64::{Angle, Length};
use uom::si::length::{kilometer, meter, mile};

use crate::astro::coords::ecef::Ecef;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::coords::teme::Teme;
use crate::astro::models::{LookAngles, ObserverTrajectory, Pass, SatellitePosition, Trajectory};
use crate::astro::models::{Sampling, SatelliteIdentifier, TimeRange};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::errors::TimestampConversionError;
use crate::domain::models::{PassesComputationMetadata, TrajectoryComputationMetadata};
use crate::transport::adapter::tle_client::tle_grpc;
use crate::transport::grpc::server::trajectory_grpc;
use crate::transport::grpc::server::trajectory_grpc::unit_settings::{AngleUnit, DistanceUnit};
use crate::transport::grpc::server::trajectory_grpc::{
    GeodeticInput, GeodeticOutput, UnitSettings, Vector3, geodetic_input,
};

pub trait ToChrono {
    fn to_chrono(&self) -> Result<DateTime<Utc>, TimestampConversionError>;
}

impl ToChrono for Timestamp {
    fn to_chrono(&self) -> Result<DateTime<Utc>, TimestampConversionError> {
        let nanos = u32::try_from(self.nanos)?;

        Utc.timestamp_opt(self.seconds, nanos)
            .single()
            .ok_or(TimestampConversionError::InvalidTimestamp)
    }
}

pub trait ToProtoTimestamp {
    fn to_proto_timestamp(&self) -> Result<Timestamp, TimestampConversionError>;
}

impl ToProtoTimestamp for DateTime<Utc> {
    fn to_proto_timestamp(&self) -> Result<Timestamp, TimestampConversionError> {
        let nanos = i32::try_from(self.timestamp_subsec_nanos())?;

        Ok(Timestamp {
            seconds: self.timestamp(),
            nanos,
        })
    }
}

impl TryFrom<GeodeticInput> for Geodetic {
    type Error = Status;

    fn try_from(value: GeodeticInput) -> Result<Self, Self::Error> {
        let lat_rad = match value
            .lat
            .ok_or_else(|| Status::invalid_argument("Missing latitude"))?
        {
            geodetic_input::Lat::LatRad(r) => r,
            geodetic_input::Lat::LatDeg(d) => d.to_radians(),
        };

        let lon_rad = match value
            .lon
            .ok_or_else(|| Status::invalid_argument("Missing longitude"))?
        {
            geodetic_input::Lon::LonRad(r) => r,
            geodetic_input::Lon::LonDeg(d) => d.to_radians(),
        };

        let alt_m = match value
            .alt
            .ok_or_else(|| Status::invalid_argument("Missing alt"))?
        {
            geodetic_input::Alt::AltM(m) => m,
            geodetic_input::Alt::AltKm(km) => km * 1000.0,
        };

        Ok(Self {
            lat: Angle::new::<radian>(lat_rad),
            lon: Angle::new::<radian>(lon_rad),
            alt: Length::new::<meter>(alt_m),
        })
    }
}

impl From<&FieldMask> for PositionComputation {
    fn from(mask: &FieldMask) -> Self {
        let has = |prefix: &str| mask.paths.iter().any(|p| p.starts_with(prefix));
        Self {
            teme: has("eci"),
            ecef: has("ecef"),
            geodetic: has("geodetic"),
        }
    }
}

impl From<&FieldMask> for LookAnglesComputation {
    fn from(mask: &FieldMask) -> Self {
        let has = |prefix: &str| mask.paths.iter().any(|p| p.starts_with(prefix));
        Self {
            azimuth: has("azimuth"),
            elevation: has("elevation"),
            range: has("range"),
        }
    }
}

impl From<SatelliteIdentifier> for tle_grpc::SatelliteIdentifier {
    fn from(identifier: SatelliteIdentifier) -> Self {
        match identifier {
            SatelliteIdentifier::NoradId(id) => Self {
                kind: Some(tle_grpc::satellite_identifier::Kind::NoradId(id)),
            },
            SatelliteIdentifier::Name(name) => Self {
                kind: Some(tle_grpc::satellite_identifier::Kind::SatelliteName(name)),
            },
        }
    }
}

impl TryFrom<trajectory_grpc::SatelliteIdentifier> for SatelliteIdentifier {
    type Error = Status;

    fn try_from(value: trajectory_grpc::SatelliteIdentifier) -> Result<Self, Self::Error> {
        match value.kind {
            Some(trajectory_grpc::satellite_identifier::Kind::NoradId(id)) => Ok(Self::NoradId(id)),
            Some(trajectory_grpc::satellite_identifier::Kind::SatelliteName(name)) => {
                Ok(Self::Name(name))
            }
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

impl TryFrom<trajectory_grpc::SamplingOptions> for Sampling {
    type Error = Status;

    fn try_from(value: trajectory_grpc::SamplingOptions) -> Result<Self, Self::Error> {
        Self::new(value.step_seconds).map_err(|_| Status::invalid_argument("Step must be > 0"))
    }
}

impl From<SatelliteIdentifier> for trajectory_grpc::SatelliteIdentifier {
    fn from(value: SatelliteIdentifier) -> Self {
        match value {
            SatelliteIdentifier::NoradId(id) => Self {
                kind: Some(trajectory_grpc::satellite_identifier::Kind::NoradId(id)),
            },
            SatelliteIdentifier::Name(name) => Self {
                kind: Some(trajectory_grpc::satellite_identifier::Kind::SatelliteName(
                    name,
                )),
            },
        }
    }
}

impl trajectory_grpc::TrajectoryComputationMetadata {
    pub fn with_units(
        metadata: TrajectoryComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Option<Self>, Status> {
        Ok(Some(Self {
            propagation_model: metadata.propagation_model,
            computation_time: Some(metadata.computation_time.to_proto_timestamp()?),
            norad_id: metadata.norad_id,
            satellite_name: metadata.satellite_name,
            tle_epoch: Some(metadata.tle_epoch.to_proto_timestamp()?),
            units,
        }))
    }
}

impl trajectory_grpc::PassesComputationMetadata {
    pub fn with_units(
        metadata: PassesComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Option<Self>, Status> {
        Ok(Some(Self {
            propagation_model: metadata.propagation_model,
            computation_time: Some(metadata.computation_time.to_proto_timestamp()?),
            norad_ids: metadata.norad_ids,
            satellite_names: metadata.satellite_names,
            tle_epoch: Some(metadata.tle_epoch.to_proto_timestamp()?),
            satellites_evaluated: metadata.satellites_evaluated,
            passes_found: metadata.passes_found,
            computation_ms: metadata.computation_ms,
            units,
        }))
    }
}

pub trait HasXYZ {
    fn x(&self) -> Length;
    fn y(&self) -> Length;
    fn z(&self) -> Length;
}

impl HasXYZ for Teme {
    fn x(&self) -> Length {
        self.x
    }
    fn y(&self) -> Length {
        self.y
    }
    fn z(&self) -> Length {
        self.z
    }
}

impl HasXYZ for Ecef {
    fn x(&self) -> Length {
        self.x
    }
    fn y(&self) -> Length {
        self.y
    }
    fn z(&self) -> Length {
        self.z
    }
}

impl Vector3 {
    pub fn from_xyz<T: HasXYZ>(
        coords: Option<&T>,
        units: Option<UnitSettings>,
    ) -> Result<Option<Self>, Status> {
        let Some(coords) = coords else {
            return Ok(None);
        };

        let distance_unit = units
            .as_ref()
            .and_then(|u| DistanceUnit::try_from(u.distance_unit).ok())
            .unwrap_or(DistanceUnit::Unspecified);

        if distance_unit == DistanceUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Distance unit is unspecified in UnitSettings",
            ));
        }

        let (x, y, z) = match distance_unit {
            DistanceUnit::Meters => (
                coords.x().get::<meter>(),
                coords.y().get::<meter>(),
                coords.z().get::<meter>(),
            ),
            DistanceUnit::Kilometers => (
                coords.x().get::<kilometer>(),
                coords.y().get::<kilometer>(),
                coords.z().get::<kilometer>(),
            ),
            DistanceUnit::Miles => (
                coords.x().get::<mile>(),
                coords.y().get::<mile>(),
                coords.z().get::<mile>(),
            ),
            DistanceUnit::Unspecified => unreachable!(),
        };

        Ok(Some(Self { x, y, z }))
    }
}

type LookAnglesValues = (Option<f64>, Option<f64>, Option<f64>);

fn look_angles_values(
    look_angles: &LookAngles,
    units: Option<UnitSettings>,
) -> Result<LookAnglesValues, Status> {
    let distance_unit = units
        .as_ref()
        .and_then(|u| DistanceUnit::try_from(u.distance_unit).ok())
        .unwrap_or(DistanceUnit::Unspecified);

    let angle_unit = units
        .as_ref()
        .and_then(|u| AngleUnit::try_from(u.angle_unit).ok())
        .unwrap_or(AngleUnit::Unspecified);

    if distance_unit == DistanceUnit::Unspecified {
        return Err(Status::invalid_argument(
            "Distance unit is unspecified in UnitSettings",
        ));
    }

    if angle_unit == AngleUnit::Unspecified {
        return Err(Status::invalid_argument(
            "Angle unit is unspecified in UnitSettings",
        ));
    }

    let azimuth = match (&look_angles.azimuth, angle_unit) {
        (Some(a), AngleUnit::Degrees) => Some(a.get::<degree>()),
        (Some(a), AngleUnit::Radians) => Some(a.get::<radian>()),
        (None, _) => None,
        (_, AngleUnit::Unspecified) => unreachable!(),
    };

    let elevation = match (&look_angles.elevation, angle_unit) {
        (Some(a), AngleUnit::Degrees) => Some(a.get::<degree>()),
        (Some(a), AngleUnit::Radians) => Some(a.get::<radian>()),
        (None, _) => None,
        (_, AngleUnit::Unspecified) => unreachable!(),
    };

    let range = match (&look_angles.range, distance_unit) {
        (Some(r), DistanceUnit::Meters) => Some(r.get::<meter>()),
        (Some(r), DistanceUnit::Kilometers) => Some(r.get::<kilometer>()),
        (Some(r), DistanceUnit::Miles) => Some(r.get::<mile>()),
        (None, _) => None,
        (_, DistanceUnit::Unspecified) => unreachable!(),
    };

    Ok((azimuth, elevation, range))
}

impl GeodeticOutput {
    pub fn from_geodetic(
        geodetic: Option<&Geodetic>,
        units: Option<UnitSettings>,
    ) -> Result<Option<Self>, Status> {
        let Some(g) = geodetic else { return Ok(None) };

        let distance_unit = units
            .as_ref()
            .and_then(|u| DistanceUnit::try_from(u.distance_unit).ok())
            .unwrap_or(DistanceUnit::Unspecified);

        let angle_unit = units
            .as_ref()
            .and_then(|u| AngleUnit::try_from(u.angle_unit).ok())
            .unwrap_or(AngleUnit::Unspecified);

        if distance_unit == DistanceUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Distance unit is unspecified in UnitSettings",
            ));
        }

        if angle_unit == AngleUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Angle unit is unspecified in UnitSettings",
            ));
        }

        let (lat, lon) = match angle_unit {
            AngleUnit::Degrees => (g.lat.get::<degree>(), g.lon.get::<degree>()),
            AngleUnit::Radians => (g.lat.get::<radian>(), g.lon.get::<radian>()),
            AngleUnit::Unspecified => unreachable!(),
        };

        let alt = match distance_unit {
            DistanceUnit::Meters => g.alt.get::<meter>(),
            DistanceUnit::Kilometers => g.alt.get::<kilometer>(),
            DistanceUnit::Miles => g.alt.get::<mile>(),
            DistanceUnit::Unspecified => unreachable!(),
        };

        Ok(Some(Self { lat, lon, alt }))
    }
}

impl trajectory_grpc::StateVector {
    pub fn from_position(
        position: &SatellitePosition,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        Ok(Self {
            datetime: Some(position.time.to_proto_timestamp()?),
            position_teme: Vector3::from_xyz(position.teme.as_ref(), units)?,
            velocity_teme: None,
            position_ecef: Vector3::from_xyz(position.ecef.as_ref(), units)?,
            velocity_ecef: None,
            geodetic: GeodeticOutput::from_geodetic(position.geodetic.as_ref(), units)?,
        })
    }
}

impl trajectory_grpc::ObserverTrajectoryPoint {
    pub fn from_look_angles(
        look_angles: &LookAngles,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        let (azimuth, elevation, range) = look_angles_values(look_angles, units)?;

        Ok(Self {
            datetime: Some(look_angles.time.to_proto_timestamp()?),
            azimuth,
            elevation,
            range,
        })
    }
}

impl trajectory_grpc::PositionResponse {
    pub fn from_position(
        position: &SatellitePosition,
        metadata: TrajectoryComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        Ok(Self {
            time: Some(position.time.to_proto_timestamp()?),
            metadata: trajectory_grpc::TrajectoryComputationMetadata::with_units(metadata, units)?,
            teme: Vector3::from_xyz(position.teme.as_ref(), units)?,
            ecef: Vector3::from_xyz(position.ecef.as_ref(), units)?,
            geodetic: GeodeticOutput::from_geodetic(position.geodetic.as_ref(), units)?,
        })
    }
}

impl trajectory_grpc::LookAnglesResponse {
    pub fn from_look_angles(
        look_angles: &LookAngles,
        metadata: TrajectoryComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        let distance_unit = units
            .as_ref()
            .and_then(|u| DistanceUnit::try_from(u.distance_unit).ok())
            .unwrap_or(DistanceUnit::Unspecified);

        let angle_unit = units
            .as_ref()
            .and_then(|u| AngleUnit::try_from(u.angle_unit).ok())
            .unwrap_or(AngleUnit::Unspecified);

        if distance_unit == DistanceUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Distance unit is unspecified in UnitSettings",
            ));
        }

        if angle_unit == AngleUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Angle unit is unspecified in UnitSettings",
            ));
        }

        let azimuth = match (&look_angles.azimuth, angle_unit) {
            (Some(a), AngleUnit::Degrees) => Some(a.get::<degree>()),
            (Some(a), AngleUnit::Radians) => Some(a.get::<radian>()),
            (None, _) => None,
            (_, AngleUnit::Unspecified) => unreachable!(),
        };

        let elevation = match (&look_angles.elevation, angle_unit) {
            (Some(a), AngleUnit::Degrees) => Some(a.get::<degree>()),
            (Some(a), AngleUnit::Radians) => Some(a.get::<radian>()),
            (None, _) => None,
            (_, AngleUnit::Unspecified) => unreachable!(),
        };

        let range = match (&look_angles.range, distance_unit) {
            (Some(r), DistanceUnit::Meters) => Some(r.get::<meter>()),
            (Some(r), DistanceUnit::Kilometers) => Some(r.get::<kilometer>()),
            (Some(r), DistanceUnit::Miles) => Some(r.get::<mile>()),
            (None, _) => None,
            (_, DistanceUnit::Unspecified) => unreachable!(),
        };

        Ok(Self {
            time: Some(look_angles.time.to_proto_timestamp()?),
            metadata: trajectory_grpc::TrajectoryComputationMetadata::with_units(metadata, units)?,
            azimuth,
            elevation,
            range,
        })
    }
}

impl trajectory_grpc::TrajectoryResponse {
    pub fn from_trajectory(
        trajectory: &Trajectory,
        metadata: TrajectoryComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: trajectory_grpc::TrajectoryComputationMetadata::with_units(metadata, units)?,
            states: trajectory
                .samples
                .iter()
                .map(|p| trajectory_grpc::StateVector::from_position(p, units))
                .collect::<Result<Vec<_>, Status>>()?,
            range: Some(trajectory_grpc::TimeRange {
                start: Some(trajectory.start.to_proto_timestamp()?),
                end: Some(trajectory.end.to_proto_timestamp()?),
            }),
            sampling: Some(trajectory_grpc::SamplingOptions {
                step_seconds: trajectory.step_seconds,
            }),
        })
    }
}

impl trajectory_grpc::ObserverTrajectoryResponse {
    pub fn from_observer_trajectory(
        observer_trajectory: &ObserverTrajectory,
        metadata: TrajectoryComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: trajectory_grpc::TrajectoryComputationMetadata::with_units(metadata, units)?,
            points: observer_trajectory
                .samples
                .iter()
                .map(|p| trajectory_grpc::ObserverTrajectoryPoint::from_look_angles(p, units))
                .collect::<Result<Vec<_>, Status>>()?,
            range: Some(trajectory_grpc::TimeRange {
                start: Some(observer_trajectory.start.to_proto_timestamp()?),
                end: Some(observer_trajectory.end.to_proto_timestamp()?),
            }),
            sampling: Some(trajectory_grpc::SamplingOptions {
                step_seconds: observer_trajectory.step_seconds,
            }),
        })
    }
}

impl trajectory_grpc::Pass {
    pub fn from_pass(pass: &Pass, units: Option<UnitSettings>) -> Result<Self, Status> {
        let angle_unit = units
            .as_ref()
            .and_then(|u| AngleUnit::try_from(u.angle_unit).ok())
            .unwrap_or(AngleUnit::Unspecified);

        if angle_unit == AngleUnit::Unspecified {
            return Err(Status::invalid_argument(
                "Angle unit is unspecified in UnitSettings",
            ));
        }

        let aos_azimuth = match (&pass.aos_azimuth, angle_unit) {
            (a, AngleUnit::Degrees) => a.get::<degree>(),
            (a, AngleUnit::Radians) => a.get::<radian>(),
            (_, AngleUnit::Unspecified) => unreachable!(),
        };

        let max_elevation = match (&pass.max_elevation, angle_unit) {
            (a, AngleUnit::Degrees) => a.get::<degree>(),
            (a, AngleUnit::Radians) => a.get::<radian>(),
            (_, AngleUnit::Unspecified) => unreachable!(),
        };

        let max_elevation_azimuth = match (&pass.max_elevation_azimuth, angle_unit) {
            (a, AngleUnit::Degrees) => a.get::<degree>(),
            (a, AngleUnit::Radians) => a.get::<radian>(),
            (_, AngleUnit::Unspecified) => unreachable!(),
        };

        let los_azimuth = match (&pass.los_azimuth, angle_unit) {
            (a, AngleUnit::Degrees) => a.get::<degree>(),
            (a, AngleUnit::Radians) => a.get::<radian>(),
            (_, AngleUnit::Unspecified) => unreachable!(),
        };

        Ok(Self {
            satellite: Some(pass.satellite.clone().into()),

            aos: Some(pass.aos.to_proto_timestamp()?),
            aos_azimuth,

            max_elevation_time: Some(pass.max_elevation_time.to_proto_timestamp()?),
            max_elevation,
            max_elevation_azimuth,

            los: Some(pass.los.to_proto_timestamp()?),
            los_azimuth,

            duration_seconds: pass.duration_seconds,
        })
    }
}

impl trajectory_grpc::GetPassesResponse {
    pub fn from_passes(
        passes: &[Pass],
        metadata: PassesComputationMetadata,
        range: TimeRange,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: trajectory_grpc::PassesComputationMetadata::with_units(metadata, units)?,
            passes: passes
                .iter()
                .map(|p| trajectory_grpc::Pass::from_pass(p, units))
                .collect::<Result<Vec<_>, Status>>()?,
            range: Some(trajectory_grpc::TimeRange {
                start: Some(range.start.to_proto_timestamp()?),
                end: Some(range.end.to_proto_timestamp()?),
            }),
        })
    }
}

impl trajectory_grpc::NextPassesResponse {
    pub fn from_passes(
        passes: &[Pass],
        metadata: PassesComputationMetadata,
        units: Option<UnitSettings>,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: trajectory_grpc::PassesComputationMetadata::with_units(metadata, units)?,
            passes: passes
                .iter()
                .map(|p| trajectory_grpc::Pass::from_pass(p, units))
                .collect::<Result<Vec<_>, Status>>()?,
        })
    }
}
