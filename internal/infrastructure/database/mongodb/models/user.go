package models

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Avatar    string             `bson:"avatar"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (u User) ToDomain() entities.User {
	return entities.User{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
