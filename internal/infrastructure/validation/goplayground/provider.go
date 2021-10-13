package goplayground

import (
	brazilian_portuguese "github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ptbr_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
)

type GoPlaygroundValidatorProvider struct{}

func NewGoPlaygroundValidatorProvider() *GoPlaygroundValidatorProvider {
	return &GoPlaygroundValidatorProvider{}
}

func (provider *GoPlaygroundValidatorProvider) ValidateStruct(value interface{}) []dtos.ValidationErrorDTO {
	var errors []dtos.ValidationErrorDTO

	ptbr := brazilian_portuguese.New()
	uni := ut.New(ptbr, ptbr)
	trans, _ := uni.GetTranslator("pt_BR")

	validate := validator.New()
	ptbr_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(value)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {

			element := dtos.ValidationErrorDTO{
				Field:   e.Field(),
				Message: e.Translate(trans),
			}

			errors = append(errors, element)
		}
	}

	return errors
}
