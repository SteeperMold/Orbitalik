use chrono::Utc;
use std::sync::Arc;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{ObserverTrajectory, Trajectory};
use crate::astro::models::{Sampling, SatelliteIdentifier, TimeRange};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::astro::propagation::propagator::Propagator;
use crate::domain::errors::ServiceError;
use crate::domain::models::TrajectoryComputationMetadata;
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct TrajectoryService {
    tle_grpc_client: Arc<TleGrpcClient>,
}

impl TrajectoryService {
    pub const fn new(tle_grpc_client: Arc<TleGrpcClient>) -> Self {
        Self { tle_grpc_client }
    }

    pub async fn get_trajectory(
        &self,
        satellite_identifier: SatelliteIdentifier,
        range: TimeRange,
        sampling: Sampling,
        compute: &PositionComputation,
    ) -> Result<(Trajectory, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let trajectory = Propagator::from_tle(&tle)?.trajectory_at(range, sampling, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((trajectory, metadata))
    }

    pub async fn get_observer_trajectory(
        &self,
        satellite_identifier: SatelliteIdentifier,
        range: TimeRange,
        sampling: Sampling,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(ObserverTrajectory, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let observer_trajectory = Propagator::from_tle(&tle)?
            .observer_trajectory_at(range, sampling, observer, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((observer_trajectory, metadata))
    }
}
