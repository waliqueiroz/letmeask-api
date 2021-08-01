package errors

import (
	"github.com/gofiber/fiber/v2"
	application "github.com/waliqueiroz/letmeask-api/application/errors"
)

func Handler(ctx *fiber.Ctx, err error) error {

	switch e := err.(type) {
	case *fiber.Error:
		return sendError(ctx, e.Code, e.Error())
	case *application.ResourceNotFoundError:
		return sendError(ctx, fiber.StatusNotFound, e.Error())
	case *application.UnauthorizedError:
		return sendError(ctx, fiber.StatusUnauthorized, e.Error())
	default:
		return sendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

}

func sendError(ctx *fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
}
