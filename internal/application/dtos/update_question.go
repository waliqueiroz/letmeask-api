package dtos

type UpdateQuestionDTO struct {
	IsAnswered    *bool `json:"is_answered" validate:"required_without=IsHighlighted"`
	IsHighlighted *bool `json:"is_highlighted" validate:"required_without=IsAnswered"`
}
