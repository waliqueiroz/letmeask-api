package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
)

func SetupRoomRoutes(router fiber.Router, authMiddleware fiber.Handler, roomController *controllers.RoomController) {
	router.Post("/rooms", roomController.Create)
	router.Get("/rooms/:roomID", roomController.FindByID)
	router.Delete("/rooms/:roomID", roomController.EndRoom)
	router.Post("/rooms/:roomID/questions", roomController.CreateQuestion)
	router.Post("/rooms/:roomID/questions/:questionID/likes", roomController.LikeQuestion)
	router.Delete("/rooms/:roomID/questions/:questionID/likes/:likeID", roomController.DeslikeQuestion)
	router.Patch("/rooms/:roomID/questions/:questionID", roomController.UpdateQuestion)
	router.Delete("/rooms/:roomID/questions/:questionID", roomController.DeleteQuestion)
}
