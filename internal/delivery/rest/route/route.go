package route

import (
	"net/http"
	"vrs-api/internal/constant"
	"vrs-api/internal/delivery/rest"
	"vrs-api/internal/delivery/rest/middleware"
	"vrs-api/internal/dto"
	"vrs-api/internal/repository/postgresql"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	UserController   *rest.UserController
	VideoController  *rest.VideoController
	RentalController *rest.RentalController
	RBACRepository   *postgresql.RBACRepository
	TokenManager     *util.TokenManager
}

func (c *RouteConfig) Setup() {
	c.App.Use(middleware.ErrorMiddleware())

	c.SetupPublicRoute()
	c.SetupPrivateRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	v1 := c.App.Group("/v1")
	v1.GET("/welcome", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, dto.Response{
			Success: true,
			Data:    "Welcome to video rental API",
		})
	})
	v1.POST("/auth/register", c.UserController.Register)
	v1.POST("/auth/login", c.UserController.Login)
}

func (c *RouteConfig) SetupPrivateRoute() {
	v1 := c.App.Group("/v1")
	v1.Use(middleware.AuthenticateMiddleware(c.TokenManager))
	v1.POST("/videos",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_VIDEOS, c.RBACRepository),
		c.VideoController.CreateVideo,
	)
	v1.GET("/videos",
		middleware.AuthorizationMiddleware(constant.PERM_READ_ALL, constant.RSC_VIDEOS, c.RBACRepository),
		c.VideoController.GetVideos,
	)
	v1.POST("/rentals",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_RENTALS, c.RBACRepository),
		c.RentalController.RentVideos,
	)
	v1.POST("/rentals/return",
		middleware.AuthorizationMiddleware(constant.PERM_UPDATE_ALL, constant.RSC_RENTALS, c.RBACRepository),
		c.RentalController.ReturnVideos,
	)
	v1.GET("/private", func(ctx *gin.Context) {
		ctx.JSON(200, dto.Response{
			Success: true,
			Data:    "success to access private route",
		})
	})
	v1.GET("/private/user",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_PAYMENTS, c.RBACRepository),
		func(ctx *gin.Context) {
			ctx.JSON(200, dto.Response{
				Success: true,
				Data:    "success to access user only private route",
			})
		})
	v1.GET("/private/admin",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_VIDEOS, c.RBACRepository),
		func(ctx *gin.Context) {
			ctx.JSON(200, dto.Response{
				Success: true,
				Data:    "success to access admin only private route",
			})
		})
}
