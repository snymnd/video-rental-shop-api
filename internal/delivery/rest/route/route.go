package route

import (
	"net/http"
	"vrs-api/internal/constant"
	"vrs-api/internal/delivery/rest"
	"vrs-api/internal/delivery/rest/middleware"
	"vrs-api/internal/dto"
	"vrs-api/internal/repository/postgresql"
	"vrs-api/internal/repository/redis"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App                 *gin.Engine
	UserController      *rest.UserController
	VideoController     *rest.VideoController
	RentalController    *rest.RentalController
	RBACRepository      *postgresql.RBACRepository
	RBACCacheRepository *redis.RBACCacheRepository
	TokenManager        *util.TokenManager
	PaymentController   *rest.PaymentController
}

func (c *RouteConfig) Setup() {
	c.App.Use(middleware.ErrorMiddleware())
	c.SetupPublicRoute()
	c.SetupPrivateRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	v1 := c.App.Group("/api/v1")
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
	v1 := c.App.Group("/api/v1")
	v1.Use(middleware.AuthenticateMiddleware(c.TokenManager))
	v1.POST("/videos",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_VIDEOS, c.RBACRepository, c.RBACCacheRepository),
		c.VideoController.CreateVideo,
	)
	v1.GET("/videos",
		middleware.AuthorizationMiddleware(constant.PERM_READ_ALL, constant.RSC_VIDEOS, c.RBACRepository, c.RBACCacheRepository),
		c.VideoController.GetVideos,
	)
	v1.POST("/rentals",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_RENTALS, c.RBACRepository, c.RBACCacheRepository),
		c.RentalController.RentVideos,
	)
	v1.POST("/rentals/return",
		middleware.AuthorizationMiddleware(constant.PERM_UPDATE_ALL, constant.RSC_RENTALS, c.RBACRepository, c.RBACCacheRepository),
		c.RentalController.ReturnVideos,
	)
	// NOTE: This endpoint/operation is intended to be called after payment has been verified by a third-party (i.e. payment-gateway, admin (for cash payment method).
	v1.GET("/payments/rentals/:method/:id",
		middleware.AuthorizationMiddleware(constant.PERM_UPDATE_ALL, constant.RSC_PAYMENTS, c.RBACRepository, c.RBACCacheRepository),
		c.PaymentController.PayRentals,
	)

	// endpoint to check rbac
	v1.GET("/private", func(ctx *gin.Context) {
		ctx.JSON(200, dto.Response{
			Success: true,
			Data:    "success to access private route",
		})
	})
	v1.GET("/private/user",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_PAYMENTS, c.RBACRepository, c.RBACCacheRepository),
		func(ctx *gin.Context) {
			ctx.JSON(200, dto.Response{
				Success: true,
				Data:    "success to access user only private route",
			})
		})
	v1.GET("/private/admin",
		middleware.AuthorizationMiddleware(constant.PERM_CREATE, constant.RSC_VIDEOS, c.RBACRepository, c.RBACCacheRepository),
		func(ctx *gin.Context) {
			ctx.JSON(200, dto.Response{
				Success: true,
				Data:    "success to access admin only private route",
			})
		})
}
