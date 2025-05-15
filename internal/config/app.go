package config

import (
	"database/sql"
	"vrs-api/internal/delivery/rest"
	"vrs-api/internal/delivery/rest/route"
	"vrs-api/internal/repository/postgresql"
	"vrs-api/internal/usecase"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB           *sql.DB
	App          *gin.Engine
	TokenManager *util.TokenManager
	Config       *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	userRepository := postgresql.NewUserRepository(config.DB)
	rbacRepository := postgresql.NewRBACRepository(config.DB)
	videoRepository := postgresql.NewVideoRepository(config.DB)

	// setup use cases
	userUseCase := usecase.NewUsersUsecase(userRepository, config.TokenManager)
	videoUsecase := usecase.NewVideoUsecase(videoRepository)

	// setup controller
	userController := rest.NewUserController(config.App, userUseCase)
	videoController := rest.NewVideoController(config.App, videoUsecase)

	routeConfig := route.RouteConfig{
		App:             config.App,
		TokenManager:    config.TokenManager,
		UserController:  userController,
		RBACRepository:  rbacRepository,
		VideoController: videoController,
	}
	routeConfig.Setup()
}
