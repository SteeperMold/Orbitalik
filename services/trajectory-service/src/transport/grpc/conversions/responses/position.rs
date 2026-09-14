use tonic::Status;

use crate::astro::models::{LookAngles, SatellitePosition};
use crate::domain::models::{TrajectoryComputationMetadata, UnitContext};
use crate::transport::grpc::conversions::timestamp::ToProtoTimestamp;
use crate::transport::grpc::server::trajectory_grpc;

impl trajectory_grpc::PositionResponse {
    pub fn build(
        position: &SatellitePosition,
        metadata: TrajectoryComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            time: Some(position.time.to_proto_timestamp()?),

            metadata: Some(trajectory_grpc::TrajectoryComputationMetadata::from_domain(
                metadata,
                unit_context,
            )?),

            teme: position
                .teme
                .as_ref()
                .map(|teme| trajectory_grpc::Vector3::from_xyz(teme, unit_context))
                .transpose()?,

            ecef: position
                .ecef
                .as_ref()
                .map(|ecef| trajectory_grpc::Vector3::from_xyz(ecef, unit_context))
                .transpose()?,

            geodetic: position
                .geodetic
                .as_ref()
                .map(|geodetic| {
                    trajectory_grpc::GeodeticOutput::from_domain(geodetic, unit_context)
                })
                .transpose()?,
        })
    }
}

impl trajectory_grpc::LookAnglesResponse {
    pub fn build(
        look_angles: &LookAngles,
        metadata: TrajectoryComputationMetadata,
        unit_context: &UnitContext,
    ) -> Result<Self, Status> {
        Ok(Self {
            time: Some(look_angles.time.to_proto_timestamp()?),

            metadata: Some(trajectory_grpc::TrajectoryComputationMetadata::from_domain(
                metadata,
                unit_context,
            )?),

            azimuth: look_angles
                .azimuth
                .map(|az| unit_context.angle(az))
                .transpose()?,

            elevation: look_angles
                .elevation
                .map(|el| unit_context.angle(el))
                .transpose()?,

            range: look_angles
                .range
                .map(|r| unit_context.distance(r))
                .transpose()?,
        })
    }
}
