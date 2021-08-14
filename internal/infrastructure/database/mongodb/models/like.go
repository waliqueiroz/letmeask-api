package models

import "time"

type Like struct {
	ID        string    `bson:"_id,omitempty"`
	Author    Author    `bson:"author"`
	CreatedAt time.Time `bson:"created_at"`
}
