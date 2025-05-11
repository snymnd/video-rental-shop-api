package route

import (
	"vrs-api/internal/delivery/rest"
	"vrs-api/internal/delivery/rest/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App            *gin.Engine
	UserController *rest.UserController
}

func (c *RouteConfig) Setup() {
	c.App.Use(middleware.ErrorMiddleware())

	c.SetupPublicRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	c.App.POST("/auth/register", c.UserController.Register)
	c.App.POST("/auth/login", c.UserController.Login)
}
