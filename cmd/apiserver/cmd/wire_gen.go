// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/http/rest"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/infra"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/setting"
	"net/http"
)

// Injectors from wire.go:

func InitializeServer() (*http.Server, error) {
	httpServerConfig := setting.ProvideHTTPServerConfig()
	engine := rest.ProvideRouter()
	appConfig := setting.ProvideAppConfig()
	logger := infra.ProvideLogger(appConfig)
	server := rest.ProvideServer(httpServerConfig, engine, logger)
	return server, nil
}
