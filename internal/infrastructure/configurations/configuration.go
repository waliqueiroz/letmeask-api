package configurations

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	Database Database `mapstructure:",squash"`
	Auth     Auth     `mapstructure:",squash"`
}

func NewConfiguration(path string) Configuration {
	viper.AddConfigPath(path) // root
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	var configuration Configuration

	if err := viper.Unmarshal(&configuration); err != nil {
		panic(fmt.Errorf("error unmarshaling config file: %w", err))
	}

	return configuration
}
