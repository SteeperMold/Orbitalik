use tonic::Status;
use uom::si::f64::Length;

use crate::astro::coords::ecef::Ecef;
use crate::astro::coords::geodetic::Geodetic;
use crate::astro::coords::teme::Teme;
use crate::astro::models::{LookAngles, SatellitePosition};
use crate::domain::models::UnitContext;
use crate::transport::grpc::conversions::timestamp::ToProtoTimestamp;
use crate::transport::grpc::server::trajectory_grpc;

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

impl trajectory_grpc::Vector3 {
    pub fn from_xyz<T: HasXYZ>(coords: &T, unit_context: &UnitContext) -> Result<Self, Status> {
        Ok(Self {
            x: unit_context.distance(coords.x())?,
            y: unit_context.distance(coords.y())?,
            z: unit_context.distance(coords.z())?,
        })
    }
}

impl trajectory_grpc::GeodeticOutput {
    pub fn from_domain(geodetic: &Geodetic, unit_context: &UnitContext) -> Result<Self, Status> {
        Ok(Self {
            lat: unit_context.angle(geodetic.lat)?,
            lon: unit_context.angle(geodetic.lon)?,
            alt: unit_context.distance(geodetic.alt)?,
        })
    }
}

impl trajectory_grpc::StateVector {
    pub fn from_domain(
        position: &SatellitePosition,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            datetime: Some(position.time.to_proto_timestamp()?),

            position_teme: position
                .teme
                .as_ref()
                .map(|teme| trajectory_grpc::Vector3::from_xyz(teme, unit_context))
                .transpose()?,

            velocity_teme: None,

            position_ecef: position
                .ecef
                .as_ref()
                .map(|ecef| trajectory_grpc::Vector3::from_xyz(ecef, unit_context))
                .transpose()?,

            velocity_ecef: None,

            geodetic: position
                .geodetic
                .as_ref()
                .map(|g| trajectory_grpc::GeodeticOutput::from_domain(g, unit_context))
                .transpose()?,
        })
    }
}

impl trajectory_grpc::ObserverTrajectoryPoint {
    pub fn from_look_angles(
        look_angles: &LookAngles,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            datetime: Some(look_angles.time.to_proto_timestamp()?),

            azimuth: look_angles
                .azimuth
                .map(|az| unit_context.angle(az))
                .transpose()?,

            elevation: look_angles
                .elevation
                .map(|el| unit_context.angle(el))
                .transpose()?,

            range: look_angles
                .range
                .map(|r| unit_context.distance(r))
                .transpose()?,
        })
    }
}
