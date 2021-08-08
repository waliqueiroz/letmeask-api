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
			return entities.Room{}, application.NewResourceNotFoundError("usuário não encontrado")
		}
		return entities.Room{}, err
	}

	return room, nil
}
