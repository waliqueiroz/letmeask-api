package errors

import (
	"github.com/gofiber/fiber/v2"
	domain "github.com/waliqueiroz/letmeask-api/internal/domain/errors"
)

func Handler(ctx *fiber.Ctx, err error) error {

	switch e := err.(type) {
	case *fiber.Error:
		return sendError(ctx, e.Code, e.Error())
	case domain.HTTPError:
		return sendError(ctx, e.Code(), err.Error())
	default:
		return sendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

}

func sendError(ctx *fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"message": message,
	})
}
