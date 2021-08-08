package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/waliqueiroz/letmeask-api/application/services"
	"github.com/waliqueiroz/letmeask-api/infra/configurations"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
	"github.com/waliqueiroz/letmeask-api/infra/database"
	"github.com/waliqueiroz/letmeask-api/infra/errors"
	"github.com/waliqueiroz/letmeask-api/infra/middlewares"
	"github.com/waliqueiroz/letmeask-api/infra/providers"
	"github.com/waliqueiroz/letmeask-api/infra/repositories"
	"github.com/waliqueiroz/letmeask-api/infra/routes"
)

func main() {
	configuration := configurations.Load()

	db, err := database.ConnectMongoDB(configuration)
	if err != nil {
		log.Fatalln(err)
	}

	authProvider := providers.NewAuthProvider(configuration)
	securityProvider := providers.NewSecurityProvider()

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, securityProvider)
	userController := controllers.NewUserController(userService)

	authService := services.NewAuthService(userRepository, securityProvider, authProvider)
	authController := controllers.NewAuthController(authService)

	roomRepository := repositories.NewRoomRepository(db)
	roomService := services.NewRoomService(roomRepository)
	roomController := controllers.NewRoomController(roomService)

	authMiddleware := middlewares.NewAuthMiddleware(configuration)

	app := fiber.New(fiber.Config{
		ErrorHandler: errors.Handler,
	})

	api := app.Group("/api")

	routes.SetupAuthRoutes(api, authController)
	routes.SetupUserRoutes(api, authMiddleware, userController)
	routes.SetupRoomRoutes(api, authMiddleware, roomController)

	app.Listen(":8080")
}
