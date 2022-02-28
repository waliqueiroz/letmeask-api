package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	domain "github.com/waliqueiroz/letmeask-api/internal/domain/errors"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/errors"
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

func TestCreateUser(t *testing.T) {
	createUserRequestSerialized, _ := ioutil.ReadFile("../../../../../test/resources/create_user_request.json")
	createUserRequestIncompleteSerialized, _ := ioutil.ReadFile("../../../../../test/resources/create_user_request_incomplete.json")
	createdUserSerialized, _ := ioutil.ReadFile("../../../../../test/resources/user.json")

	var createdUser entities.User
	json.Unmarshal(createdUserSerialized, &createdUser)

	tests := []struct {
		name                 string
		input                *bytes.Buffer
		expectedCreateResult entities.User
		expectedCreateCalls  int
		expectedCreateError  error
		expectedStatusCode   int
	}{
		{
			name:                 "Create user",
			input:                bytes.NewBuffer(createUserRequestSerialized),
			expectedCreateResult: createdUser,
			expectedCreateError:  nil,
			expectedCreateCalls:  1,
			expectedStatusCode:   fiber.StatusCreated,
		},
		{
			name:                 "Validation error while creating user",
			input:                bytes.NewBuffer(createUserRequestIncompleteSerialized),
			expectedCreateResult: entities.User{},
			expectedCreateError:  nil,
			expectedCreateCalls:  0,
			expectedStatusCode:   fiber.StatusUnprocessableEntity,
		},
		{
			name:                 "Unmarshal error while creating user",
			input:                bytes.NewBuffer(nil),
			expectedCreateResult: entities.User{},
			expectedCreateError:  nil,
			expectedCreateCalls:  0,
			expectedStatusCode:   fiber.StatusBadRequest,
		},
		{
			name:                 "Error while creating user",
			input:                bytes.NewBuffer(createUserRequestSerialized),
			expectedCreateResult: entities.User{},
			expectedCreateError:  assert.AnError,
			expectedCreateCalls:  1,
			expectedStatusCode:   fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userServiceMock := mocks.NewUserServiceMock()
			userServiceMock.On("Create", mock.AnythingOfType("entities.User")).Return(test.expectedCreateResult, test.expectedCreateError)

			userController := controllers.NewUserController(userServiceMock, validationProvider)

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			req := httptest.NewRequest(fiber.MethodPost, routes.CREATE_USER_ROUTE, test.input)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			userServiceMock.AssertNumberOfCalls(t, "Create", test.expectedCreateCalls)

			if response.StatusCode == fiber.StatusCreated {
				body, _ := ioutil.ReadAll(response.Body)
				var user entities.User
				json.Unmarshal(body, &user)

				assert.Equal(t, test.expectedCreateResult, user)
			}
		})
	}

}

