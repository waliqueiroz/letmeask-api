package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/application/dtos"
	"github.com/waliqueiroz/letmeask-api/application/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService,
	}
}

func (controller *AuthController) Login(ctx *fiber.Ctx) error {
	var credentials dtos.CredentialsDTO

	err := ctx.BodyParser(&credentials)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	response, err := controller.authService.Login(credentials)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
