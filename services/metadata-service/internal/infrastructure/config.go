package infrastructure

import (
	"log"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/config"
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	HTTPPort          string
	GRPCPort          string
	ContextTimeout    time.Duration
	IngestionInterval time.Duration
	IngestionTimeout  time.Duration
	MaxPageSize       int
	DefaultPageSize   int
	DB                *postgres.Config
	Redis             *RedisConfig
	UCS               *UCSConfig
	Celestrak         *CelestrakConfig
	SatNOGS           *SatNOGSConfig
}

type RedisConfig struct {
	Host                    string
	Port                    string
	DirtySatellitesSetKey   string
	DirtySatellitesQueueKey string
}

type UCSConfig struct {
	SourceURL        string
	FetchTimeout     time.Duration
	FetchRetries     int
	BatchSize        int
	FallbackFilePath string
}

type CelestrakConfig struct {
	SourceURL        string
	FetchTimeout     time.Duration
	FetchRetries     int
	BatchSize        int
	FallbackFilePath string
}

type SatNOGSConfig struct {
	SourceURL    string
	FetchTimeout time.Duration
	FetchRetries int
	BatchSize    int
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv:            config.GetEnv("APP_ENV", "development"),
		HTTPPort:          config.GetEnv("HTTP_PORT", "8080"),
		GRPCPort:          config.GetEnv("GRPC_PORT", "50051"),
		ContextTimeout:    config.GetEnvAsDuration("CONTEXT_TIMEOUT", 10*time.Second),
		IngestionInterval: config.GetEnvAsDuration("INGESTION_INTERVAL", 12*time.Hour),
		IngestionTimeout:  config.GetEnvAsDuration("INGESTION_TIMEOUT", 30*time.Minute),
		MaxPageSize:       config.GetEnvAsInt("MAX_PAGE_SIZE", 500),
		DefaultPageSize:   config.GetEnvAsInt("DEFAULT_PAGE_SIZE", 50),

		DB: &postgres.Config{
			Host:              config.GetEnv("DB_HOST", "postgres"),
			Port:              config.GetEnv("DB_PORT", "5432"),
			Name:              config.GetEnv("DB_NAME", "metadata_db"),
			User:              config.GetEnv("DB_USER", "user"),
			Password:          config.GetEnv("DB_PASSWORD", "123456789admin"),
			ConnectionTimeout: config.GetEnvAsDuration("DB_CONNECTION_TIMEOUT", 10*time.Second),
		},

		Redis: &RedisConfig{
			Host:                    config.GetEnv("REDIS_HOST", "redis"),
			Port:                    config.GetEnv("REDIS_PORT", "6379"),
			DirtySatellitesSetKey:   config.GetEnv("REDIS_DIRTY_SATELLITES_SET_KEY", "dirty-satellites"),
			DirtySatellitesQueueKey: config.GetEnv("REDIS_DIRTY_SATELLITES_QUEUE_KEY", "satellite-dirty"),
		},

		UCS: &UCSConfig{
			SourceURL:        config.GetEnv("UCS_SOURCE_URL", "https://www.ucs.org/sites/default/files/2024-01/UCS-Satellite-Database%205-1-2023%20%28text%29.txt"),
			FetchTimeout:     config.GetEnvAsDuration("UCS_FETCH_TIMEOUT", 3*time.Minute),
			FetchRetries:     config.GetEnvAsInt("UCS_FETCH_RETRIES", 5),
			BatchSize:        config.GetEnvAsInt("UCS_BATCH_SIZE", 1000),
			FallbackFilePath: config.GetEnv("UCS_FALLBACK_FILE_PATH", "assets/fallback/ucs-snapshot-24-04-2026.txt"),
		},

		Celestrak: &CelestrakConfig{
			SourceURL:        config.GetEnv("CELESTRAK_SOURCE_URL", "https://celestrak.org/pub/satcat.txt"),
			FetchTimeout:     config.GetEnvAsDuration("CELESTRAK_FETCH_TIMEOUT", 30*time.Second),
			FetchRetries:     config.GetEnvAsInt("CELESTRAK_FETCH_RETRIES", 5),
			BatchSize:        config.GetEnvAsInt("CELESTRAK_BATCH_SIZE", 1000),
			FallbackFilePath: config.GetEnv("CELESTRAK_FALLBACK_FILE_PATH", "assets/fallback/satcat-snapshot-24-04-2026.txt"),
		},

		SatNOGS: &SatNOGSConfig{
			SourceURL:    config.GetEnv("SATNOGS_SOURCE_URL", "https://db.satnogs.org/api"),
			FetchTimeout: config.GetEnvAsDuration("SATNOGS_FETCH_TIMEOUT", 30*time.Second),
			FetchRetries: config.GetEnvAsInt("SATNOGS_FETCH_RETRIES", 5),
			BatchSize:    config.GetEnvAsInt("SATNOGS_BATCH_SIZE", 1000),
		},
	}
}
