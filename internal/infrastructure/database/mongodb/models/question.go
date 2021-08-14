package models

import "time"

type Question struct {
	ID            string    `bson:"_id"`
	Content       string    `bson:"content"`
	IsHighlighted bool      `bson:"is_highlighted"`
	IsAnswered    bool      `bson:"is_answered"`
	Author        Author    `bson:"author"`
	Likes         []Like    `bson:"likes,omitempty"`
	CreatedAt     time.Time `bson:"created_at"`
}