func TestFindUserByID(t *testing.T) {
	userSerialized, _ := ioutil.ReadFile("../../../../../test/resources/user.json")

	var user entities.User
	json.Unmarshal(userSerialized, &user)

	tests := []struct {
		name                   string
		userID                 string
		expectedFindByIDResult entities.User
		expectedFindByIDCalls  int
		expectedFindByIDError  error
		expectedStatusCode     int
	}{
		{
			name:                   "Find user by ID",
			userID:                 "6117e377b6e7bae09f52c483",
			expectedFindByIDResult: user,
			expectedFindByIDError:  nil,
			expectedFindByIDCalls:  1,
			expectedStatusCode:     fiber.StatusOK,
		},
		{
			name:                   "User not found while finding by ID",
			userID:                 "6117e377b6e7bae09f5333383",
			expectedFindByIDResult: entities.User{},
			expectedFindByIDError:  domain.NewResourceNotFoundError(),
			expectedFindByIDCalls:  1,
			expectedStatusCode:     fiber.StatusNotFound,
		},
		{
			name:                   "Error finding user by ID",
			userID:                 "6117e377b6e7bae09f52c483",
			expectedFindByIDResult: entities.User{},
			expectedFindByIDError:  assert.AnError,
			expectedFindByIDCalls:  1,
			expectedStatusCode:     fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userServiceMock := mocks.NewUserServiceMock()
			userServiceMock.On("FindByID", test.userID).Return(test.expectedFindByIDResult, test.expectedFindByIDError)

			userController := controllers.NewUserController(userServiceMock, validationProvider)

			app := fiber.New(fiber.Config{
				ErrorHandler: errors.Handler,
			})

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			route := strings.Replace(routes.FIND_USER_BY_ID_ROUTE, ":userID", test.userID, 1)

			req := httptest.NewRequest(fiber.MethodGet, route, nil)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			userServiceMock.AssertNumberOfCalls(t, "FindByID", test.expectedFindByIDCalls)

			if response.StatusCode == fiber.StatusOK {
				body, _ := ioutil.ReadAll(response.Body)
				var user entities.User
				json.Unmarshal(body, &user)

				assert.Equal(t, test.expectedFindByIDResult, user)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	updateUserRequestSerialized, _ := ioutil.ReadFile("../../../../../test/resources/update_user_request.json")
	updateUserRequestIncompleteSerialized, _ := ioutil.ReadFile("../../../../../test/resources/update_user_request_incomplete.json")
	updatedUserSerialized, _ := ioutil.ReadFile("../../../../../test/resources/user.json")

	var updatedUser entities.User
	json.Unmarshal(updatedUserSerialized, &updatedUser)

	tests := []struct {
		name                 string
		userID               string
		input                *bytes.Buffer
		expectedUpdateResult entities.User
		expectedUpdateCalls  int
		expectedUpdateError  error
		expectedStatusCode   int
	}{
		{
			name:                 "Update user",
			userID:               "6117e377b6e7bae09f52c483",
			input:                bytes.NewBuffer(updateUserRequestSerialized),
			expectedUpdateResult: updatedUser,
			expectedUpdateError:  nil,
			expectedUpdateCalls:  1,
			expectedStatusCode:   fiber.StatusOK,
		},
		{
			name:                 "Validation error while updating user",
			userID:               "6117e377b6e7bae09f52c483",
			input:                bytes.NewBuffer(updateUserRequestIncompleteSerialized),
			expectedUpdateResult: entities.User{},
			expectedUpdateError:  nil,
			expectedUpdateCalls:  0,
			expectedStatusCode:   fiber.StatusUnprocessableEntity,
		},
		{
			name:                 "Unmarshal error while updating user",
			userID:               "6117e377b6e7bae09f52c483",
			input:                bytes.NewBuffer(nil),
			expectedUpdateResult: entities.User{},
			expectedUpdateError:  nil,
			expectedUpdateCalls:  0,
			expectedStatusCode:   fiber.StatusBadRequest,
		},
		{
			name:                 "Error while updating user",
			userID:               "6117e377b6e7bae09f52c483",
			input:                bytes.NewBuffer(updateUserRequestSerialized),
			expectedUpdateResult: entities.User{},
			expectedUpdateError:  assert.AnError,
			expectedUpdateCalls:  1,
			expectedStatusCode:   fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userServiceMock := mocks.NewUserServiceMock()
			userServiceMock.On("Update", test.userID, mock.AnythingOfType("dtos.UserDTO")).Return(test.expectedUpdateResult, test.expectedUpdateError)

			userController := controllers.NewUserController(userServiceMock, validationProvider)

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			route := strings.Replace(routes.UPDATE_USER_ROUTE, ":userID", test.userID, 1)

			req := httptest.NewRequest(fiber.MethodPut, route, test.input)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			userServiceMock.AssertNumberOfCalls(t, "Update", test.expectedUpdateCalls)

			if response.StatusCode == fiber.StatusOK {
				body, _ := ioutil.ReadAll(response.Body)
				var user entities.User
				json.Unmarshal(body, &user)

				assert.Equal(t, test.expectedUpdateResult, user)
			}
		})
	}

}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name                string
		userID              string
		expectedDeleteCalls int
		expectedDeleteError error
		expectedStatusCode  int
	}{
		{
			name:                "Delete user",
			userID:              "6117e377b6e7bae09f52c483",
			expectedDeleteError: nil,
			expectedDeleteCalls: 1,
			expectedStatusCode:  fiber.StatusOK,
		},
		{
			name:                "Error deleting user",
			userID:              "6117e377b6e7bae09f52c483",
			expectedDeleteError: assert.AnError,
			expectedDeleteCalls: 1,
			expectedStatusCode:  fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userServiceMock := mocks.NewUserServiceMock()
			userServiceMock.On("Delete", test.userID).Return(test.expectedDeleteError)

			userController := controllers.NewUserController(userServiceMock, validationProvider)

			app := fiber.New(fiber.Config{
				ErrorHandler: errors.Handler,
			})

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			route := strings.Replace(routes.DELETE_USER_ROUTE, ":userID", test.userID, 1)

			req := httptest.NewRequest(fiber.MethodDelete, route, nil)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			userServiceMock.AssertNumberOfCalls(t, "Delete", test.expectedDeleteCalls)
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	updatePasswordRequestSerialized, _ := ioutil.ReadFile("../../../../../test/resources/update_password_request.json")
	updatePasswordRequestIncompleteSerialized, _ := ioutil.ReadFile("../../../../../test/resources/update_password_request_incomplete.json")

	tests := []struct {
		name                        string
		userID                      string
		input                       *bytes.Buffer
		expectedUpdatePasswordCalls int
		expectedUpdatePasswordError error
		expectedStatusCode          int
	}{
		{
			name:                        "Update user password",
			userID:                      "6117e377b6e7bae09f52c483",
			input:                       bytes.NewBuffer(updatePasswordRequestSerialized),
			expectedUpdatePasswordError: nil,
			expectedUpdatePasswordCalls: 1,
			expectedStatusCode:          fiber.StatusOK,
		},
		{
			name:                        "Validation error while updating user password ",
			userID:                      "6117e377b6e7bae09f52c483",
			input:                       bytes.NewBuffer(updatePasswordRequestIncompleteSerialized),
			expectedUpdatePasswordError: nil,
			expectedUpdatePasswordCalls: 0,
			expectedStatusCode:          fiber.StatusUnprocessableEntity,
		},
		{
			name:                        "Unmarshal error while updating user password",
			userID:                      "6117e377b6e7bae09f52c483",
			input:                       bytes.NewBuffer(nil),
			expectedUpdatePasswordError: nil,
			expectedUpdatePasswordCalls: 0,
			expectedStatusCode:          fiber.StatusBadRequest,
		},
		{
			name:                        "Error while updating user password",
			userID:                      "6117e377b6e7bae09f52c483",
			input:                       bytes.NewBuffer(updatePasswordRequestSerialized),
			expectedUpdatePasswordError: assert.AnError,
			expectedUpdatePasswordCalls: 1,
			expectedStatusCode:          fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userServiceMock := mocks.NewUserServiceMock()
			userServiceMock.On("UpdatePassword", test.userID, mock.AnythingOfType("dtos.PasswordDTO")).Return(test.expectedUpdatePasswordError)

			userController := controllers.NewUserController(userServiceMock, validationProvider)

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			route := strings.Replace(routes.UPDATE_USER_PASSWORD_ROUTE, ":userID", test.userID, 1)

			req := httptest.NewRequest(fiber.MethodPost, route, test.input)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			userServiceMock.AssertNumberOfCalls(t, "UpdatePassword", test.expectedUpdatePasswordCalls)
		})
	}

}