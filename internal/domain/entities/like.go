package entities

import "time"

type Like struct {
	ID        string    `json:"id"`
	Author    Author    `json:"author" validate:"dive"`
	CreatedAt time.Time `json:"created_at"`
}
