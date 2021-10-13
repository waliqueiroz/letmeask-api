package providers

import "github.com/waliqueiroz/letmeask-api/internal/application/dtos"

type ValidationProvider interface {
	ValidateStruct(value interface{}) []dtos.ValidationErrorDTO
}
