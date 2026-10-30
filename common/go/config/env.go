package config

import (
	"os"
	"slices"
	"strconv"
	"strings"
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

func GetEnvAsBool(name string, defaultVal bool) bool {
	valueStr := GetEnv(name, "")
	if valueStr == "" {
		return defaultVal
	}

	truthyValues := []string{
		"true",
		"yes",
		"1",
	}

	return slices.Contains(truthyValues, strings.ToLower(valueStr))
}
