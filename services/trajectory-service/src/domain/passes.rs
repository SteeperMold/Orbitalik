use async_trait::async_trait;

use crate::astro::models::{Pass, SatelliteIdentifier};
use crate::domain::errors::ServiceError;
use crate::domain::models::PassesComputationMetadata;
use crate::service::passes::{GetPassesOptions, NextPassesOptions};

#[cfg_attr(test, mockall::automock)]
#[async_trait]
pub trait PassesServiceApi: Send + Sync {
    async fn get_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &GetPassesOptions,
    ) -> Result<(Vec<Pass>, PassesComputationMetadata), ServiceError>;

    async fn next_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &NextPassesOptions,
    ) -> Result<(Vec<Pass>, PassesComputationMetadata), ServiceError>;
}
