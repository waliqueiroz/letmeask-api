package mongodb

import (
	"context"
	"fmt"

	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(configuration configurations.Configuration) (*mongo.Database, error) {
	ctx := context.Background()

	credential := options.Credential{
		Username:      configuration.Database.Username,
		Password:      configuration.Database.Password,
		AuthMechanism: "SCRAM-SHA-1",
	}

	dbURI := fmt.Sprintf("mongodb://%s:%s/", configuration.Database.Host, configuration.Database.Port)

	clientOptions := options.Client().ApplyURI(dbURI).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(configuration.Database.Database), nil
}
