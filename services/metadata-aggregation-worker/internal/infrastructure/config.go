package infrastructure

import (
	"log"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/config"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	DB     *DBConfig
	Redis  *RedisConfig
}

type DBConfig struct {
	Host              string
	Port              string
	Name              string
	User              string
	Password          string
	ConnectionTimeout time.Duration
}

type RedisConfig struct {
	Host                    string
	Port                    string
	DirtySatellitesSetKey   string
	DirtySatellitesQueueKey string
	GroupName               string
	ConsumerName            string
	StreamsCount            int
	StreamsBlock            time.Duration
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env, using defaults: %v\n", err)
	}

	return &Config{
		AppEnv: config.GetEnv("APP_ENV", "development"),

		DB: &DBConfig{
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
			GroupName:               config.GetEnv("REDIS_GROUP_NAME", "aggregators"),
			ConsumerName:            config.GetEnv("REDIS_CONSUMER_NAME", "worker-1"),
			StreamsCount:            config.GetEnvAsInt("REDIS_STREAMS_COUNT", 10),
			StreamsBlock:            config.GetEnvAsDuration("REDIS_STREAMS_BLOCK", 5*time.Second),
		},
	}
}
