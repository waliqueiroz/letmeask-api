package models

import (
	"fmt"
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type Question struct {
	Content       string    `bson:"content"`
	IsHighlighted bool      `bson:"is_highlighted"`
	IsAnswered    bool      `bson:"is_answered"`
	Author        Author    `bson:"author"`
	Likes         []Like    `bson:"likes,omitempty"`
	CreatedAt     time.Time `bson:"created_at"`
}

func (q Question) ToDomain(ID int) entities.Question {
	var likes []entities.Like
	for key, like := range q.Likes {
		likes = append(likes, like.ToDomain(key))
	}

	return entities.Question{
		ID:            fmt.Sprintf("%d", ID),
		Content:       q.Content,
		IsHighlighted: q.IsHighlighted,
		IsAnswered:    q.IsAnswered,
		Author:        q.Author.ToDomain(),
		Likes:         likes,
		CreatedAt:     q.CreatedAt,
	}
}
