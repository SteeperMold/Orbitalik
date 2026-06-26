use rayon::prelude::*;

use crate::astro::errors::PropagationError;
use crate::astro::models::Pass;
use crate::astro::passes::context::SatelliteContext;
use crate::astro::passes::detector::PassPredictionOptions;

pub struct PassPredictor;

impl PassPredictor {
    pub fn predict_many(
        satellites: &[SatelliteContext],
        options: &PassPredictionOptions,
    ) -> Result<Vec<Pass>, PropagationError> {
        let mut all: Vec<Pass> = satellites
            .par_iter()
            .map(|sat| sat.propagator.predict_passes(options))
            .collect::<Result<Vec<_>, _>>()?
            .into_par_iter()
            .flat_map(|v| v)
            .collect();

        all.sort_by_key(|p| p.aos);

        if let Some(max) = options.max_results {
            all.truncate(max);
        }

        Ok(all)
    }
}
