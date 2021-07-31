package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/waliqueiroz/letmeask-api/infra/configurations"
)

func NewAuthMiddleware(configuration configurations.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(configuration.Auth.SecretKey),
	})
}
