package entities

import "time"

type Question struct {
	ID            string    `json:"id" bson:"_id"`
	Content       string    `json:"content" bson:"content"`
	IsHighlighted bool      `json:"is_highlighted" bson:"is_highlighted"`
	IsAnswered    bool      `json:"is_answered" bson:"is_answered"`
	Author        Author    `json:"author" bson:"author"`
	Likes         []Like    `json:"likes,omitempty" bson:"likes,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

func (question *Question) AddLike(like Like) {
	question.Likes = append(question.Likes, like)
}
