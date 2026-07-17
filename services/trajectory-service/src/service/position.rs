use chrono::{DateTime, Utc};
use std::sync::Arc;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{LookAngles, SatelliteIdentifier, SatellitePosition};
use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::astro::propagation::propagator::Propagator;
use crate::domain::errors::ServiceError;
use crate::domain::models::TrajectoryComputationMetadata;
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct PositionService {
    tle_grpc_client: Arc<TleGrpcClient>,
}

impl PositionService {
    pub const fn new(tle_grpc_client: Arc<TleGrpcClient>) -> Self {
        Self { tle_grpc_client }
    }

    pub async fn get_position(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        compute: &PositionComputation,
    ) -> Result<(SatellitePosition, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let position = Propagator::from_tle(&tle)?.position_at(datetime, compute)?;

        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((position, metadata))
    }

    pub async fn get_look_angles(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(LookAngles, TrajectoryComputationMetadata), ServiceError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let look_angles =
            Propagator::from_tle(&tle)?.look_angles_at(datetime, observer, compute)?;
        
        let metadata = TrajectoryComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((look_angles, metadata))
    }
}
