package models

import (
	"time"
)

type Room struct {
	ID        string     `bson:"_id"`
	Title     string     `bson:"title"`
	Questions []Question `bson:"questions,omitempty"`
	Author    Author     `bson:"author"`
	EndedAt   *time.Time `bson:"ended_at,omitempty"`
	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at"`
}
