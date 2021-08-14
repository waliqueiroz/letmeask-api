package models

import (
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Author struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name"`
	Avatar string             `bson:"avatar"`
}

func (a Author) ToDomain() entities.Author {
	return entities.Author{
		ID:     a.ID.Hex(),
		Name:   a.Name,
		Avatar: a.Avatar,
	}
}
