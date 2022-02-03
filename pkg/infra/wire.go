//+build wireinject

package infra

import (
	"github.com/google/wire"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/setting"
	"sync"
)

var logger *Logger
var once sync.Once
var sentry_once  sync.Once

// ProvideLogger ...
func ProvideLogger(appConfig *setting.AppConfig) *Logger {
	once.Do(func() {
		logger = MustNewLogger(appConfig)
	})
	return logger
}

// InitializeLogger ...
func InitializeLogger() *Logger {
	panic(wire.Build(ProvideLogger, setting.DefaultSuperSet))
}

// ProvideLoggerSuperSet ...
var ProvideLoggerSuperSet = wire.NewSet(
	ProvideLogger,
	setting.DefaultSuperSet,
)