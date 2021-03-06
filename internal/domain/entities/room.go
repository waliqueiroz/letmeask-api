package entities

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/errors"
)

type Room struct {
	ID        string     `json:"id"`
	Title     string     `json:"title" validate:"required"`
	Questions []Question `json:"questions,omitempty"`
	Author    Author     `json:"author" validate:"required"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (room *Room) AddQuestion(question Question) {
	question.CreatedAt = time.Now()
	room.Questions = append(room.Questions, question)
}

func (room *Room) LikeQuestion(questionID string, like Like) error {
	like.CreatedAt = time.Now()
	for key, question := range room.Questions {
		if question.ID == questionID {
			room.Questions[key].AddLike(like)
			return nil
		}
	}

	return errors.NewResourceNotFoundError("pergunta não encontrada.")
}

func (room *Room) DeslikeQuestion(questionID string, likeID string) error {
	for key, question := range room.Questions {
		if question.ID == questionID {
			room.Questions[key].RemoveLike(likeID)
			return nil
		}
	}

	return errors.NewResourceNotFoundError("pergunta não encontrada.")
}

func (room *Room) MarkQuestionAsAnswered(questionID string) error {
	for key, question := range room.Questions {
		if question.ID == questionID {
			room.Questions[key].IsAnswered = true
			return nil
		}
	}

	return errors.NewResourceNotFoundError("pergunta não encontrada.")
}

func (room *Room) UpdateQuestionHighlight(questionID string, highligh bool) error {
	for key, question := range room.Questions {
		if question.ID == questionID {
			room.Questions[key].IsHighlighted = highligh
			return nil
		}
	}

	return errors.NewResourceNotFoundError("pergunta não encontrada.")
}

func (room *Room) DeleteQuestion(questionID string) {
	for key, question := range room.Questions {
		if question.ID == questionID {
			room.Questions = append(room.Questions[:key], room.Questions[key+1:]...)
			break
		}
	}
}
