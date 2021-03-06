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
	ctx := context.Background()
	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$unwind": bson.M{"path": "$questions", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"questions.created_at": -1}},
		{"$group": bson.M{
			"_id":        "$_id",
			"title":      bson.M{"$first": "$title"},
			"author":     bson.M{"$first": "$author"},
			"created_at": bson.M{"$first": "$created_at"},
			"updated_at": bson.M{"$first": "$updated_at"},
			"ended_at":   bson.M{"$first": "$ended_at"},
			"questions":  bson.M{"$push": "$questions"},
		}},
		{"$addFields": bson.M{
			"questions": bson.M{
				"$filter": bson.M{
					"input": "$questions",
					"as":    "question",
					"cond": bson.M{"$and": []interface{}{
						bson.M{"$ne": []interface{}{"$$question", nil}},
						bson.M{"$ne": []interface{}{"$$question.content", ""}},
					}},
				},
			},
		}},
	}

	result, err := repository.roomCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return entities.Room{}, err
	}

	if result.Next(ctx) {
		var room models.Room

		err := result.Decode(&room)
		if err != nil {
			return entities.Room{}, err
		}

		return room.ToDomain(), nil
	}

	return entities.Room{}, domain.NewResourceNotFoundError("sala n??o encontrada")
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

	fields := bson.M{
		"title":      room.Title,
		"questions":  questions,
		"updated_at": time.Now(),
	}

	if room.EndedAt != nil {
		fields["ended_at"] = room.EndedAt
	}

	update := bson.M{
		"$set": fields,
	}

	_, err = repository.roomCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Room{}, domain.NewResourceNotFoundError("sala n??o encontrada")
		}
		return entities.Room{}, err
	}

	return repository.FindByID(roomID)
}

func (repository *RoomRepository) entityQuestionsToModelQuestions(entityQuestions []entities.Question) ([]models.Question, error) {
	var questions []models.Question

	for _, entityQuestion := range entityQuestions {
		var questionID primitive.ObjectID
		var err error

		if entityQuestion.ID == "" {
			questionID = primitive.NewObjectID()
		} else {
			questionID, err = primitive.ObjectIDFromHex(entityQuestion.ID)
			if err != nil {
				return []models.Question{}, err
			}
		}

		authorID, err := primitive.ObjectIDFromHex(entityQuestion.Author.ID)
		if err != nil {
			return []models.Question{}, err
		}

		likes, err := repository.entityLikesToModelLikes(entityQuestion.Likes)
		if err != nil {
			return []models.Question{}, err
		}

		question := models.Question{
			ID:            questionID,
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
		var likeID primitive.ObjectID
		var err error

		if entityLike.ID == "" {
			likeID = primitive.NewObjectID()
		} else {
			likeID, err = primitive.ObjectIDFromHex(entityLike.ID)
			if err != nil {
				return []models.Like{}, err
			}
		}

		authorID, err := primitive.ObjectIDFromHex(entityLike.Author.ID)
		if err != nil {
			return []models.Like{}, err
		}

		like := models.Like{
			ID: likeID,
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
