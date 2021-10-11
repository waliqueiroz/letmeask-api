package configurations

import (
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Database Database
	Auth     Auth
}

func Load() (Configuration, error) {
	if err := godotenv.Load(); err != nil {
		return Configuration{}, err
	}

	configuration := Configuration{
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Database: os.Getenv("DB_DATABASE"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Auth: Auth{
			SecretKey: os.Getenv("SECRET_KEY"),
		},
	}

	return configuration, nil
}
