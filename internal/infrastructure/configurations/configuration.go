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
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBDatabase: os.Getenv("DB_DATABASE"),
		},
		Auth: Auth{
			SecretKey: os.Getenv("SECRET_KEY"),
		},
	}

	return configuration, nil
}
