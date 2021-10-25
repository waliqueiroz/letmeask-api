package dtos

type PasswordDTO struct {
	Current string `json:"current"  validate:"required"`
	New     string `json:"new"  validate:"required"`
}
