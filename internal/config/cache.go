package config

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper) *redis.Client {
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetInt("REDIS_PORT")
	redisPassword := viper.GetString("REDIS_PASSWORD")
	redisDB := viper.GetInt("REDIS_DB")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})

	return client
}
