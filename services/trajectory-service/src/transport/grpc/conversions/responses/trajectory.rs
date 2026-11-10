use tonic::Status;

use crate::astro::models::{ObserverTrajectory, Trajectory};
use crate::domain::models::{TrajectoryComputationMetadata, UnitContext};
use crate::transport::grpc::conversions::timestamp::ToProtoTimestamp;
use crate::transport::grpc::server::trajectory_grpc;

impl trajectory_grpc::TrajectoryComputationMetadata {
    pub fn from_domain(
        metadata: TrajectoryComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            propagation_model: metadata.propagation_model,
            computation_time: Some(metadata.computation_time.to_proto_timestamp()?),
            norad_id: metadata.norad_id,
            satellite_name: metadata.satellite_name,
            tle_epoch: Some(metadata.tle_epoch.to_proto_timestamp()?),
            units: Some(unit_context.into()),
        })
    }
}

impl trajectory_grpc::TrajectoryResponse {
    pub fn build(
        trajectory: &Trajectory,
        metadata: TrajectoryComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: Some(trajectory_grpc::TrajectoryComputationMetadata::from_domain(
                metadata,
                unit_context,
            )?),

            states: trajectory
                .samples
                .iter()
                .map(|p| trajectory_grpc::StateVector::from_domain(p, unit_context))
                .collect::<Result<Vec<_>, Status>>()?,

            range: Some(trajectory_grpc::TimeRange {
                start: Some(trajectory.start.to_proto_timestamp()?),
                end: Some(trajectory.end.to_proto_timestamp()?),
            }),

            sampling: Some(trajectory_grpc::SamplingOptions {
                step_seconds: trajectory.step_seconds,
            }),
        })
    }
}

impl trajectory_grpc::ObserverTrajectoryResponse {
    pub fn build(
        observer_trajectory: &ObserverTrajectory,
        metadata: TrajectoryComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            metadata: Some(trajectory_grpc::TrajectoryComputationMetadata::from_domain(
                metadata,
                unit_context,
            )?),

            points: observer_trajectory
                .samples
                .iter()
                .map(|p| {
                    trajectory_grpc::ObserverTrajectoryPoint::from_look_angles(p, unit_context)
                })
                .collect::<Result<Vec<_>, Status>>()?,

            range: Some(trajectory_grpc::TimeRange {
                start: Some(observer_trajectory.start.to_proto_timestamp()?),
                end: Some(observer_trajectory.end.to_proto_timestamp()?),
            }),

            sampling: Some(trajectory_grpc::SamplingOptions {
                step_seconds: observer_trajectory.step_seconds,
            }),
        })
    }
}
