package dtos

type UpdateQuestionDTO struct {
	IsAnswered    *bool `json:"is_answered"`
	IsHighlighted *bool `json:"is_highlighted"`
}
