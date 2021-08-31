package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type RoomController struct {
	roomService services.RoomService
}

func NewRoomController(roomService services.RoomService) *RoomController {
	return &RoomController{
		roomService,
	}
}

func (controller *RoomController) Create(ctx *fiber.Ctx) error {
	var room entities.Room

	err := ctx.BodyParser(&room)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	room, err = controller.roomService.Create(room)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(room)
}

func (controller *RoomController) EndRoom(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")

	room, err := controller.roomService.EndRoom(roomID)
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

	room, err := controller.roomService.CreateQuestion(roomID, question)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(room)
}

func (controller *RoomController) UpdateQuestion(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")
	questionID := ctx.Params("questionID")

	var questionData dtos.UpdateQuestionDTO

	err := ctx.BodyParser(&questionData)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	room, err := controller.roomService.UpdateQuestion(roomID, questionID, questionData)
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

	room, err := controller.roomService.DeleteQuestion(roomID, questionID)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}
