use chrono::{DateTime, Utc};
use std::sync::Arc;

use crate::astro::models::SatellitePosition;
use crate::astro::propagator::Propagator;
use crate::domain::errors::PropagationError;
use crate::domain::models::SatelliteIdentifier;
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct PositionService {
    tle_grpc_client: Arc<TleGrpcClient>,
}

impl PositionService {
    pub const fn new(tle_grpc_client: Arc<TleGrpcClient>) -> Self {
        Self { tle_grpc_client }
    }

    pub async fn calculate_position(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
    ) -> Result<SatellitePosition, PropagationError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let propagator = Propagator::from_tle(&tle)?;
        propagator.position_at(datetime)
    }
}
