use dotenv::dotenv;

pub struct AppConfig {
    pub app_env: String,
    pub http_port: u16,
    pub grpc_port: u16,
    pub tle_service_address: String,
}

impl AppConfig {
    pub fn from_dotenv() -> Self {
        dotenv().ok();
        Self {
            app_env: env_string("APP_ENV", "development"),
            http_port: env_u16("HTTP_PORT", 8080),
            grpc_port: env_u16("GRPC_PORT", 50051),
            tle_service_address: env_string(
                "TLE_SERVICE_ADDRESS",
                "grpc://tle-ingestion-service:50051",
            ),
        }
    }
}

fn env_string(key: &str, default: &str) -> String {
    std::env::var(key).unwrap_or_else(|_| default.to_string())
}

fn env_u16(key: &str, default: u16) -> u16 {
    std::env::var(key)
        .ok()
        .and_then(|v| v.parse::<u16>().ok())
        .unwrap_or(default)
}
