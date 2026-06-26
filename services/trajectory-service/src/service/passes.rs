use chrono::Utc;
use futures::future;
use std::sync::Arc;
use std::time::Instant;

use crate::astro::models::{Pass, SatelliteIdentifier};
use crate::astro::passes::context::SatelliteContext;
use crate::astro::passes::detector::PassPredictionOptions;
use crate::astro::passes::predictor::PassPredictor;
use crate::domain::errors::ServiceError;
use crate::domain::models::PassesComputationMetadata;
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct PassesService {
    tle_grpc_client: Arc<TleGrpcClient>,
    max_satellites: usize,
}

struct PredictionContext {
    contexts: Vec<SatelliteContext>,
    metadata: PassesComputationMetadata,
}

impl PassesService {
    pub const fn new(tle_grpc_client: Arc<TleGrpcClient>, max_satellites: usize) -> Self {
        Self {
            tle_grpc_client,
            max_satellites,
        }
    }

    pub async fn get_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &PassPredictionOptions<'_>,
    ) -> Result<(Vec<Pass>, PassesComputationMetadata), ServiceError> {
        let mut ctx = self.prepare_prediction(satellites).await?;

        let start = Instant::now();

        let passes = PassPredictor::predict_many(&ctx.contexts, options)?;

        ctx.metadata.passes_found = passes.len() as u32;
        ctx.metadata.computation_ms = start.elapsed().as_millis() as u32;

        Ok((passes, ctx.metadata))
    }

    async fn prepare_prediction(
        &self,
        satellites: Vec<SatelliteIdentifier>,
    ) -> Result<PredictionContext, ServiceError> {
        if satellites.is_empty() {
            return Err(ServiceError::NoSatellites);
        }

        if satellites.len() > self.max_satellites {
            return Err(ServiceError::TooManySatellites {
                provided: satellites.len(),
                max: self.max_satellites,
            });
        }

        let tles = future::try_join_all(
            satellites
                .iter()
                .cloned()
                .map(|sat| self.tle_grpc_client.get_tle(sat)),
        )
        .await?;
        let tle_epoch = tles
            .first()
            .map(|tle| tle.epoch)
            .ok_or(ServiceError::NoSatellites)?;

        let mut norad_ids = Vec::with_capacity(tles.len());
        let mut satellite_names = Vec::with_capacity(tles.len());

        let contexts: Vec<_> = tles
            .iter()
            .map(|tle| {
                norad_ids.push(tle.norad_id);
                satellite_names.push(tle.satellite_name.clone());

                SatelliteContext::from_tle(tle)
            })
            .collect::<Result<_, _>>()?;

        let metadata = PassesComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_ids,
            satellite_names,
            tle_epoch,
            satellites_evaluated: contexts.len() as u32,

            passes_found: 0,
            computation_ms: 0,
        };

        Ok(PredictionContext { contexts, metadata })
    }
}
