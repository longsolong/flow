package cmd

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/http/rest/route"
)

func bootstrap(srv *http.Server) {
	routes := route.InitializeRoutes(srv.Handler.(*gin.Engine))
	routes.Setup()
}
