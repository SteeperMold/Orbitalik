use async_trait::async_trait;
use chrono::{Duration, Utc};
use futures::future;
use std::sync::Arc;
use std::time::Instant;
use uom::si::f64::Angle;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::models::{Pass, SatelliteIdentifier, TimeRange};
use crate::astro::passes::predictor::PassPredictionOptions;
use crate::astro::passes::predictor::PassPredictor;
use crate::astro::propagation::propagator::Propagator;
use crate::domain::errors::ServiceError;
use crate::domain::models::PassesComputationMetadata;
use crate::domain::passes::PassesServiceApi;
use crate::domain::tle_provider::TleProvider;

pub struct PassesService {
    tle_grpc_client: Arc<dyn TleProvider>,
    max_satellites: usize,
    next_passes_lookahead: Duration,
}

pub struct GetPassesOptions {
    pub range: TimeRange,
    pub observer: Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub max_results: Option<usize>,
}

pub struct NextPassesOptions {
    pub observer: Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub passes_count: usize,
}

struct PredictionContext {
    propagators: Vec<Propagator>,
    metadata: PassesComputationMetadata,
}

impl PassesService {
    pub const fn new(
        tle_grpc_client: Arc<dyn TleProvider>,
        max_satellites: usize,
        next_passes_lookahead: Duration,
    ) -> Self {
        Self {
            tle_grpc_client,
            max_satellites,
            next_passes_lookahead,
        }
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

        let propagators: Vec<_> = tles
            .iter()
            .map(|tle| {
                norad_ids.push(tle.norad_id);
                satellite_names.push(tle.satellite_name.clone());

                Propagator::from_tle(tle)
            })
            .collect::<Result<_, _>>()?;

        let metadata = PassesComputationMetadata {
            propagation_model: "SGP4".to_string(),
            computation_time: Utc::now(),
            norad_ids,
            satellite_names,
            tle_epoch,
            satellites_evaluated: u32::try_from(propagators.len())?,

            passes_found: 0,
            computation_ms: 0,
        };

        Ok(PredictionContext {
            propagators,
            metadata,
        })
    }
}

#[async_trait]
impl PassesServiceApi for PassesService {
    async fn get_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &GetPassesOptions,
    ) -> Result<(Vec<Pass>, PassesComputationMetadata), ServiceError> {
        let mut ctx = self.prepare_prediction(satellites).await?;

        let astro_options = &PassPredictionOptions {
            range: options.range,
            observer: &options.observer,
            min_elevation: options.min_elevation,
            min_peak_elevation: options.min_peak_elevation,
            max_results: options.max_results,
        };

        let start = Instant::now();

        let passes = PassPredictor::predict_many(&ctx.propagators, astro_options)?;

        ctx.metadata.passes_found = u32::try_from(passes.len())?;
        ctx.metadata.computation_ms = u32::try_from(start.elapsed().as_millis())?;

        Ok((passes, ctx.metadata))
    }

    async fn next_passes(
        &self,
        satellites: Vec<SatelliteIdentifier>,
        options: &NextPassesOptions,
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
            observer: &options.observer,
            min_elevation: options.min_elevation,
            min_peak_elevation: options.min_peak_elevation,
            max_results: Some(options.passes_count),
        };

        let passes = PassPredictor::predict_many(&ctx.propagators, astro_options)?;

        ctx.metadata.passes_found = u32::try_from(passes.len())?;
        ctx.metadata.computation_ms = u32::try_from(start.elapsed().as_millis())?;

        Ok((passes, ctx.metadata))
    }
}

#[allow(clippy::unwrap_used)]
#[cfg(test)]
mod tests {
    use super::*;

    use chrono::{Duration, Utc};
    use std::sync::Arc;
    use uom::si::angle::degree;
    use uom::si::f64::Angle;

    use crate::astro::models::SatelliteIdentifier;
    use crate::astro::test_utils::{test_observer, test_tle};
    use crate::domain::tle_provider::MockTleProvider;

    fn satellite(id: u32) -> SatelliteIdentifier {
        SatelliteIdentifier::NoradId(id)
    }

    fn get_passes_options(observer: Geodetic) -> GetPassesOptions {
        let start = Utc::now();

        GetPassesOptions {
            range: TimeRange {
                start,
                end: start + Duration::hours(2),
            },
            observer,
            min_elevation: Angle::new::<degree>(0.0),
            min_peak_elevation: Angle::new::<degree>(0.0),
            max_results: None,
        }
    }

    fn next_passes_options(observer: Geodetic) -> NextPassesOptions {
        NextPassesOptions {
            observer,
            min_elevation: Angle::new::<degree>(0.0),
            min_peak_elevation: Angle::new::<degree>(0.0),
            passes_count: 5,
        }
    }

