package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	domain "github.com/waliqueiroz/letmeask-api/internal/domain/errors"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/database/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ctx := context.Background()
	result, err := repository.userCollection.Find(ctx, bson.M{})

	if err != nil {
		return []entities.User{}, err
	}

	defer result.Close(ctx)

	var users []entities.User

	for result.Next(ctx) {
		var user models.User

		err := result.Decode(&user)
		if err != nil {
			return []entities.User{}, err
		}

		users = append(users, user.ToDomain())
	}

	return users, nil
}

func (repository *UserRepository) Create(user entities.User) (entities.User, error) {
	newUser := models.User{
		Name:      user.Name,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := repository.userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return entities.User{}, err
	}

	objectID := result.InsertedID.(primitive.ObjectID)

	return repository.FindByID(objectID.Hex())
}

func (repository *UserRepository) FindByID(userID string) (entities.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return entities.User{}, err
	}

	filter := bson.M{"_id": id}

	result := repository.userCollection.FindOne(context.Background(), filter)

	var user models.User

	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, domain.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.User{}, err
	}

	return user.ToDomain(), nil
}

func (repository *UserRepository) Delete(userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	_, err = repository.userCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) Update(userID string, user entities.User) (entities.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return entities.User{}, err
	}

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"updated_at": time.Now(),
		},
	}

	_, err = repository.userCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, domain.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.User{}, err
	}

	return repository.FindByID(userID)
}

func (repository *UserRepository) UpdatePassword(userID string, password string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"password":   password,
			"updated_at": time.Now(),
		},
	}

	_, err = repository.userCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.NewResourceNotFoundError("usuário não encontrado")
		}
		return err
	}

	return nil
}

func (repository *UserRepository) FindByEmail(email string) (entities.User, error) {
	filter := bson.M{"email": email}

	result := repository.userCollection.FindOne(context.Background(), filter)

	var user models.User

	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, domain.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.User{}, err
	}

	return user.ToDomain(), nil
}
