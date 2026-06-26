use chrono::{Duration, Utc};
use futures::future;
use std::sync::Arc;
use std::time::Instant;
use uom::si::f64::Angle;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{Pass, SatelliteIdentifier, TimeRange};
use crate::astro::passes::context::SatelliteContext;
use crate::astro::passes::predictor::PassPredictionOptions;
use crate::astro::passes::predictor::PassPredictor;
use crate::domain::errors::ServiceError;
use crate::domain::models::PassesComputationMetadata;
use crate::transport::adapter::tle_client::TleGrpcClient;

pub struct PassesService {
    tle_grpc_client: Arc<TleGrpcClient>,
    max_satellites: usize,
    next_passes_lookahead: Duration,
}

pub struct GetPassesOptions<'a> {
    pub range: TimeRange,
    pub observer: &'a Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub max_results: Option<usize>,
}

pub struct NextPassesOptions<'a> {
    pub observer: &'a Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub passes_count: usize,
}

struct PredictionContext {
    contexts: Vec<SatelliteContext>,
    metadata: PassesComputationMetadata,
}

impl PassesService {
    pub const fn new(
        tle_grpc_client: Arc<TleGrpcClient>,
        max_satellites: usize,
        next_passes_lookahead: Duration,
    ) -> Self {
        Self {
            tle_grpc_client,
            max_satellites,
            next_passes_lookahead,
        }
    }

    pub async fn get_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &GetPassesOptions<'_>,
    ) -> Result<(Vec<Pass>, PassesComputationMetadata), ServiceError> {
        let mut ctx = self.prepare_prediction(satellites).await?;

        let astro_options = &PassPredictionOptions {
            range: options.range,
            observer: options.observer,
            min_elevation: options.min_elevation,
            min_peak_elevation: options.min_peak_elevation,
            max_results: options.max_results,
        };

        let start = Instant::now();

        let passes = PassPredictor::predict_many(&ctx.contexts, astro_options)?;

        ctx.metadata.passes_found = passes.len() as u32;
        ctx.metadata.computation_ms = start.elapsed().as_millis() as u32;

        Ok((passes, ctx.metadata))
    }

    pub async fn next_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &NextPassesOptions<'_>,
    ) -> Result<(Vec<Pass>, PassesComputationMetadata), ServiceError> {
        let mut ctx = self.prepare_prediction(satellites).await?;

        let start = Instant::now();
        let now = Utc::now();

        let range = TimeRange {
            start: now,
            end: now + self.next_passes_lookahead,
        };

        let astro_options = &PassPredictionOptions {
            range,
            observer: options.observer,
            min_elevation: options.min_elevation,
            min_peak_elevation: options.min_peak_elevation,
            max_results: Some(options.passes_count),
        };

        let passes = PassPredictor::predict_many(&ctx.contexts, astro_options)?;

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
