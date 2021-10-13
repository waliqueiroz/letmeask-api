package dtos

type ValidationErrorDTO struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
