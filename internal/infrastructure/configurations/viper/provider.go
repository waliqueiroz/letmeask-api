package viper

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
)

type ViperProvider struct {
}

func NewViperProvider() *ViperProvider {
	return &ViperProvider{}
}

func (provider *ViperProvider) LoadConfiguration(path string) configurations.Configuration {
	viper.AddConfigPath(path) // root
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	var configuration configurations.Configuration

	if err := viper.Unmarshal(&configuration); err != nil {
		panic(fmt.Errorf("error unmarshaling config file: %w", err))
	}

	return configuration
}
