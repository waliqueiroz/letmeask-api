package utils

import "github.com/gofiber/fiber/v2"

func FilterUnauthenticatedRoutes(ctx *fiber.Ctx) bool {
	path := ctx.Path()
	method := ctx.Method()

	postUnauthenticatedRoutes := []string{
		"/api/login",
		"/api/users",
	}

	switch method {
	case fiber.MethodPost:
		return stringInSlice(path, postUnauthenticatedRoutes)
	}

	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
