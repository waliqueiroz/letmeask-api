package entities

import "time"

type Room struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Title     string     `json:"title" bson:"title"`
	Questions []Question `json:"quations" bson:"questions"`
	Author    Author     `json:"author" bson:"author"`
	EndedAt   time.Time  `json:"ended_at" bson:"ended_at"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
}
