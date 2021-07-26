package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/application/services"
	"github.com/waliqueiroz/letmeask-api/infra/configuration"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
	"github.com/waliqueiroz/letmeask-api/infra/database"
	"github.com/waliqueiroz/letmeask-api/infra/repositories"
	"github.com/waliqueiroz/letmeask-api/infra/routes"
)

func main() {
	configuration := configuration.Load()

	db, err := database.ConnectMongoDB(configuration)
	if err != nil {
		log.Fatalln(err)
	}

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	app := fiber.New()

	api := app.Group("/api")

	routes.SetupUserRoutes(api, userController)

	app.Listen(":8080")
}
