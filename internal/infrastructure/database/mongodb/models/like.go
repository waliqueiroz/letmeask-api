package models

import (
	"fmt"
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type Like struct {
	Author    Author    `bson:"author"`
	CreatedAt time.Time `bson:"created_at"`
}

func (l Like) ToDomain(ID int) entities.Like {
	return entities.Like{
		ID:        fmt.Sprintf("%d", ID),
		Author:    l.Author.ToDomain(),
		CreatedAt: l.CreatedAt,
	}
}
