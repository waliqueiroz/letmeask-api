package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		userCollection: db.Collection("users"),
	}
}

func (repository *UserRepository) FindAll() ([]entities.User, error) {
	ctx := context.TODO()
	result, err := repository.userCollection.Find(ctx, bson.M{})

	if err != nil {
		return []entities.User{}, err
	}

	defer result.Close(ctx)

	var users []entities.User

	for result.Next(ctx) {
		var user entities.User

		err := result.Decode(&user)
		if err != nil {
			return []entities.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository *UserRepository) Create(user entities.User) (entities.User, error) {
	ctx := context.TODO()

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := repository.userCollection.InsertOne(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return repository.FindById(user.ID)
}

func (repository *UserRepository) FindById(userID string) (entities.User, error) {
	ctx := context.TODO()
	filter := bson.M{"_id": userID}

	result := repository.userCollection.FindOne(ctx, filter)

	var user entities.User

	err := result.Decode(&user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