    fn service(
        provider: MockTleProvider,
        max_satellites: usize,
        lookahead: Duration,
    ) -> PassesService {
        PassesService::new(Arc::new(provider), max_satellites, lookahead)
    }

    #[tokio::test]
    async fn get_passes_rejects_empty_satellite_list() {
        let provider = MockTleProvider::new();

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let result = service.get_passes(Vec::new(), &options).await;

        assert!(matches!(result, Err(ServiceError::NoSatellites)));
    }

    #[tokio::test]
    async fn next_passes_rejects_empty_satellite_list() {
        let provider = MockTleProvider::new();

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = next_passes_options(observer);

        let result = service.next_passes(Vec::new(), &options).await;

        assert!(matches!(result, Err(ServiceError::NoSatellites)));
    }

    #[tokio::test]
    async fn get_passes_rejects_too_many_satellites() {
        let provider = MockTleProvider::new();

        let service = service(provider, 2, Duration::hours(24));

        let satellites = vec![satellite(1), satellite(2), satellite(3)];

        let observer = test_observer();
        let options = get_passes_options(observer);

        let result = service.get_passes(satellites, &options).await;

        assert!(matches!(
            result,
            Err(ServiceError::TooManySatellites {
                provided: 3,
                max: 2,
            })
        ));
    }

    #[tokio::test]
    async fn next_passes_rejects_too_many_satellites() {
        let provider = MockTleProvider::new();

        let service = service(provider, 2, Duration::hours(24));

        let satellites = vec![satellite(1), satellite(2), satellite(3)];

        let observer = test_observer();
        let options = next_passes_options(observer);

        let result = service.next_passes(satellites, &options).await;

        assert!(matches!(
            result,
            Err(ServiceError::TooManySatellites {
                provided: 3,
                max: 2,
            })
        ));
    }

