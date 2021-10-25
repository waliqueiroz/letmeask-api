package dtos

type UserDTO struct {
	Name   string `json:"name" validate:"required"`
	Avatar string `json:"avatar" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
}
