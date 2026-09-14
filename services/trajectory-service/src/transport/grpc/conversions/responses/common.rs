use crate::astro::models::SatelliteIdentifier;
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::models::UnitContext;
use crate::transport::adapter::tle_client::tle_grpc;
use crate::transport::grpc::server::trajectory_grpc;
use crate::transport::grpc::server::trajectory_grpc::satellite_identifier;
use prost_types::FieldMask;

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

impl From<SatelliteIdentifier> for trajectory_grpc::SatelliteIdentifier {
    fn from(value: SatelliteIdentifier) -> Self {
        match value {
            SatelliteIdentifier::NoradId(id) => Self {
                kind: Some(satellite_identifier::Kind::NoradId(id)),
            },
            SatelliteIdentifier::Name(name) => Self {
                kind: Some(satellite_identifier::Kind::SatelliteName(name)),
            },
        }
    }
}

impl From<&UnitContext> for trajectory_grpc::UnitSettings {
    fn from(value: &UnitContext) -> Self {
        value.settings
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
