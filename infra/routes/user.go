package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
)

func SetupUserRoutes(router fiber.Router, userController *controllers.UserController) {
	router.Get("/users", userController.Index)
	router.Get("/users/:userID", userController.FindByID)
	router.Post("/users", userController.Create)
	router.Put("/users/:userID", userController.Update)
	router.Delete("/users/:userID", userController.Delete)
}
