use chrono::{DateTime, Utc};

pub struct TrajectoryComputationMetadata {
    pub propagation_model: String,

    pub computation_time: DateTime<Utc>,

    pub norad_id: u32,
    pub satellite_name: String,

    pub tle_epoch: DateTime<Utc>,
}

pub struct PassesComputationMetadata {
    pub propagation_model: String,

    pub computation_time: DateTime<Utc>,

    pub norad_ids: Vec<u32>,
    pub satellite_names: Vec<String>,

    pub tle_epoch: DateTime<Utc>,

    pub satellites_evaluated: u32,
    pub passes_found: u32,

    pub computation_ms: u32,
}

