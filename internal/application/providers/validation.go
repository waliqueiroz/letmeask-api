package providers

import "github.com/waliqueiroz/letmeask-api/internal/application/dtos"

type Validator interface {
	ValidateStruct(value interface{}) []dtos.ValidationErrorDTO
}
