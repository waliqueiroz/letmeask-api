package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
)

func SetupAuthRoutes(router fiber.Router, authController *controllers.AuthController) {
	router.Post("/login", authController.Login)
}
