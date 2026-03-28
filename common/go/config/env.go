package config

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key string, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	return value
}

func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}

	return value
}

func GetEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valueStr := GetEnv(name, "")
	if valueStr == "" {
		return defaultVal
	}

	dur, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultVal
	}
	return dur
}
