use chrono::{DateTime, Utc};
use std::sync::Arc;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::look_angles::LookAnglesComputation;
use crate::astro::models::LookAngles;
use crate::astro::propagator::Propagator;
use crate::domain::errors::PropagationError;
use crate::domain::models::{ComputationMetadata, SatelliteIdentifier};
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct LookAnglesService {
    tle_grpc_client: Arc<TleGrpcClient>,
}

impl LookAnglesService {
    pub const fn new(tle_grpc_client: Arc<TleGrpcClient>) -> Self {
        Self { tle_grpc_client }
    }

    pub async fn get_look_angles_with_metadata(
        &self,
        satellite_identifier: SatelliteIdentifier,
        datetime: DateTime<Utc>,
        observer: &Geodetic,
        compute: &LookAnglesComputation,
    ) -> Result<(LookAngles, ComputationMetadata), PropagationError> {
        let tle = self
            .tle_grpc_client
            .get_tle(satellite_identifier.clone())
            .await?;

        let look_angles =
            Propagator::from_tle(&tle)?.look_angles_at(datetime, observer, compute)?;
        let metadata = ComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: datetime,
            norad_id: tle.norad_id,
            satellite_name: tle.satellite_name,
            tle_epoch: tle.epoch,
        };

        Ok((look_angles, metadata))
    }
}
