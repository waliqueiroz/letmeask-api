package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/authentication/jwt"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

func TestCreateRoom(t *testing.T) {
	if err := godotenv.Load("../../../../../.env"); err != nil {
		log.Fatalln(err)
	}

	configuration := configurations.Load()

	createRoomRequestSerialized, _ := ioutil.ReadFile("../../../../../test/resources/create_room_request.json")
	createRoomRequestIncompleteSerialized, _ := ioutil.ReadFile("../../../../../test/resources/create_room_request_incomplete.json")
	roomSerialized, _ := ioutil.ReadFile("../../../../../test/resources/room.json")

	var room entities.Room
	json.Unmarshal(roomSerialized, &room)

	tests := []struct {
		name                 string
		input                *bytes.Buffer
		expectedCreateResult entities.Room
		expectedCreateCalls  int
		expectedCreateError  error
		expectedStatusCode   int
	}{
		{
			name:                 "Create room",
			input:                bytes.NewBuffer(createRoomRequestSerialized),
			expectedCreateResult: room,
			expectedCreateCalls:  1,
			expectedCreateError:  nil,
			expectedStatusCode:   fiber.StatusCreated,
		},
		{
			name:                 "Try to create room with incomplete data",
			input:                bytes.NewBuffer(createRoomRequestIncompleteSerialized),
			expectedCreateResult: entities.Room{},
			expectedCreateCalls:  0,
			expectedCreateError:  nil,
			expectedStatusCode:   fiber.StatusUnprocessableEntity,
		},
		{
			name:                 "Try to create room with invalid data",
			input:                bytes.NewBuffer(nil),
			expectedCreateResult: entities.Room{},
			expectedCreateCalls:  0,
			expectedCreateError:  nil,
			expectedStatusCode:   fiber.StatusBadRequest,
		},
		{
			name:                 "Error creating room",
			input:                bytes.NewBuffer(createRoomRequestSerialized),
			expectedCreateResult: entities.Room{},
			expectedCreateCalls:  1,
			expectedCreateError:  assert.AnError,
			expectedStatusCode:   fiber.StatusInternalServerError,
		},
	}

	validationProvider := goplayground.NewGoPlaygroundValidatorProvider()
	authProvider := jwt.NewJwtProvider(configuration)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			roomServiceMock := mocks.NewRoomServiceMock()
			roomServiceMock.On("Create", mock.AnythingOfType("entities.Room")).Return(test.expectedCreateResult, test.expectedCreateError)

			roomController := controllers.NewRoomController(roomServiceMock, authProvider, validationProvider)

			app := fiber.New()

			routes.SetupRoomRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, roomController)

			req := httptest.NewRequest(fiber.MethodPost, routes.CREATE_ROOM_ROUTE, test.input)
			req.Header.Set("Content-Type", "application/json")

			response, _ := app.Test(req)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			roomServiceMock.AssertNumberOfCalls(t, "Create", test.expectedCreateCalls)

			if response.StatusCode == fiber.StatusCreated {
				body, _ := ioutil.ReadAll(response.Body)
				var room entities.Room
				json.Unmarshal(body, &room)

				assert.Equal(t, test.expectedCreateResult, room)
			}
		})
	}

}
