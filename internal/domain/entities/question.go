package entities

import "time"

type Question struct {
	ID            string    `json:"id"`
	Content       string    `json:"content" validate:"required"`
	IsHighlighted bool      `json:"is_highlighted"`
	IsAnswered    bool      `json:"is_answered"`
	Author        Author    `json:"author" validate:"dive"`
	Likes         []Like    `json:"likes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (question *Question) AddLike(newLike Like) {
	for _, like := range question.Likes {
		if like.Author.ID == newLike.Author.ID {
			return
		}
	}

	question.Likes = append(question.Likes, newLike)
}

func (question *Question) RemoveLike(likeID string) {
	for key, like := range question.Likes {
		if like.ID == likeID {
			question.Likes = append(question.Likes[:key], question.Likes[key+1:]...)
			break
		}
	}
}
