package infrastructure

import (
	"log"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/config"
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	HTTPPort           string
	GRPCPort           string
	ContextTimeout     time.Duration
	TLESourceUrl       string
	TLEFetchInterval   time.Duration
	TLEFetchTimeout    time.Duration
	TLEFetchMaxRetries int
	TLECleanupInterval time.Duration
	TLERetentionPeriod time.Duration
	TLECleanupTimeout  time.Duration
	DB                 *postgres.Config
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env filesource, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv:             config.GetEnv("APP_ENV", "development"),
		HTTPPort:           config.GetEnv("HTTP_PORT", "8080"),
		GRPCPort:           config.GetEnv("GRPC_PORT", "50051"),
		ContextTimeout:     config.GetEnvAsDuration("CONTEXT_TIMEOUT", 10*time.Second),
		TLESourceUrl:       config.GetEnv("TLE_SOURCE_URL", "https://celestrak.org/NORAD/elements/gp.php?GROUP=active&FORMAT=tle"),
		TLEFetchInterval:   config.GetEnvAsDuration("TLE_FETCH_INTERVAL", 6*time.Hour),
		TLEFetchTimeout:    config.GetEnvAsDuration("TLE_FETCH_TIMEOUT", 60*time.Second),
		TLEFetchMaxRetries: config.GetEnvAsInt("TLE_FETCH_MAX_RETRIES", 10),
		TLECleanupInterval: config.GetEnvAsDuration("TLE_CLEANUP_INTERVAL", 24*time.Hour),
		TLERetentionPeriod: config.GetEnvAsDuration("TLE_RETENTION_PERIOD", 31*24*time.Hour),
		TLECleanupTimeout:  config.GetEnvAsDuration("TLE_CLEANUP_TIMEOUT", 5*time.Minute),

		DB: &postgres.Config{
			Host:              config.GetEnv("DB_HOST", "metadata-service"),
			Port:              config.GetEnv("DB_PORT", "5432"),
			Name:              config.GetEnv("DB_NAME", "metadata_db"),
			User:              config.GetEnv("DB_USER", "user"),
			Password:          config.GetEnv("DB_PASSWORD", "123456789admin"),
			ConnectionTimeout: config.GetEnvAsDuration("DB_CONNECTION_TIMEOUT", 10*time.Second),
		},
	}
}
