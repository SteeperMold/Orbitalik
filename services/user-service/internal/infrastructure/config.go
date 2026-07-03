package infrastructure

import (
	"log"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/config"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	HTTPPort       string
	GRPCPort       string
	ContextTimeout time.Duration
	DB             *DBConfig
}

type DBConfig struct {
	Host              string
	Port              string
	Name              string
	User              string
	Password          string
	ConnectionTimeout time.Duration
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env filesource, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv:         config.GetEnv("APP_ENV", "development"),
		HTTPPort:       config.GetEnv("HTTP_PORT", "8080"),
		GRPCPort:       config.GetEnv("GRPC_PORT", "50055"),
		ContextTimeout: config.GetEnvAsDuration("CONTEXT_TIMEOUT", 10*time.Second),

		DB: &DBConfig{
			Host:              config.GetEnv("DB_HOST", "user-service"),
			Port:              config.GetEnv("DB_PORT", "5432"),
			Name:              config.GetEnv("DB_NAME", "user_db"),
			User:              config.GetEnv("DB_USER", "user"),
			Password:          config.GetEnv("DB_PASSWORD", "123456789admin"),
			ConnectionTimeout: config.GetEnvAsDuration("DB_CONNECTION_TIMEOUT", 10*time.Second),
		},
	}
}
