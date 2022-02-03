package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/infra"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/setting"
	"time"
	"net/http"
)

// ProvideRouter create a new gin router
func ProvideRouter() *gin.Engine {
	return gin.Default()
}

// ProvideServer create a new http server
func ProvideServer(httpServerConfig *setting.HTTPServerConfig, router *gin.Engine, logger *infra.Logger) *http.Server {
	logger.SugarZap.Infof("Server listen on %s", httpServerConfig.ServerAddr)
	server := &http.Server{
		Addr:         httpServerConfig.ServerAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}

	return server
}

