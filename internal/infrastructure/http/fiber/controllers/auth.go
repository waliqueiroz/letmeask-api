package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/providers"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
)

type AuthController struct {
	authService services.AuthService
	validator   providers.Validator
}

func NewAuthController(authService services.AuthService, validationProvider providers.Validator) *AuthController {
	return &AuthController{
		authService,
		validationProvider,
	}
}

func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	var credentials dtos.CredentialsDTO

	err := ctx.BodyParser(&credentials)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors := controller.validator.ValidateStruct(credentials)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	response, err := controller.authService.Login(credentials)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
