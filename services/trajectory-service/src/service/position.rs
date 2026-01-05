use chrono::{DateTime, Utc};
use std::sync::Arc;

use crate::astro::models::SatellitePosition;
use crate::astro::position::PositionComputation;
use crate::astro::propagator::Propagator;
use crate::domain::errors::PropagationError;
use crate::domain::models::{ComputationMetadata, SatelliteIdentifier};
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct PositionService {
    tle_grpc_client: Arc<TleGrpcClient>,
}

impl PositionService {
    pub const fn new(tle_grpc_client: Arc<TleGrpcClient>) -> Self {
        Self { tle_grpc_client }
    }

    pub async fn get_position_with_metadata(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        compute: &PositionComputation,
    ) -> Result<(SatellitePosition, ComputationMetadata), PropagationError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let position = Propagator::from_tle(&tle)?.position_at(datetime, compute)?;
        let metadata = ComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: datetime,
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((position, metadata))
    }
}
