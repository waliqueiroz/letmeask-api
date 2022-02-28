package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/authentication/jwt"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations/env"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	infrastructure "github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/errors"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

var _ = Describe("Room", func() {
	var configuration configurations.Configuration

	BeforeEach(func() {
		envProvider := env.NewEnvProvider()
		configuration = envProvider.LoadConfigurationFromFile("../../../../../.env.test")
	})

	Describe("Create room", func() {
		var input *bytes.Buffer
		var response *http.Response
		var mockCtrl *gomock.Controller
		var roomController *controllers.RoomController

		JustBeforeEach(func() {
			var err error

			app := fiber.New(fiber.Config{
				ErrorHandler: infrastructure.Handler,
			})

			routes.SetupRoomRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, roomController)

			req := httptest.NewRequest(fiber.MethodPost, routes.CREATE_ROOM_ROUTE, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("create room with success", func() {
			var expectedCreateResult entities.Room

			BeforeEach(func() {
				createRoomRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_room_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomSerialized, err := ioutil.ReadFile("../../../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomSerialized, &expectedCreateResult)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(createRoomRequestSerialized, &room)
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(createRoomRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().Create(room).Return(expectedCreateResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()
				authProvider := jwt.NewJwtProvider(configuration)

				roomController = controllers.NewRoomController(mockRoomService, authProvider, validationProvider)
			})

			It("response status code should be equal to 201 Created", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusCreated))
			})

			It("response body should be equal to roomService.Create result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedCreateResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("validation fails while creating room", func() {
			BeforeEach(func() {
				createRoomRequestIncompleteSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_room_request_incomplete.json")
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(createRoomRequestIncompleteSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()
				authProvider := jwt.NewJwtProvider(configuration)

				roomController = controllers.NewRoomController(mockRoomService, authProvider, validationProvider)
			})

			It("response status code should be equal to 422 Unprocessable Entity", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusUnprocessableEntity))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("request body comes with an invalid payload", func() {
			BeforeEach(func() {
				input = bytes.NewBuffer(nil)

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()
				authProvider := jwt.NewJwtProvider(configuration)

				roomController = controllers.NewRoomController(mockRoomService, authProvider, validationProvider)
			})

			It("response status code should be equal to 400 Bad Request", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusBadRequest))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("a general error occurs while creating user", func() {
			BeforeEach(func() {
				createRoomRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_room_request.json")
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(createRoomRequestSerialized, &room)
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(createRoomRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().Create(room).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()
				authProvider := jwt.NewJwtProvider(configuration)

				roomController = controllers.NewRoomController(mockRoomService, authProvider, validationProvider)
			})

			It("response status code should be equal to 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

})
