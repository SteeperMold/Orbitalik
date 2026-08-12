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
	UserServiceAddress string
	DB                 *postgres.Config
	JWT                *JWTConfig
}

type JWTConfig struct {
	PrivateKeyPath  string
	PublicKeyPath   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env filesource, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv:             config.GetEnv("APP_ENV", "development"),
		HTTPPort:           config.GetEnv("HTTP_PORT", "8080"),
		GRPCPort:           config.GetEnv("GRPC_PORT", "50055"),
		ContextTimeout:     config.GetEnvAsDuration("CONTEXT_TIMEOUT", 10*time.Second),
		UserServiceAddress: config.GetEnv("USER_SERVICE_ADDRESS", "user-service:50051"),

		DB: &postgres.Config{
			Host:              config.GetEnv("DB_HOST", "user-service"),
			Port:              config.GetEnv("DB_PORT", "5432"),
			Name:              config.GetEnv("DB_NAME", "auth_db"),
			User:              config.GetEnv("DB_USER", "user"),
			Password:          config.GetEnv("DB_PASSWORD", "123456789admin"),
			ConnectionTimeout: config.GetEnvAsDuration("DB_CONNECTION_TIMEOUT", 10*time.Second),
		},

		JWT: &JWTConfig{
			PrivateKeyPath:  config.GetEnv("JWT_PRIVATE_KEY_PATH", "keys/private.pem"),
			PublicKeyPath:   config.GetEnv("JWT_PUBLIC_KEY_PATH", "keys/public.pem"),
			AccessTokenTTL:  config.GetEnvAsDuration("JWT_ACCESS_TOKEN_TTL", 20*time.Minute),
			RefreshTokenTTL: config.GetEnvAsDuration("JWT_REFRESH_TOKEN_TTL", 24*time.Hour),
		},
	}
}
