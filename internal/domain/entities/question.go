package entities

import "time"

type Question struct {
	ID            string    `json:"id"`
	Content       string    `json:"content"`
	IsHighlighted bool      `json:"is_highlighted"`
	IsAnswered    bool      `json:"is_answered"`
	Author        Author    `json:"author"`
	Likes         []Like    `json:"likes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (question *Question) AddLike(like Like) {
	question.Likes = append(question.Likes, like)
}
