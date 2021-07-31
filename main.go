package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/waliqueiroz/letmeask-api/application/services"
	"github.com/waliqueiroz/letmeask-api/infra/configurations"
	"github.com/waliqueiroz/letmeask-api/infra/controllers"
	"github.com/waliqueiroz/letmeask-api/infra/database"
	"github.com/waliqueiroz/letmeask-api/infra/providers"
	"github.com/waliqueiroz/letmeask-api/infra/repositories"
	"github.com/waliqueiroz/letmeask-api/infra/routes"
	"github.com/waliqueiroz/letmeask-api/infra/utils"
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

	app := fiber.New()

	api := app.Group("/api")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(configuration.Auth.SecretKey),
		Filter:     utils.FilterUnauthenticatedRoutes,
	}))

	routes.SetupAuthRoutes(api, authController)
	routes.SetupUserRoutes(api, userController)

	app.Listen(":8080")
}
