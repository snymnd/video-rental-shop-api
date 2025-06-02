package dbcommand

import (
	"fmt"

	"github.com/spf13/viper"
)

func newViper() *viper.Viper {
	config := viper.New()

	config.SetConfigFile(".env")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return config
}
