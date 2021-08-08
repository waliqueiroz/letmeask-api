package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	application "github.com/waliqueiroz/letmeask-api/application/errors"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
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
	ctx := context.TODO()

	room.ID = uuid.New().String()
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	_, err := repository.roomCollection.InsertOne(ctx, room)
	if err != nil {
		return entities.Room{}, err
	}

	return repository.FindByID(room.ID)
}

func (repository *RoomRepository) FindByID(roomID string) (entities.Room, error) {
	ctx := context.TODO()
	filter := bson.M{"_id": roomID}

	result := repository.roomCollection.FindOne(ctx, filter)

	var room entities.Room

	if err := result.Decode(&room); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Room{}, application.NewResourceNotFoundError("sala não encontrada")
		}
		return entities.Room{}, err
	}

	return room, nil
}

func (repository *RoomRepository) Update(roomID string, room entities.Room) (entities.Room, error) {
	ctx := context.TODO()
	filter := bson.M{"_id": roomID}

	update := bson.M{
		"$set": bson.M{
			"title":      room.Title,
			"questions":  room.Questions,
			"author":     room.Author,
			"ended_at":   room.EndedAt,
			"updated_at": time.Now(),
		},
	}

	_, err := repository.roomCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Room{}, application.NewResourceNotFoundError("sala não encontrada")
		}
		return entities.Room{}, err
	}

	return repository.FindByID(roomID)
}
