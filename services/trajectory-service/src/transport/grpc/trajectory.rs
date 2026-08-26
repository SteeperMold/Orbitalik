use tonic::{Request, Response, Status};

use crate::astro::propagation::look_angles::LookAnglesComputation;
use crate::astro::propagation::position::PositionComputation;
use crate::domain::trajectory::TrajectoryServiceApi;
use crate::transport::grpc::server::TrajectoryGrpcServer;
use crate::transport::grpc::server::trajectory_grpc::{
    ObserverTrajectoryRequest, ObserverTrajectoryResponse, TrajectoryRequest, TrajectoryResponse,
};

impl<P, T, Pa> TrajectoryGrpcServer<P, T, Pa>
where
    P: Sync,
    T: TrajectoryServiceApi,
    Pa: Sync,
{
    pub async fn handle_get_trajectory(
        &self,
        request: Request<TrajectoryRequest>,
    ) -> Result<Response<TrajectoryResponse>, Status> {
        let req = request.into_inner();

        let identifier = req
            .identifier
            .ok_or_else(|| Status::invalid_argument("Missing satellite identifier"))?
            .try_into()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let sampling = req
            .sampling
            .ok_or_else(|| Status::invalid_argument("Missing sampling"))?
            .try_into()?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(PositionComputation::default, PositionComputation::from);

        let (trajectory, metadata) = self
            .trajectory
            .get_trajectory(identifier, range, sampling, &compute)
            .await?;

        let response = TrajectoryResponse::from_trajectory(&trajectory, metadata, req.units)?;

        Ok(Response::new(response))
    }

    pub async fn handle_get_observer_trajectory(
        &self,
        request: Request<ObserverTrajectoryRequest>,
    ) -> Result<Response<ObserverTrajectoryResponse>, Status> {
        let req = request.into_inner();

        let identifier = req
            .identifier
            .ok_or_else(|| Status::invalid_argument("Missing satellite identifier"))?
            .try_into()?;

        let range = req
            .range
            .ok_or_else(|| Status::invalid_argument("Missing range"))?
            .try_into()?;

        let sampling = req
            .sampling
            .ok_or_else(|| Status::invalid_argument("Missing sampling"))?
            .try_into()?;

        let observer = req
            .observer
            .ok_or_else(|| Status::invalid_argument("Missing observer"))?
            .try_into()?;

        let mask = req.output_mask.as_ref();
        let compute = mask.map_or_else(LookAnglesComputation::default, LookAnglesComputation::from);

        let (observer_trajectory, metadata) = self
            .trajectory
            .get_observer_trajectory(identifier, range, sampling, &observer, &compute)
            .await?;

        let response = ObserverTrajectoryResponse::from_observer_trajectory(
            &observer_trajectory,
            metadata,
            req.units,
        )?;

        Ok(Response::new(response))
    }
}
