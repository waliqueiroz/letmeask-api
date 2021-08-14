package configurations

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Database Database
	Auth     Auth
}

func Load() Configuration {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
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

	return configuration
}
