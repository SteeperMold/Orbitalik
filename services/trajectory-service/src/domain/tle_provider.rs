use async_trait::async_trait;

use crate::astro::models::{SatelliteIdentifier, Tle};

#[cfg_attr(test, mockall::automock)]
#[async_trait]
pub trait TleProvider: Send + Sync {
    async fn get_tle(
        &self,
        satellite_identifier: SatelliteIdentifier,
    ) -> Result<Tle, tonic::Status>;
}
