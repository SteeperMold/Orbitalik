use dotenv::dotenv;
use std::str::FromStr;

pub struct AppConfig {
    pub app_env: String,
    pub http_port: u16,
    pub grpc_port: u16,
    pub tle_service_address: String,
    pub max_satellites: usize,
}

impl AppConfig {
    pub fn from_dotenv() -> Self {
        dotenv().ok();
        Self {
            app_env: parse_env("APP_ENV", "development".to_string()),
            http_port: parse_env("HTTP_PORT", 8080),
            grpc_port: parse_env("GRPC_PORT", 50051),
            tle_service_address: parse_env(
                "TLE_SERVICE_ADDRESS",
                "grpc://tle-ingestion-service:50051".to_string(),
            ),
            max_satellites: parse_env("MAX_SATELLITES", 30),
        }
    }
}

fn parse_env<T>(key: &str, default: T) -> T
where
    T: FromStr,
{
    std::env::var(key)
        .ok()
        .and_then(|v| v.parse::<T>().ok())
        .unwrap_or(default)
}
