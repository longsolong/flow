package infra

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents a logger
type Logger struct {
	Log *zap.Logger
}

// CreateLogger creates a logger instance for all components
func CreateLogger(level int) (*Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return &Logger{
		Log: logger,
	}, nil
}