package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/application/dtos"
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

func (controller *UserController) Index(ctx *fiber.Ctx) error {
	users, err := controller.userService.FindAll()
	if err != nil {
		return err
	}

	return ctx.JSON(users)
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

func (controller *UserController) FindByID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	users, err := controller.userService.FindByID(userID)
	if err != nil {
		return err
	}

	return ctx.JSON(users)
}

func (controller *UserController) Update(ctx *fiber.Ctx) error {
	var user entities.User

	err := ctx.BodyParser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	userID := ctx.Params("userID")

	updatedUser, err := controller.userService.Update(userID, user)
	if err != nil {
		return err
	}

	return ctx.JSON(updatedUser)
}

func (controller *UserController) Delete(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	err := controller.userService.Delete(userID)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (controller *UserController) UpdatePassword(ctx *fiber.Ctx) error {
	var password dtos.PasswordDTO

	err := ctx.BodyParser(&password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	userID := ctx.Params("userID")

	if err := controller.userService.UpdatePassword(userID, password); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
