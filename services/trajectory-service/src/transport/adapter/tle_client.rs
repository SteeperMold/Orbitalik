use tonic::transport::Channel;

use crate::astro::models::Tle;
use crate::domain::models::SatelliteIdentifier;
use crate::transport::adapter::tle_client::tle_grpc::GetTleRequest;
use crate::transport::grpc::converters::ToChrono;

use tle_grpc::tle_service_client::TleServiceClient;

pub mod tle_grpc {
    tonic::include_proto!("tle");
}

#[derive(Clone)]
pub struct TleGrpcClient {
    inner: TleServiceClient<Channel>,
}

impl TleGrpcClient {
    pub async fn new(endpoint: String) -> Result<Self, tonic::transport::Error> {
        let client = TleServiceClient::connect(endpoint).await?;
        Ok(Self { inner: client })
    }

    pub async fn get_tle(
        &self,
        satellite_identifier: SatelliteIdentifier,
    ) -> Result<Tle, tonic::Status> {
        let mut client = self.inner.clone();

        let request = GetTleRequest {
            identifier: Some(satellite_identifier.clone().into()),
        };
        let response = client.get_tle(request).await?.into_inner();

        let t = response.tle.ok_or_else(|| {
            tonic::Status::not_found(format!("TLE not found for {satellite_identifier}"))
        })?;

        let epoch = t
            .epoch
            .as_ref()
            .ok_or_else(|| tonic::Status::invalid_argument("Missing epoch"))?
            .to_chrono()?;

        Ok(Tle {
            norad_id: t.norad_id,
            satellite_name: t.satellite_name,
            line1: t.line1,
            line2: t.line2,
            epoch,
        })
    }
}
