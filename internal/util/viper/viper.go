package viper

import (
	"errors"
	"vrs-api/internal/util/logger"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()
	log := logger.GetLogger()

	// find .env file with relative path
	var err error
	path := "./"
	maxFolderDepth := 10
	for range maxFolderDepth {
		viper.AddConfigPath(path)
		viper.SetConfigName(".env")
		err = viper.ReadInConfig()
		if err == nil {
			break
		}

		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			path = path + "/"
			continue
		}

		log.Fatalf("failed to parse config, error: " + err.Error())
	}
	return config
}
