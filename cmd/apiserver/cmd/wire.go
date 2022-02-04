//+build wireinject

package cmd

import (
	"github.com/google/wire"
	"net/http"
	"github.com/longsolong/flow/pkg/http/rest"
	"github.com/longsolong/flow/pkg/infra"
	"github.com/longsolong/flow/pkg/setting"
)

func InitializeServer() (*http.Server, error) {
	panic(wire.Build(rest.ProvideServer, rest.ProvideRouter, infra.ProvideLogger, setting.DefaultSuperSet))
}