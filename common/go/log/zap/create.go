package zaplog

import (
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"go.uber.org/zap"
)

func NewLogger(appEnv string) (applog.Logger, error) {
	var logger *zap.Logger
	var err error

	switch appEnv {
	case "production", "test":
		logger, err = zap.NewProduction()
	case "development":
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{l: logger}, nil
}
