package models

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Author    Author             `bson:"author"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (l Like) ToDomain() entities.Like {
	return entities.Like{
		ID:        l.ID.Hex(),
		Author:    l.Author.ToDomain(),
		CreatedAt: l.CreatedAt,
	}
}