    #[tokio::test]
    async fn exactly_max_satellites_is_allowed() {
        let tle1 = test_tle();
        let id1 = tle1.norad_id;

        let mut tle2 = test_tle();
        tle2.norad_id += 1;
        let id2 = tle2.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| matches!(id, SatelliteIdentifier::NoradId(id) if *id == id1))
            .return_once(move |_| Ok(tle1));

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| matches!(id, SatelliteIdentifier::NoradId(id) if *id == id2))
            .return_once(move |_| Ok(tle2));

        let service = service(provider, 2, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let result = service
            .get_passes(vec![satellite(id1), satellite(id2)], &options)
            .await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn get_passes_fetches_tle_for_each_satellite() {
        let tle1 = test_tle();
        let id1 = tle1.norad_id;

        let mut tle2 = test_tle();
        tle2.norad_id += 1;
        let id2 = tle2.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| {
                matches!(
                    id,
                    SatelliteIdentifier::NoradId(id) if *id == id1
                )
            })
            .return_once(move |_| Ok(tle1));

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| {
                matches!(
                    id,
                    SatelliteIdentifier::NoradId(id) if *id == id2
                )
            })
            .return_once(move |_| Ok(tle2));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let result = service
            .get_passes(vec![satellite(id1), satellite(id2)], &options)
            .await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn get_passes_propagates_tle_error() {
        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(|_| Err(tonic::Status::not_found("TLE not found")));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let result = service.get_passes(vec![satellite(25544)], &options).await;

        assert!(result.is_err());
    }

    #[tokio::test]
    async fn next_passes_propagates_tle_error() {
        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(|_| Err(tonic::Status::not_found("TLE not found")));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = next_passes_options(observer);

        let result = service.next_passes(vec![satellite(25544)], &options).await;

        assert!(result.is_err());
    }

    #[tokio::test]
    async fn get_passes_populates_metadata_from_tle() {
        let tle = test_tle();

        let expected_norad_id = tle.norad_id;
        let expected_name = tle.satellite_name.clone();
        let expected_epoch = tle.epoch;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| {
                matches!(
                    id,
                    SatelliteIdentifier::NoradId(id)
                        if *id == expected_norad_id
                )
            })
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let (_, metadata) = service
            .get_passes(vec![satellite(expected_norad_id)], &options)
            .await
            .unwrap();

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_ids, vec![expected_norad_id]);
        assert_eq!(metadata.satellite_names, vec![expected_name]);
        assert_eq!(metadata.tle_epoch, expected_epoch);
        assert_eq!(metadata.satellites_evaluated, 1);
        assert_eq!(metadata.passes_found, 0);
    }

    #[tokio::test]
    async fn get_passes_metadata_contains_all_satellites() {
        let tle1 = test_tle();
        let id1 = tle1.norad_id;
        let name1 = tle1.satellite_name.clone();

        let mut tle2 = test_tle();
        tle2.norad_id += 1;
        let id2 = tle2.norad_id;
        let name2 = tle2.satellite_name.clone();

        let epoch = tle1.epoch;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| {
                matches!(
                    id,
                    SatelliteIdentifier::NoradId(id) if *id == id1
                )
            })
            .return_once(move |_| Ok(tle1));

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| {
                matches!(
                    id,
                    SatelliteIdentifier::NoradId(id) if *id == id2
                )
            })
            .return_once(move |_| Ok(tle2));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let (_, metadata) = service
            .get_passes(vec![satellite(id1), satellite(id2)], &options)
            .await
            .unwrap();

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_ids, vec![id1, id2]);
        assert_eq!(metadata.satellite_names, vec![name1, name2]);
        assert_eq!(metadata.tle_epoch, epoch);
        assert_eq!(metadata.satellites_evaluated, 2);
    }

    #[tokio::test]
    async fn next_passes_populates_metadata() {
        let tle = test_tle();

        let expected_id = tle.norad_id;
        let expected_name = tle.satellite_name.clone();
        let expected_epoch = tle.epoch;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |id| {
                matches!(
                    id,
                    SatelliteIdentifier::NoradId(id)
                        if *id == expected_id
                )
            })
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = next_passes_options(observer);

        let (_, metadata) = service
            .next_passes(vec![satellite(expected_id)], &options)
            .await
            .unwrap();

        assert_eq!(metadata.propagation_model, "SGP4");
        assert_eq!(metadata.norad_ids, vec![expected_id]);
        assert_eq!(metadata.satellite_names, vec![expected_name]);
        assert_eq!(metadata.tle_epoch, expected_epoch);
        assert_eq!(metadata.satellites_evaluated, 1);
    }

    #[tokio::test]
    async fn passes_found_matches_number_of_returned_passes() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let (passes, metadata) = service
            .get_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        assert_eq!(metadata.passes_found as usize, passes.len());
    }

    #[tokio::test]
    async fn computation_time_is_between_before_and_after() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let before = Utc::now();

        let (_, metadata) = service
            .get_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        let after = Utc::now();

        assert!(metadata.computation_time >= before);
        assert!(metadata.computation_time <= after);
    }

    #[tokio::test]
    async fn computation_time_ms_is_non_negative() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let (_, metadata) = service
            .get_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        assert!(metadata.computation_ms > 0);
    }

    #[tokio::test]
    async fn get_passes_accepts_max_results_zero() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();

        let mut options = get_passes_options(observer);
        options.max_results = Some(0);

        let (passes, metadata) = service
            .get_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        assert!(passes.is_empty());
        assert_eq!(metadata.passes_found, 0);
    }

    #[tokio::test]
    async fn get_passes_respects_max_results() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::days(7));

        let observer = test_observer();

        let mut options = get_passes_options(observer);
        options.range.end = options.range.start + Duration::days(7);
        options.max_results = Some(1);

        let (passes, metadata) = service
            .get_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        assert!(passes.len() <= 1);
        assert_eq!(metadata.passes_found as usize, passes.len());
    }

    #[tokio::test]
    async fn next_passes_respects_passes_count() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::days(7));

        let observer = test_observer();

        let mut options = next_passes_options(observer);
        options.passes_count = 1;

        let (passes, metadata) = service
            .next_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        assert!(passes.len() <= 1);
        assert_eq!(metadata.passes_found as usize, passes.len());
    }

    #[tokio::test]
    async fn next_passes_with_zero_count_returns_no_passes() {
        let tle = test_tle();
        let id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::days(7));

        let observer = test_observer();

        let mut options = next_passes_options(observer);
        options.passes_count = 0;

        let (passes, metadata) = service
            .next_passes(vec![satellite(id)], &options)
            .await
            .unwrap();

        assert!(passes.is_empty());
        assert_eq!(metadata.passes_found, 0);
    }

    #[tokio::test]
    async fn get_passes_requests_expected_satellite_identifier() {
        let tle = test_tle();
        let expected_id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |identifier| {
                matches!(
                    identifier,
                    SatelliteIdentifier::NoradId(id)
                        if *id == expected_id
                )
            })
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = get_passes_options(observer);

        let result = service
            .get_passes(vec![satellite(expected_id)], &options)
            .await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn next_passes_requests_expected_satellite_identifier() {
        let tle = test_tle();
        let expected_id = tle.norad_id;

        let mut provider = MockTleProvider::new();

        provider
            .expect_get_tle()
            .times(1)
            .withf(move |identifier| {
                matches!(
                    identifier,
                    SatelliteIdentifier::NoradId(id)
                        if *id == expected_id
                )
            })
            .return_once(move |_| Ok(tle));

        let service = service(provider, 10, Duration::hours(24));

        let observer = test_observer();
        let options = next_passes_options(observer);

        let result = service
            .next_passes(vec![satellite(expected_id)], &options)
            .await;

        assert!(result.is_ok());
    }
}
