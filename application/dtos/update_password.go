package dtos

type PasswordDTO struct {
	Current string `json:"current"`
	New     string `json:"new"`
}
