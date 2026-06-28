use dotenv::dotenv;
use std::str::FromStr;

use crate::domain::errors::StartupError;

pub struct AppConfig {
    pub app_env: String,

    pub http_port: u16,
    pub grpc_port: u16,

    pub tle_service_address: String,

    pub max_satellites: usize,
    pub next_passes_lookahead: chrono::Duration,
}

impl AppConfig {
    pub fn from_dotenv() -> Result<Self, StartupError> {
        dotenv().ok();
        Ok(Self {
            app_env: parse_env("APP_ENV", "development".to_string()),
            http_port: parse_env("HTTP_PORT", 8080),
            grpc_port: parse_env("GRPC_PORT", 50051),
            tle_service_address: parse_env(
                "TLE_SERVICE_ADDRESS",
                "grpc://tle-ingestion-service:50051".to_string(),
            ),
            max_satellites: parse_env("MAX_SATELLITES", 30),
            next_passes_lookahead: parse_duration_env("NEXT_PASSES_LOOKAHEAD", "15d").map_err(
                |_| StartupError::InvalidConfig("invalid next passes lookahead".to_string()),
            )?,
        })
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

fn parse_duration_env(
    key: &str,
    default: &str,
) -> Result<chrono::Duration, Box<dyn std::error::Error>> {
    let raw = std::env::var(key).unwrap_or_else(|_| default.to_string());
    let std_dur: std::time::Duration = humantime::parse_duration(&raw)?;
    Ok(chrono::Duration::from_std(std_dur)?)
}
