package entities

type Author struct {
	ID     string `json:"id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Avatar string `json:"avatar" validate:"required"`
}
