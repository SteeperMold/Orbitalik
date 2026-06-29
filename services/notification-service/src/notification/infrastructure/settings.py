import enum

import pydantic_settings


class AppEnv(str, enum.Enum):
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
    metadata_service_url: str

    grpc_port: int = 50051
    http_port: int = 8080

    @property
    def database_url(self) -> str:
        return (
            f"postgresql+asyncpg://"
            f"{self.db_user}:{self.db_password}"
            f"@{self.db_host}:{self.db_port}"
            f"/{self.db_name}"
        )
