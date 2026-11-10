use async_trait::async_trait;
use chrono::{DateTime, Utc};

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{LookAngles, SatelliteIdentifier, SatellitePosition};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::errors::ServiceError;
use crate::domain::models::TrajectoryComputationMetadata;

#[cfg_attr(test, mockall::automock)]
#[async_trait]
pub trait PositionServiceApi: Send + Sync {
    async fn get_position(
        &self,
        identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        compute: &PositionComputation,
    ) -> Result<(SatellitePosition, TrajectoryComputationMetadata), ServiceError>;

    async fn get_look_angles(
        &self,
        identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(LookAngles, TrajectoryComputationMetadata), ServiceError>;
}
