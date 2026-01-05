use tracing_subscriber::{EnvFilter, fmt};

use crate::domain::errors::StartupError;

pub fn init_logger(app_env: &str) -> Result<(), StartupError> {
    match app_env {
        "development" | "production" | "test" => {}
        _ => return Err(StartupError::InvalidAppEnv(app_env.to_string())),
    }

    let filter = EnvFilter::try_from_default_env().unwrap_or_else(|_| match app_env {
        "production" | "test" => EnvFilter::new("info"),
        "development" => EnvFilter::new("debug"),
        _ => unreachable!(),
    });

    match app_env {
        "production" | "test" => fmt()
            .with_env_filter(filter)
            .json()
            .with_target(false)
            .with_level(true)
            .with_thread_names(true)
            .with_timer(fmt::time::UtcTime::rfc_3339())
            .with_writer(std::io::stdout)
            .try_init()
            .map_err(|e| StartupError::LoggerInit(e.to_string()))?,
        "development" => fmt()
            .with_env_filter(filter)
            .pretty()
            .with_level(true)
            .with_thread_names(true)
            .with_writer(std::io::stdout)
            .try_init()
            .map_err(|e| StartupError::LoggerInit(e.to_string()))?,
        _ => unreachable!(),
    }

    Ok(())
}
