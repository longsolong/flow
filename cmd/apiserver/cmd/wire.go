//+build wireinject

package cmd

import (
	"github.com/google/wire"
	"net/http"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/http/rest"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/infra"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/setting"
)

func InitializeServer() (*http.Server, error) {
	panic(wire.Build(rest.ProvideServer, rest.ProvideRouter, infra.ProvideLogger, setting.DefaultSuperSet))
}