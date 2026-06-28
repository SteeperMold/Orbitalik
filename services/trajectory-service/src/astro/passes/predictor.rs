use rayon::prelude::*;
use std::sync::atomic::AtomicUsize;
use uom::si::f64::Angle;

use crate::astro::coords::geodetic::Geodetic;
use crate::astro::errors::PropagationError;
use crate::astro::models::{Pass, TimeRange};
use crate::astro::propagation::propagator::Propagator;

pub struct PassPredictor;

pub struct PassPredictionOptions<'a> {
    pub range: TimeRange,
    pub observer: &'a Geodetic,
    pub min_elevation: Angle,
    pub min_peak_elevation: Angle,
    pub max_results: Option<usize>,
}

impl PassPredictor {
    pub fn predict_many(
        satellites: &[Propagator],
        options: &PassPredictionOptions,
    ) -> Result<Vec<Pass>, PropagationError> {
        let remaining = AtomicUsize::new(options.max_results.unwrap_or(usize::MAX));

        let mut all: Vec<Pass> = satellites
            .par_iter()
            .map(|propagator| propagator.predict_passes(options, &remaining))
            .collect::<Result<Vec<_>, _>>()?
            .into_par_iter()
            .flatten()
            .collect();

        all.sort_by_key(|p| p.aos);

        Ok(all)
    }
}
