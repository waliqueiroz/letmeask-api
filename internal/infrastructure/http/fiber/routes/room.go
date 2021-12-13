package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
)

const CREATE_ROOM_ROUTE = "/rooms"
const FIND_ROOM_BY_ID_ROUTE = "/rooms/:roomID"
const END_ROOM_ROUTE = "/rooms/:roomID"
const CREATE_QUESTION_ROUTE = "/rooms/:roomID/questions"
const LIKE_QUESTION_ROUTE = "/rooms/:roomID/questions/:questionID/likes"
const DESLIKE_QUESTION_ROUTE = "/rooms/:roomID/questions/:questionID/likes/:likeID"
const UPDATE_QUESTION_ROUTE = "/rooms/:roomID/questions/:questionID"
const DELETE_QUESTION_ROUTE = "/rooms/:roomID/questions/:questionID"

func SetupRoomRoutes(router fiber.Router, authMiddleware fiber.Handler, roomController *controllers.RoomController) {
	router.Post(CREATE_ROOM_ROUTE, authMiddleware, roomController.Create)
	router.Get(FIND_ROOM_BY_ID_ROUTE, roomController.FindByID)
	router.Delete(END_ROOM_ROUTE, authMiddleware, roomController.EndRoom)
	router.Post(CREATE_QUESTION_ROUTE, authMiddleware, roomController.CreateQuestion)
	router.Post(LIKE_QUESTION_ROUTE, authMiddleware, roomController.LikeQuestion)
	router.Delete(DESLIKE_QUESTION_ROUTE, authMiddleware, roomController.DeslikeQuestion)
	router.Patch(UPDATE_QUESTION_ROUTE, authMiddleware, roomController.UpdateQuestion)
	router.Delete(DELETE_QUESTION_ROUTE, authMiddleware, roomController.DeleteQuestion)
}
