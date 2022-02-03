package route

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/authz"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/marvincaspar/go-web-app-boilerplate/pkg/http/rest/handler/response"
)

// APIV1RouterGroup ...
type APIV1RouterGroup struct {
	Gin *gin.Engine
}

// Setup ...
func (r APIV1RouterGroup) Setup() {
	// authorization
	e, err := casbin.NewEnforcer("configs/authz_model.conf", "configs/authz_policy.csv")
	if err != nil {
		panic(err)
	}

	apiv1 := r.Gin.Group("/api/v1", gin.BasicAuth(gin.Accounts{"fakeuser": "fakepasswd"}))
	// middleware
	apiv1.Use(authz.NewAuthorizer(e))
	apiv1.Use(requestid.New())

	// inject middleware
	apiv1.Use(func(c *gin.Context) {
		// before request
		c.Next()
		// after request
	})
	apiv1.GET("/current_user", response.Wrapper(func(c *gin.Context) (gin.H, error) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		return gin.H{"user": user}, nil
	}))
}
