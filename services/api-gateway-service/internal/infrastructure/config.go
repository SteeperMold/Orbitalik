package infrastructure

import (
	"log"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/config"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	HTTPPort        string
	ContextTimeout  time.Duration
	MaxPageSize     int
	DefaultPageSize int

	MetadataAddress   string
	TLEAddress        string
	TrajectoryAddress string
	AuthAddress       string

	UseTLS bool
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv:          config.GetEnv("APP_ENV", "development"),
		HTTPPort:        config.GetEnv("HTTP_PORT", "8080"),
		ContextTimeout:  config.GetEnvAsDuration("CONTEXT_TIMEOUT", 10*time.Second),
		MaxPageSize:     config.GetEnvAsInt("MAX_PAGE_SIZE", 500),
		DefaultPageSize: config.GetEnvAsInt("DEFAULT_PAGE_SIZE", 50),

		MetadataAddress:   config.GetEnv("METADATA_ADDRESS", "metadata-service:50051"),
		TLEAddress:        config.GetEnv("TLE_ADDRESS", "tle-ingestion-service:50051"),
		TrajectoryAddress: config.GetEnv("TRAJECTORY_ADDRESS", "trajectory-service:50051"),
		AuthAddress:       config.GetEnv("AUTH_ADDRESS", "auth-service:50051"),

		UseTLS: config.GetEnvAsBool("USE_TLS", false),
	}
}
