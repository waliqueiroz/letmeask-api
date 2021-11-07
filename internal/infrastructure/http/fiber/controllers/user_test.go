package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

func TestUserIndex(t *testing.T) {
	usersSerialized, _ := ioutil.ReadFile("../../../../../test/resources/users.json")

	var users []entities.User
	json.Unmarshal(usersSerialized, &users)

	tests := []struct {
		name                  string
		expectedFindAllResult []entities.User
		expectedFindAllCalls  int
		expectedFindAllError  error
		expectedStatusCode    int
	}{
		{
			name:                  "Find all users",
			expectedFindAllResult: users,
			expectedFindAllError:  nil,
			expectedFindAllCalls:  1,
			expectedStatusCode:    fiber.StatusOK,
		},
		{
			name:                  "Find all users error",
			expectedFindAllResult: nil,
			expectedFindAllError:  assert.AnError,
			expectedFindAllCalls:  1,
			expectedStatusCode:    fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userServiceMock := mocks.NewUserServiceMock()
			userServiceMock.On("FindAll").Return(test.expectedFindAllResult, test.expectedFindAllError)

			userController := controllers.NewUserController(userServiceMock, validationProvider)

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			req := httptest.NewRequest(fiber.MethodGet, routes.FIND_ALL_USERS_ROUTE, nil)

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			userServiceMock.AssertNumberOfCalls(t, "FindAll", test.expectedFindAllCalls)

			if response.StatusCode == fiber.StatusOK {
				body, _ := ioutil.ReadAll(response.Body)
				var users []entities.User
				json.Unmarshal(body, &users)

				assert.Equal(t, test.expectedFindAllResult, users)
			}
		})
	}
}
