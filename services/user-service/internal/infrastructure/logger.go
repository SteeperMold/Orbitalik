package infrastructure

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger(appEnv string) *zap.Logger {
	var logger *zap.Logger
	var err error

	switch appEnv {
	case "production", "test":
		logger, err = zap.NewProduction()
	case "development":
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	return logger
}

func LoggerSync(logger *zap.Logger) {
	err := logger.Sync()
	if err != nil {
		log.Printf("logger sync error: %v\n", err)
	}
}
