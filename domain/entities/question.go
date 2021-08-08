package entities

import "time"

type Question struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Content       string    `json:"content" bson:"content"`
	IsHighlighted bool      `json:"is_highlighted" bson:"is_highlighted"`
	IsAnswered    bool      `json:"is_answered" bson:"is_answered"`
	Author        Author    `json:"author" bson:"author"`
	Likes         []Like    `json:"likes" bson:"likes"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}
