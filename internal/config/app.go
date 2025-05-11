package config

import (
	"database/sql"
	"vrs-api/internal/delivery/rest"
	"vrs-api/internal/delivery/rest/route"
	postgressql "vrs-api/internal/repository"
	"vrs-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       *sql.DB
	App      *gin.Engine
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := postgressql.NewUserRepository(config.DB)

	// setup use cases
	userUseCase := usecase.NewUsersUsecase(userRepository)

	// setup controller
	userController := rest.NewUserController(config.App, userUseCase)

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
	}
	routeConfig.Setup()
}