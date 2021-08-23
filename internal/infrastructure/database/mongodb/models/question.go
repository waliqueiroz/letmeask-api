package models

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Content       string             `bson:"content"`
	IsHighlighted bool               `bson:"is_highlighted"`
	IsAnswered    bool               `bson:"is_answered"`
	Author        Author             `bson:"author"`
	Likes         []Like             `bson:"likes,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
}

func (q Question) ToDomain() entities.Question {
	var likes []entities.Like
	for _, like := range q.Likes {
		likes = append(likes, like.ToDomain())
	}

	return entities.Question{
		ID:            q.ID.Hex(),
		Content:       q.Content,
		IsHighlighted: q.IsHighlighted,
		IsAnswered:    q.IsAnswered,
		Author:        q.Author.ToDomain(),
		Likes:         likes,
		CreatedAt:     q.CreatedAt,
	}
}
