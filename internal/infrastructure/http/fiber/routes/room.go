package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
)

func SetupRoomRoutes(router fiber.Router, authMiddleware fiber.Handler, roomController *controllers.RoomController) {
	router.Post("/rooms", authMiddleware, roomController.Create)
	router.Get("/rooms/:roomID", roomController.FindByID)
	router.Delete("/rooms/:roomID", authMiddleware, roomController.EndRoom)
	router.Post("/rooms/:roomID/questions", authMiddleware, roomController.CreateQuestion)
	router.Post("/rooms/:roomID/questions/:questionID/likes", authMiddleware, roomController.LikeQuestion)
	router.Delete("/rooms/:roomID/questions/:questionID/likes/:likeID", authMiddleware, roomController.DeslikeQuestion)
	router.Patch("/rooms/:roomID/questions/:questionID", authMiddleware, roomController.UpdateQuestion)
	router.Delete("/rooms/:roomID/questions/:questionID", authMiddleware, roomController.DeleteQuestion)
}
