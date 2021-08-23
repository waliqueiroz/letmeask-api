package models

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Questions []Question         `bson:"questions,omitempty"`
	Author    Author             `bson:"author"`
	EndedAt   *time.Time         `bson:"ended_at,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (r Room) ToDomain() entities.Room {
	var questions []entities.Question
	for _, question := range r.Questions {
		questions = append(questions, question.ToDomain())
	}

	return entities.Room{
		ID:        r.ID.Hex(),
		Title:     r.Title,
		Questions: questions,
		Author:    r.Author.ToDomain(),
		EndedAt:   r.EndedAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
