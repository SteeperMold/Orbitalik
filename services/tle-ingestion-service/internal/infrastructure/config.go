package infrastructure

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	HTTPPort         string
	GRPCPort         string
	ContextTimeout   time.Duration
	TLEFetchInterval time.Duration
	TLEFetchTimeout  time.Duration
	TLESourceUrl     string
	DB               *DBConfig
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
		log.Printf("failed to load .env file, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		HTTPPort:         getEnv("HTTP_PORT", "8080"),
		GRPCPort:         getEnv("GRPC_PORT", "50051"),
		ContextTimeout:   getEnvAsDuration("CONTEXT_TIMEOUT_MS", 10000) * time.Millisecond,
		TLEFetchInterval: getEnvAsDuration("TLE_FETCH_INTERVAL_H", 6) * time.Hour,
		TLEFetchTimeout:  getEnvAsDuration("TLE_FETCH_TIMEOUT_MS", 10000) * time.Millisecond,
		TLESourceUrl:     getEnv("TLE_SOURCE_URL", "https://celestrak.org/NORAD/elements/gp.php?GROUP=active&FORMAT=tle"),

		DB: &DBConfig{
			Host:              getEnv("DB_HOST", "tle-ingestion-service"),
			Port:              getEnv("DB_PORT", "5432"),
			Name:              getEnv("DB_NAME", "tle_database"),
			User:              getEnv("DB_USER", "user"),
			Password:          getEnv("DB_PASSWORD", "123456789admin"),
			ConnectionTimeout: getEnvAsDuration("DB_CONNECTION_TIMEOUT_MS", 10_000) * time.Millisecond,
		},
	}
}

func getEnv(key string, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	return value
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}

	return value
}

func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	value := getEnvAsInt(name, int(defaultVal))
	return time.Duration(value)
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valueStr := getEnv(name, "")
	if valueStr == "" {
		return defaultVal
	}

	split := strings.Split(valueStr, sep)
	result := make([]string, 0, len(split))
	for _, v := range split {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultVal
	}

	return result
}
