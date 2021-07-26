package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/application/services"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService,
	}
}

func (controller *UserController) Create(ctx *fiber.Ctx) error {
	var user entities.User

	err := ctx.BodyParser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	user, err = controller.userService.Create(user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}
