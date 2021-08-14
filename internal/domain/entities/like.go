package entities

import "time"

type Like struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Author    Author    `json:"author" bson:"author"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
