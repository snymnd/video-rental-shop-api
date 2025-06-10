package dbcommand

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

func newViper() *viper.Viper {
	config := viper.New()

	// Set up config to read from environment variables as well
	config.AutomaticEnv()

	// find .env file with relative path
	var err error
	path := "./"
	maxFolderDepth := 10
	for range maxFolderDepth {
		config.AddConfigPath(path)
		// Find .env file with relative path
		config.SetConfigName(".env")
		config.SetConfigType("env")
		err = config.ReadInConfig()
		if err == nil {
			break
		}
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			path += "../"
			continue
		}
		log.Fatalf("failed to parse config, error: %s", err.Error())
	}
	return config
}
