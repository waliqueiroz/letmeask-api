package entities

import "time"

type Room struct {
	ID        string     `json:"id" bson:"_id"`
	Title     string     `json:"title" bson:"title"`
	Questions []Question `json:"questions,omitempty" bson:"questions,omitempty"`
	Author    Author     `json:"author" bson:"author"`
	EndedAt   *time.Time `json:"ended_at,omitempty" bson:"ended_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
}
