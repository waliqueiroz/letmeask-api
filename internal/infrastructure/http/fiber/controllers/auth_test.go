package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

func TestLogin(t *testing.T) {
	credentialsSerialized, _ := ioutil.ReadFile("../../../../../test/resources/credentials.json")
	incompleteCredentialsSerialized, _ := ioutil.ReadFile("../../../../../test/resources/incomplete_credentials.json")
	authSerialized, _ := ioutil.ReadFile("../../../../../test/resources/auth.json")

	var auth dtos.AuthDTO
	json.Unmarshal(authSerialized, &auth)

	tests := []struct {
		name                string
		input               *bytes.Buffer
		expectedLoginResult dtos.AuthDTO
		expectedLoginCalls  int
		expectedLoginError  error
		expectedStatusCode  int
	}{
		{
			name:                "Login with success",
			input:               bytes.NewBuffer(credentialsSerialized),
			expectedLoginResult: auth,
			expectedLoginError:  nil,
			expectedLoginCalls:  1,
			expectedStatusCode:  fiber.StatusOK,
		},
		{
			name:                "Error on login",
			input:               bytes.NewBuffer(credentialsSerialized),
			expectedLoginResult: dtos.AuthDTO{},
			expectedLoginError:  assert.AnError,
			expectedLoginCalls:  1,
			expectedStatusCode:  fiber.StatusInternalServerError,
		},
		{
			name:                "Login with incomplete credentials",
			input:               bytes.NewBuffer(incompleteCredentialsSerialized),
			expectedLoginResult: dtos.AuthDTO{},
			expectedLoginError:  nil,
			expectedLoginCalls:  0,
			expectedStatusCode:  fiber.StatusUnprocessableEntity,
		},
		{
			name:                "Error unmarshaling credentials",
			input:               bytes.NewBuffer(nil),
			expectedLoginResult: dtos.AuthDTO{},
			expectedLoginError:  nil,
			expectedLoginCalls:  0,
			expectedStatusCode:  fiber.StatusBadRequest,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authServicemock := mocks.NewAuthServiceMock()
			authServicemock.On("Login", mock.AnythingOfType("dtos.CredentialsDTO")).Return(test.expectedLoginResult, test.expectedLoginError)

			authController := controllers.NewAuthController(authServicemock, validationProvider)

			app := fiber.New()

			routes.SetupAuthRoutes(app, authController)

			req := httptest.NewRequest(fiber.MethodPost, routes.LOGIN_ROUTE, test.input)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			authServicemock.AssertNumberOfCalls(t, "Login", test.expectedLoginCalls)

			if response.StatusCode == fiber.StatusOK {
				body, _ := ioutil.ReadAll(response.Body)
				var auth dtos.AuthDTO
				json.Unmarshal(body, &auth)

				assert.Equal(t, test.expectedLoginResult, auth)
			}
		})
	}
}
