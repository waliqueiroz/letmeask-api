package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
)

func SetupUserRoutes(router fiber.Router, userController *controllers.UserController) {
	router.Post("/users", userController.Create)
}
