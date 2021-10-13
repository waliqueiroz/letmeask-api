package goplayground

import (
	"github.com/go-playground/validator/v10"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
)

type GoPlaygroundValidatorProvider struct{}

func NewGoPlaygroundValidatorProvider() *GoPlaygroundValidatorProvider {
	return &GoPlaygroundValidatorProvider{}
}

func (provider *GoPlaygroundValidatorProvider) ValidateStruct(value interface{}) []dtos.ValidationErrorDTO {
	var errors []dtos.ValidationErrorDTO
	validate := validator.New()
	err := validate.Struct(value)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			element := dtos.ValidationErrorDTO{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			}

			errors = append(errors, element)
		}
	}

	return errors
}
