package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/providers"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type UserController struct {
	userService        services.UserService
	validationProvider providers.ValidationProvider
}

func NewUserController(userService services.UserService, validationProvider providers.ValidationProvider) *UserController {
	return &UserController{
		userService,
		validationProvider,
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
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors := controller.validationProvider.ValidateStruct(user)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	user, err = controller.userService.Create(user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (controller *UserController) FindByID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")

	user, err := controller.userService.FindByID(userID)
	if err != nil {
		return err
	}

	return ctx.JSON(user)
}

func (controller *UserController) Update(ctx *fiber.Ctx) error {
	var user dtos.UserDTO

	err := ctx.BodyParser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors := controller.validationProvider.ValidateStruct(user)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
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

	errors := controller.validationProvider.ValidateStruct(password)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	userID := ctx.Params("userID")

	if err := controller.userService.UpdatePassword(userID, password); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
