package dtos

import "github.com/waliqueiroz/letmeask-api/domain/entities"

type AuthResponseDTO struct {
	User        entities.User `json:"user"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int64         `json:"expires_in"`
}
