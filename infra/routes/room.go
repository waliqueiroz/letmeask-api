package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
)

func SetupRoomRoutes(router fiber.Router, authMiddleware fiber.Handler, roomController *controllers.RoomController) {
	router.Post("/rooms", roomController.Create)
	router.Get("/rooms/:roomID", roomController.FindByID)
	router.Post("/rooms/:roomID/questions", roomController.CreateQuestion)
}
