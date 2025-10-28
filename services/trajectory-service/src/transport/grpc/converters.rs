use chrono::{DateTime, TimeZone, Utc};
use tonic::Status;
use uom::si::angle::{degree, radian};
use uom::si::f64::{Angle, Length};
use uom::si::length::{kilometer, meter};

use crate::astro::coords::{ecef::Ecef, eci::Eci, geodetic::Geodetic};
use crate::astro::models::{LookAngles, SatellitePosition};
use crate::domain::errors::TimestampConversionError;
use crate::domain::models::SatelliteIdentifier;
use crate::transport::adapter::tle_client::tle_grpc;
use crate::transport::grpc::trajectory::trajectory_grpc;
use crate::transport::grpc::trajectory::trajectory_grpc::{
    GeodeticInput, GeodeticOutput, LookAnglesResponse, PositionResponse, Vector3, geodetic_input,
};

pub trait ToChrono {
    fn to_chrono(&self) -> Result<DateTime<Utc>, TimestampConversionError>;
}

impl ToChrono for prost_types::Timestamp {
    fn to_chrono(&self) -> Result<DateTime<Utc>, TimestampConversionError> {
        let nanos = u32::try_from(self.nanos)?;

        Utc.timestamp_opt(self.seconds, nanos)
            .single()
            .ok_or(TimestampConversionError::InvalidTimestamp)
    }
}

impl From<Eci> for Vector3 {
    fn from(eci: Eci) -> Self {
        Self {
            x_km: eci.x.get::<kilometer>(),
            y_km: eci.y.get::<kilometer>(),
            z_km: eci.z.get::<kilometer>(),
            x_m: eci.x.get::<meter>(),
            y_m: eci.y.get::<meter>(),
            z_m: eci.z.get::<meter>(),
        }
    }
}

impl From<Ecef> for Vector3 {
    fn from(ecef: Ecef) -> Self {
        Self {
            x_km: ecef.x.get::<kilometer>(),
            y_km: ecef.y.get::<kilometer>(),
            z_km: ecef.z.get::<kilometer>(),
            x_m: ecef.x.get::<meter>(),
            y_m: ecef.y.get::<meter>(),
            z_m: ecef.z.get::<meter>(),
        }
    }
}

impl From<Geodetic> for GeodeticOutput {
    fn from(g: Geodetic) -> Self {
        Self {
            lat_deg: g.lat.get::<degree>(),
            lon_deg: g.lon.get::<degree>(),
            lat_rad: g.lat.get::<radian>(),
            lon_rad: g.lon.get::<radian>(),
            alt_km: g.alt.get::<kilometer>(),
            alt_m: g.alt.get::<meter>(),
        }
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

impl From<SatellitePosition> for PositionResponse {
    fn from(value: SatellitePosition) -> Self {
        Self {
            eci: Some(value.eci.into()),
            ecef: Some(value.ecef.into()),
            geodetic: Some(value.geodetic.into()),
        }
    }
}

impl From<LookAngles> for LookAnglesResponse {
    fn from(value: LookAngles) -> Self {
        Self {
            azimuth_rad: value.azimuth.get::<radian>(),
            azimuth_deg: value.azimuth.get::<degree>(),
            elevation_rad: value.elevation.get::<radian>(),
            elevation_deg: value.elevation.get::<degree>(),
            range_km: value.range.get::<kilometer>(),
            range_m: value.range.get::<meter>(),
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
