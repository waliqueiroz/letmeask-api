package dtos

type ValidationErrorDTO struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
