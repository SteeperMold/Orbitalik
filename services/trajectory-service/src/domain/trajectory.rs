use async_trait::async_trait;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{
    ObserverTrajectory, Sampling, SatelliteIdentifier, TimeRange, Trajectory,
};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::errors::ServiceError;
use crate::domain::models::TrajectoryComputationMetadata;

#[cfg_attr(test, mockall::automock)]
#[async_trait]
pub trait TrajectoryServiceApi: Send + Sync {
    async fn get_trajectory(
        &self,
        satellite_identifier: SatelliteIdentifier,
        range: TimeRange,
        sampling: Sampling,
        compute: &PositionComputation,
    ) -> Result<(Trajectory, TrajectoryComputationMetadata), ServiceError>;

    async fn get_observer_trajectory(
        &self,
        satellite_identifier: SatelliteIdentifier,
        range: TimeRange,
        sampling: Sampling,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(ObserverTrajectory, TrajectoryComputationMetadata), ServiceError>;
}
