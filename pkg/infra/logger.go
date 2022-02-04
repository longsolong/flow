package infra

import (
	"go.uber.org/zap"
	"github.com/longsolong/flow/pkg/setting"
)

// Logger represents a logger
type Logger struct {
	Zap      *zap.Logger
	SugarZap *zap.SugaredLogger
}

// MustNewLogger creates a logger instance for all components
func MustNewLogger(appConfig *setting.AppConfig) *Logger {
	var logger *zap.Logger
	var err error
	if appConfig.Env == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return &Logger{
		Zap:      logger,
		SugarZap: logger.Sugar(),
	}
}
