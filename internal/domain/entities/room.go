package entities

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/domain/errors"
)

type Room struct {
	ID        string     `json:"id" bson:"_id"`
	Title     string     `json:"title" bson:"title"`
	Questions []Question `json:"questions,omitempty" bson:"questions,omitempty"`
	Author    Author     `json:"author" bson:"author"`
	EndedAt   *time.Time `json:"ended_at,omitempty" bson:"ended_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
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

	return errors.NewQuestionNotFoundError("pergunta não encontrada.")
}
