use async_trait::async_trait;
use tonic::transport::Channel;

use crate::astro::models::{SatelliteIdentifier, Tle};
use crate::domain::tle_provider::TleProvider;
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
}

#[async_trait]
impl TleProvider for TleGrpcClient {
    async fn get_tle(
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
            .ok_or_else(|| tonic::Status::internal("TLE service returned a TLE without epoch"))?
            .to_chrono()
            .map_err(|_| tonic::Status::internal("TLE service returned an invalid epoch"))?;

        Ok(Tle {
            norad_id: t.norad_id,
            satellite_name: t.satellite_name,
            line1: t.line1,
            line2: t.line2,
            epoch,
        })
    }
}

#[allow(clippy::panic)]
#[cfg(test)]
mod tests {
    use super::*;

    use chrono::{TimeZone, Utc};
    use prost_types::Timestamp;
    use tokio::net::TcpListener;
    use tonic::transport::Server;
    use tonic::{Request, Response, Status};

    use crate::astro::models::SatelliteIdentifier;

    //noinspection RsUnresolvedPath
    use tle_grpc::tle_service_server::{TleService, TleServiceServer};
    use tle_grpc::{
        GetTleRequest, GetTleResponse, ListTlesRequest, ListTlesResponse, Tle as ProtoTle,
    };

    #[derive(Clone)]
    struct MockTleService {
        response: Result<GetTleResponse, Status>,
    }

    #[tonic::async_trait]
    impl TleService for MockTleService {
        async fn get_tle(
            &self,
            _request: Request<GetTleRequest>,
        ) -> Result<Response<GetTleResponse>, Status> {
            match &self.response {
                Ok(response) => Ok(Response::new(response.clone())),
                Err(status) => Err(status.clone()),
            }
        }

        async fn list_tles(
            &self,
            _request: Request<ListTlesRequest>,
        ) -> Result<Response<ListTlesResponse>, Status> {
            panic!("list_tles should not be called by these tests");
        }
    }

    async fn start_mock_server(response: Result<GetTleResponse, Status>) -> String {
        let service = MockTleService { response };

        let listener = TcpListener::bind("127.0.0.1:0")
            .await
            .expect("failed to bind test server");

        let address = listener
            .local_addr()
            .expect("failed to get test server address");

        tokio::spawn(async move {
            Server::builder()
                .add_service(TleServiceServer::new(service))
                .serve_with_incoming(tokio_stream::wrappers::TcpListenerStream::new(listener))
                .await
                .expect("test gRPC server failed");
        });

        format!("http://{address}")
    }

    fn valid_proto_tle() -> ProtoTle {
        ProtoTle {
            norad_id: 25544,
            satellite_name: "ISS (ZARYA)".to_string(),
            line1: "1 25544U 98067A   26001.50000000  .00000000  00000-0  00000-0 0  9999"
                .to_string(),
            line2: "2 25544  51.6000  10.0000 0005000  20.0000  30.0000 15.50000000123456"
                .to_string(),
            epoch: Some(Timestamp {
                seconds: Utc
                    .with_ymd_and_hms(2026, 1, 1, 12, 0, 0)
                    .unwrap()
                    .timestamp(),
                nanos: 0,
            }),
        }
    }

    #[tokio::test]
    async fn get_tle_returns_tle() {
        let expected_epoch = Utc.with_ymd_and_hms(2026, 1, 1, 12, 0, 0).unwrap();

        let endpoint = start_mock_server(Ok(GetTleResponse {
            tle: Some(valid_proto_tle()),
        }))
        .await;

        let client = TleGrpcClient::new(endpoint)
            .await
            .expect("failed to create client");

        let result = client
            .get_tle(SatelliteIdentifier::NoradId(25544))
            .await
            .expect("get_tle failed");

        assert_eq!(result.norad_id, 25544);
        assert_eq!(result.satellite_name, "ISS (ZARYA)");

        assert_eq!(
            result.line1,
            "1 25544U 98067A   26001.50000000  .00000000  00000-0  00000-0 0  9999"
        );

        assert_eq!(
            result.line2,
            "2 25544  51.6000  10.0000 0005000  20.0000  30.0000 15.50000000123456"
        );

        assert_eq!(result.epoch, expected_epoch);
    }

    #[tokio::test]
    async fn get_tle_returns_not_found_when_tle_is_missing() {
        let endpoint = start_mock_server(Ok(GetTleResponse { tle: None })).await;

        let client = TleGrpcClient::new(endpoint)
            .await
            .expect("failed to create client");

        let error = match client.get_tle(SatelliteIdentifier::NoradId(25544)).await {
            Ok(_) => panic!("expected get_tle to fail"),
            Err(error) => error,
        };

        assert_eq!(error.code(), tonic::Code::NotFound);
        assert_eq!(error.message(), "TLE not found for NORAD ID 25544");
    }

    #[tokio::test]
    async fn get_tle_returns_invalid_argument_when_epoch_is_missing() {
        let mut tle = valid_proto_tle();
        tle.epoch = None;

        let endpoint = start_mock_server(Ok(GetTleResponse { tle: Some(tle) })).await;

        let client = TleGrpcClient::new(endpoint)
            .await
            .expect("failed to create client");

        let error = match client.get_tle(SatelliteIdentifier::NoradId(25544)).await {
            Ok(_) => panic!("expected get_tle to fail"),
            Err(error) => error,
        };

        assert_eq!(error.code(), tonic::Code::Internal);
        assert_eq!(error.message(), "TLE service returned a TLE without epoch");
    }

    #[tokio::test]
    async fn get_tle_propagates_grpc_error() {
        let endpoint = start_mock_server(Err(Status::unavailable("TLE service unavailable"))).await;

        let client = TleGrpcClient::new(endpoint)
            .await
            .expect("failed to create client");

        let error = match client.get_tle(SatelliteIdentifier::NoradId(25544)).await {
            Ok(_) => panic!("expected get_tle to fail"),
            Err(error) => error,
        };

        assert_eq!(error.code(), tonic::Code::Unavailable);
        assert_eq!(error.message(), "TLE service unavailable");
    }

    #[tokio::test]
    async fn get_tle_rejects_invalid_epoch() {
        let mut tle = valid_proto_tle();

        tle.epoch = Some(Timestamp {
            seconds: 0,
            nanos: 2_000_000_000,
        });

        let endpoint = start_mock_server(Ok(GetTleResponse { tle: Some(tle) })).await;

        let client = TleGrpcClient::new(endpoint)
            .await
            .expect("failed to create client");

        let error = match client.get_tle(SatelliteIdentifier::NoradId(25544)).await {
            Ok(_) => panic!("expected get_tle to fail"),
            Err(error) => error,
        };

        assert_eq!(error.code(), tonic::Code::Internal);
    }
}
