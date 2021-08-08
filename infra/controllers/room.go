package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/application/services"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
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

func (controller *RoomController) FindByID(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomID")

	room, err := controller.roomService.FindByID(roomID)
	if err != nil {
		return err
	}

	return ctx.JSON(room)
}
