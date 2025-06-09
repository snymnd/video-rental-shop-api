package config

import (
	"database/sql"
	"vrs-api/internal/delivery/rest"
	"vrs-api/internal/delivery/rest/route"
	"vrs-api/internal/repository/postgresql"
	redisRepo "vrs-api/internal/repository/redis"
	"vrs-api/internal/usecase"
	"vrs-api/internal/util/logger"
	"vrs-api/internal/util/token"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB           *sql.DB
	Cache        *redis.Client
	App          *gin.Engine
	TokenManager *token.TokenManager
	Logger       *logger.Logger
	Config       *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := postgresql.NewUserRepository(config.DB)
	rbacRepository := postgresql.NewRBACRepository(config.DB)
	videoRepository := postgresql.NewVideoRepository(config.DB)
	rentalRepository := postgresql.NewRentalRepository(config.DB)
	paymentRepository := postgresql.NewPaymentRepository(config.DB)
	txRepository := postgresql.NewTxRepository(config.DB)
	rbacCacheRepository := redisRepo.NewRBACCacheRepository(config.Cache)
	videoCacheRepository := redisRepo.NewVideoCacheRepository(config.Cache)

	// setup use cases
	userUseCase := usecase.NewUsersUsecase(userRepository, config.TokenManager)
	videoUsecase := usecase.NewVideoUsecase(videoRepository, videoCacheRepository)
	rentalUsecase := usecase.NewRentalUsecase(rentalRepository, videoRepository, paymentRepository, txRepository)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository, rentalRepository, txRepository)

	// setup controller
	userController := rest.NewUserController(userUseCase)
	videoController := rest.NewVideoController(videoUsecase)
	rentalController := rest.NewRentalController(rentalUsecase)
	paymentController := rest.NewPaymentController(paymentUsecase)

	routeConfig := route.RouteConfig{
		App:                 config.App,
		TokenManager:        config.TokenManager,
		Logger:              config.Logger,
		UserController:      userController,
		RBACRepository:      rbacRepository,
		RBACCacheRepository: rbacCacheRepository,
		VideoController:     videoController,
		RentalController:    rentalController,
		PaymentController:   paymentController,
	}
	routeConfig.Setup()
}
