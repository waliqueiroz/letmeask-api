package env

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
)

type EnvProvider struct{}

func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

func (provider *EnvProvider) LoadConfiguration() configurations.Configuration {
	var configuration configurations.Configuration

	if err := env.Parse(&configuration); err != nil {
		panic(err)
	}

	return configuration
}

func (provider *EnvProvider) LoadConfigurationFromFile(path string) configurations.Configuration {
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}

	return provider.LoadConfiguration()
}
