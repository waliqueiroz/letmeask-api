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

	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		configuration.Database.Username,
		configuration.Database.Password,
		configuration.Database.Host,
		configuration.Database.Port,
	)

	clientOptions := options.Client().ApplyURI(dbURI)

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
