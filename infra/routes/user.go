package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
)

func SetupUserRoutes(router fiber.Router, authMiddleware fiber.Handler, userController *controllers.UserController) {
	router.Post("/users", userController.Create)
	router.Get("/users", authMiddleware, userController.Index)
	router.Get("/users/:userID", authMiddleware, userController.FindByID)
	router.Put("/users/:userID", authMiddleware, userController.Update)
	router.Delete("/users/:userID", authMiddleware, userController.Delete)
	router.Post("/users/:userID/update-password", authMiddleware, userController.UpdatePassword)
}
