package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	infra "github.com/waliqueiroz/letmeask-api/infrastructure/errors"
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

	return repository.FindByID(user.ID)
}

func (repository *UserRepository) FindByID(userID string) (entities.User, error) {
	ctx := context.TODO()
	filter := bson.M{"_id": userID}

	result := repository.userCollection.FindOne(ctx, filter)

	var user entities.User

	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, infra.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.User{}, err
	}

	return user, nil
}

func (repository *UserRepository) Delete(userID string) error {
	ctx := context.TODO()
	filter := bson.M{"_id": userID}

	_, err := repository.userCollection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) Update(userID string, user entities.User) (entities.User, error) {
	ctx := context.TODO()
	filter := bson.M{"_id": userID}

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": time.Now(),
		},
	}

	_, err := repository.userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, infra.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.User{}, err
	}

	return repository.FindByID(userID)
}

func (repository *UserRepository) UpdatePassword(userID string, password string) error {
	ctx := context.TODO()
	filter := bson.M{"_id": userID}

	update := bson.M{
		"$set": bson.M{
			"password":   password,
			"updated_at": time.Now(),
		},
	}

	_, err := repository.userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return infra.NewResourceNotFoundError("usuário não encontrado")
		}
		return err
	}

	return nil
}

func (repository *UserRepository) FindByEmail(email string) (entities.User, error) {
	ctx := context.TODO()
	filter := bson.M{"email": email}

	result := repository.userCollection.FindOne(ctx, filter)

	var user entities.User

	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, infra.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.User{}, err
	}

	return user, nil
}
