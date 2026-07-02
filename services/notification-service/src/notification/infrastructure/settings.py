import datetime as dt
import enum

import pydantic_settings


class AppEnv(enum.Enum):
    DEVELOPMENT = "development"
    PRODUCTION = "production"
    TEST = "test"


class Settings(pydantic_settings.BaseSettings):
    model_config = pydantic_settings.SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
    )

    app_env: AppEnv = AppEnv.DEVELOPMENT

    db_host: str
    db_port: int = 5432
    db_user: str
    db_password: str
    db_name: str

    trajectory_service_url: str

    smtp_host: str = ""
    smtp_port: int = 0

    grpc_port: int = 50051
    http_port: int = 8080

    request_timeout_seconds: int = 5

    max_page_size: int = 100
    scheduling_interval_seconds: int = 60
    scheduling_refill_threshold: dt.timedelta = dt.timedelta(hours=1)
    scheduling_window: dt.timedelta = dt.timedelta(hours=2)
    worker_polling_interval_seconds: int = 2

    @property
    def database_url(self) -> str:
        return (
            f"postgresql+asyncpg://"
            f"{self.db_user}:{self.db_password}"
            f"@{self.db_host}:{self.db_port}"
            f"/{self.db_name}"
        )
