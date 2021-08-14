package database

import (
	"context"
	"fmt"

	"github.com/waliqueiroz/letmeask-api/infrastructure/configurations"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(configuration configurations.Configuration) (*mongo.Database, error) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/", configuration.Database.DBHost, configuration.Database.DBPort))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(configuration.Database.DBDatabase), nil
}
