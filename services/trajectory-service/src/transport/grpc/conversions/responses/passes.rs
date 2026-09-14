use tonic::Status;

use crate::astro::models::{Pass, TimeRange};
use crate::domain::models::{PassesComputationMetadata, UnitContext};
use crate::transport::grpc::conversions::timestamp::ToProtoTimestamp;
use crate::transport::grpc::server::trajectory_grpc;

impl trajectory_grpc::PassesComputationMetadata {
    pub fn from_domain(
        metadata: PassesComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            propagation_model: metadata.propagation_model,
            computation_time: Some(metadata.computation_time.to_proto_timestamp()?),
            norad_ids: metadata.norad_ids,
            satellite_names: metadata.satellite_names,
            tle_epoch: Some(metadata.tle_epoch.to_proto_timestamp()?),
            satellites_evaluated: metadata.satellites_evaluated,
            passes_found: metadata.passes_found,
            computation_ms: metadata.computation_ms,
            units: Some(unit_context.into()),
        })
    }
}

impl trajectory_grpc::Pass {
    pub fn from_domain(pass: &Pass, units: &UnitContext) -> Result<Self, Status> {
        Ok(Self {
            satellite: Some(pass.satellite.clone().into()),

            aos: Some(pass.aos.to_proto_timestamp()?),
            aos_azimuth: units.angle(pass.aos_azimuth)?,

            max_elevation_time: Some(pass.max_elevation_time.to_proto_timestamp()?),
            max_elevation: units.angle(pass.max_elevation)?,
            max_elevation_azimuth: units.angle(pass.max_elevation_azimuth)?,

            los: Some(pass.los.to_proto_timestamp()?),
            los_azimuth: units.angle(pass.los_azimuth)?,

            duration_seconds: pass.duration_seconds,
        })
    }
}

impl trajectory_grpc::GetPassesResponse {
    pub fn build(
        passes: &[Pass],
        metadata: PassesComputationMetadata,
        range: TimeRange,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: Some(trajectory_grpc::PassesComputationMetadata::from_domain(
                metadata,
                unit_context,
            )?),

            passes: passes
                .iter()
                .map(|p| trajectory_grpc::Pass::from_domain(p, unit_context))
                .collect::<Result<Vec<_>, Status>>()?,

            range: Some(trajectory_grpc::TimeRange {
                start: Some(range.start.to_proto_timestamp()?),
                end: Some(range.end.to_proto_timestamp()?),
            }),
        })
    }
}

impl trajectory_grpc::NextPassesResponse {
    pub fn build(
        passes: &[Pass],
        metadata: PassesComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: Some(trajectory_grpc::PassesComputationMetadata::from_domain(
                metadata,
                unit_context,
            )?),

            passes: passes
                .iter()
                .map(|p| trajectory_grpc::Pass::from_domain(p, unit_context))
                .collect::<Result<Vec<_>, Status>>()?,
        })
    }
}
