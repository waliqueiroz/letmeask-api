package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/providers"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type RoomController struct {
	roomService   services.RoomService
	authenticator providers.Authenticator
	validator     providers.Validator
}

func NewRoomController(roomService services.RoomService, authProvider providers.Authenticator, validationProvider providers.Validator) *RoomController {
	return &RoomController{
		roomService,
		authProvider,
		validationProvider,
	}
}

func (controller *RoomController) Create(ctx *fiber.Ctx) error {
	var room entities.Room

	err := ctx.BodyParser(&room)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors := controller.validator.ValidateStruct(room)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	room, err = controller.roomService.Create(room)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(room)
}

func (controller *RoomController) EndRoom(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")

	userID, err := controller.authenticator.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	room, err := controller.roomService.EndRoom(userID, roomID)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (controller *RoomController) FindByID(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")

	room, err := controller.roomService.FindByID(roomID)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (controller *RoomController) CreateQuestion(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")

	var question entities.Question

	err := ctx.BodyParser(&question)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	errors := controller.validator.ValidateStruct(question)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	room, err := controller.roomService.CreateQuestion(roomID, question)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(room)
}

func (controller *RoomController) UpdateQuestion(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")
	questionID := ctx.Params("questionID")

	userID, err := controller.authenticator.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var questionData dtos.UpdateQuestionDTO

	err = ctx.BodyParser(&questionData)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	errors := controller.validator.ValidateStruct(questionData)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	room, err := controller.roomService.UpdateQuestion(userID, roomID, questionID, questionData)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (controller *RoomController) LikeQuestion(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")
	questionID := ctx.Params("questionID")

	var like entities.Like

	err := ctx.BodyParser(&like)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	errors := controller.validator.ValidateStruct(like)
	if errors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	room, err := controller.roomService.LikeQuestion(roomID, questionID, like)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (controller *RoomController) DeslikeQuestion(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")
	questionID := ctx.Params("questionID")
	likeID := ctx.Params("likeID")

	room, err := controller.roomService.DeslikeQuestion(roomID, questionID, likeID)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}

func (controller *RoomController) DeleteQuestion(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")
	questionID := ctx.Params("questionID")

	userID, err := controller.authenticator.ExtractUserID(ctx.Locals("user"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	room, err := controller.roomService.DeleteQuestion(userID, roomID, questionID)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}
