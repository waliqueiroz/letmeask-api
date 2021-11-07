package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
)

const (
	FIND_ALL_USERS_ROUTE       = "/users"
	CREATE_USER_ROUTE          = "/users"
	FIND_USER_BY_ID_ROUTE      = "/users/:userID"
	UPDATE_USER_ROUTE          = "/users/:userID"
	DELETE_USER_ROUTE          = "/users/:userID"
	UPDATE_USER_PASSWORD_ROUTE = "/users/:userID/update-password"
)

func SetupUserRoutes(router fiber.Router, authMiddleware fiber.Handler, userController *controllers.UserController) {
	router.Post(CREATE_USER_ROUTE, userController.Create)
	router.Get(FIND_ALL_USERS_ROUTE, authMiddleware, userController.Index)
	router.Get(FIND_USER_BY_ID_ROUTE, authMiddleware, userController.FindByID)
	router.Put(UPDATE_USER_ROUTE, authMiddleware, userController.Update)
	router.Delete(DELETE_USER_ROUTE, authMiddleware, userController.Delete)
	router.Post(UPDATE_USER_PASSWORD_ROUTE, authMiddleware, userController.UpdatePassword)
}
