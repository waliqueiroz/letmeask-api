package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/database/mongodb/models"
	infrastructure "github.com/waliqueiroz/letmeask-api/internal/infrastructure/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository struct {
	roomCollection *mongo.Collection
}

func NewRoomRepository(db *mongo.Database) *RoomRepository {
	return &RoomRepository{
		roomCollection: db.Collection("rooms"),
	}
}

func (repository *RoomRepository) Create(room entities.Room) (entities.Room, error) {
	authorID, err := primitive.ObjectIDFromHex(room.Author.ID)
	if err != nil {
		return entities.Room{}, err
	}

	newRoom := models.Room{
		Title: room.Title,
		Author: models.Author{
			ID:     authorID,
			Name:   room.Author.Name,
			Avatar: room.Author.Avatar,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := repository.roomCollection.InsertOne(context.Background(), newRoom)
	if err != nil {
		return entities.Room{}, err
	}

	objectID := result.InsertedID.(primitive.ObjectID)

	return repository.FindByID(objectID.Hex())
}

func (repository *RoomRepository) FindByID(roomID string) (entities.Room, error) {
	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	filter := bson.M{"_id": id}

	result := repository.roomCollection.FindOne(context.Background(), filter)

	var room models.Room

	if err := result.Decode(&room); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Room{}, infrastructure.NewResourceNotFoundError("sala não encontrada")
		}
		return entities.Room{}, err
	}

	return room.ToDomain(), nil
}

func (repository *RoomRepository) Update(roomID string, room entities.Room) (entities.Room, error) {
	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	filter := bson.M{"_id": id}

	questions, err := repository.entityQuestionsToModelQuestions(room.Questions)
	if err != nil {
		return entities.Room{}, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":      room.Title,
			"questions":  questions,
			"updated_at": time.Now(),
		},
	}

	if room.EndedAt != nil {
		update["ended_at"] = room.EndedAt
	}

	_, err = repository.roomCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Room{}, infrastructure.NewResourceNotFoundError("sala não encontrada")
		}
		return entities.Room{}, err
	}

	return repository.FindByID(roomID)
}
func (repository *RoomRepository) entityQuestionsToModelQuestions(entityQuestions []entities.Question) ([]models.Question, error) {
	var questions []models.Question

	for _, entityQuestion := range entityQuestions {

		authorID, err := primitive.ObjectIDFromHex(entityQuestion.Author.ID)
		if err != nil {
			fmt.Println("Deu pau aqui")
			return []models.Question{}, err
		}

		likes, err := repository.entityLikesToModelLikes(entityQuestion.Likes)
		if err != nil {
			return []models.Question{}, err
		}

		question := models.Question{
			Content:       entityQuestion.Content,
			IsHighlighted: entityQuestion.IsHighlighted,
			IsAnswered:    entityQuestion.IsAnswered,
			Author: models.Author{
				ID:     authorID,
				Name:   entityQuestion.Author.Name,
				Avatar: entityQuestion.Author.Avatar,
			},
			Likes:     likes,
			CreatedAt: entityQuestion.CreatedAt,
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (repository *RoomRepository) entityLikesToModelLikes(entityLikes []entities.Like) ([]models.Like, error) {
	var likes []models.Like
	for _, entityLike := range entityLikes {

		authorID, err := primitive.ObjectIDFromHex(entityLike.Author.ID)
		if err != nil {
			fmt.Println("Deu pau aqui")
			return []models.Like{}, err
		}

		like := models.Like{
			Author: models.Author{
				ID:     authorID,
				Name:   entityLike.Author.Name,
				Avatar: entityLike.Author.Avatar,
			},
			CreatedAt: entityLike.CreatedAt,
		}

		likes = append(likes, like)
	}

	return likes, nil
}
