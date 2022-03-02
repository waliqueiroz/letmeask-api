package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	authMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/authentication/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	infrastructure "github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/errors"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

var _ = Describe("Room", func() {
	Describe("Creating a room", func() {
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

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().Create(room).Return(expectedCreateResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().Create(room).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Ending a room", func() {
		var roomID string
		var response *http.Response
		var mockCtrl *gomock.Controller
		var roomController *controllers.RoomController

		JustBeforeEach(func() {
			var err error

			app := fiber.New(fiber.Config{
				ErrorHandler: infrastructure.Handler,
			})

			routes.SetupRoomRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, roomController)

			route := strings.Replace(routes.END_ROOM_ROUTE, ":roomID", roomID, 1)

			req := httptest.NewRequest(fiber.MethodDelete, route, nil)

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("end room with success", func() {
			var expectedEndRoomResult entities.Room

			BeforeEach(func() {
				endedRoomSerialized, err := ioutil.ReadFile("../../../../../test/resources/ended_room.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(endedRoomSerialized, &expectedEndRoomResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				userID := "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().EndRoom(userID, roomID).Return(expectedEndRoomResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to roomService.EndRoom result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedEndRoomResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while extracting user ID", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return("", errors.New("an error")).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusUnauthorized))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an general error occurs while ending room", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"
				userID := "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().EndRoom(userID, roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Finding a room by ID", func() {
		var roomID string
		var response *http.Response
		var mockCtrl *gomock.Controller
		var roomController *controllers.RoomController

		JustBeforeEach(func() {
			var err error

			app := fiber.New(fiber.Config{
				ErrorHandler: infrastructure.Handler,
			})

			routes.SetupRoomRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, roomController)

			route := strings.Replace(routes.FIND_ROOM_BY_ID_ROUTE, ":roomID", roomID, 1)

			req := httptest.NewRequest(fiber.MethodGet, route, nil)

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("find a room by ID with success", func() {
			var expectedFindByIDResult entities.Room

			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to roomService.FindByID result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedFindByIDResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("a general error occurs while finding a room by ID", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Creating a question", func() {
		var roomID string
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

			route := strings.Replace(routes.CREATE_QUESTION_ROUTE, ":roomID", roomID, 1)

			req := httptest.NewRequest(fiber.MethodPost, route, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("create a question with success", func() {
			var expectedCreateQuestionResult entities.Room

			BeforeEach(func() {
				createQuestionRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedCreateQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				var question entities.Question
				err = json.Unmarshal(createQuestionRequestSerialized, &question)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				input = bytes.NewBuffer(createQuestionRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().CreateQuestion(roomID, question).Return(expectedCreateQuestionResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 201 Created", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusCreated))
			})

			It("response body should be equal to roomService.CreateQuestion result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedCreateQuestionResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("validation fails while creating a question", func() {
			BeforeEach(func() {
				createQuestionIncompleteRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_question_request_incomplete.json")
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				input = bytes.NewBuffer(createQuestionIncompleteRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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
				roomID = "621f5ec1e07fdbb81c8221f7"

				input = bytes.NewBuffer(nil)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 400 Bad Request", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusBadRequest))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("a general error occurs while creating question", func() {
			BeforeEach(func() {
				createQuestionRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				var question entities.Question
				err = json.Unmarshal(createQuestionRequestSerialized, &question)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				input = bytes.NewBuffer(createQuestionRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().CreateQuestion(roomID, question).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Updating a question", func() {
		var roomID string
		var questionID string
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

			route := strings.Replace(routes.UPDATE_QUESTION_ROUTE, ":roomID", roomID, 1)
			route = strings.Replace(route, ":questionID", questionID, 1)

			req := httptest.NewRequest(fiber.MethodPatch, route, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("highlight question with success", func() {
			var expectedUpdateQuestionResult entities.Room

			BeforeEach(func() {
				highlightQuestionRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/highlight_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionHighlightedSerialized, err := ioutil.ReadFile("../../../../../test/resources/room_with_question_highlighted.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionHighlightedSerialized, &expectedUpdateQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				var questionData dtos.UpdateQuestionDTO
				err = json.Unmarshal(highlightQuestionRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID := "621f5e02e07fdbb81c8221f5"

				input = bytes.NewBuffer(highlightQuestionRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().UpdateQuestion(userID, roomID, questionID, questionData).Return(expectedUpdateQuestionResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to roomService.UpdateQuestion result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedUpdateQuestionResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("mark question as answered with success", func() {
			var expectedUpdateQuestionResult entities.Room

			BeforeEach(func() {
				markQuestionAsAnsweredRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/mark_question_as_answered_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionAnsweredSerialized, err := ioutil.ReadFile("../../../../../test/resources/room_with_question_answered.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionAnsweredSerialized, &expectedUpdateQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				var questionData dtos.UpdateQuestionDTO
				err = json.Unmarshal(markQuestionAsAnsweredRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID := "621f5e02e07fdbb81c8221f5"

				input = bytes.NewBuffer(markQuestionAsAnsweredRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().UpdateQuestion(userID, roomID, questionID, questionData).Return(expectedUpdateQuestionResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to roomService.UpdateQuestion result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedUpdateQuestionResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("validation fails while updating question", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID := "621f5e02e07fdbb81c8221f5"

				input = bytes.NewBuffer([]byte(`{}`))

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID := "621f5e02e07fdbb81c8221f5"

				input = bytes.NewBuffer(nil)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 400 Bad Request", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusBadRequest))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while extracting user ID", func() {
			BeforeEach(func() {
				markQuestionAsAnsweredRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/mark_question_as_answered_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				input = bytes.NewBuffer(markQuestionAsAnsweredRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return("", errors.New("an error")).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusUnauthorized))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("a general error occurs while updating question", func() {
			BeforeEach(func() {
				highlightQuestionRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/highlight_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				var questionData dtos.UpdateQuestionDTO
				err = json.Unmarshal(highlightQuestionRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID := "621f5e02e07fdbb81c8221f5"

				input = bytes.NewBuffer(highlightQuestionRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().ExtractUserID(gomock.Any()).Return(userID, nil).Times(1)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().UpdateQuestion(userID, roomID, questionID, questionData).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Liking a question", func() {
		var roomID string
		var questionID string
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

			route := strings.Replace(routes.LIKE_QUESTION_ROUTE, ":roomID", roomID, 1)
			route = strings.Replace(route, ":questionID", questionID, 1)

			req := httptest.NewRequest(fiber.MethodPost, route, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("like a question with success", func() {
			var expectedLikeQuestionResult entities.Room

			BeforeEach(func() {
				likeQuestionRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/like_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionLikedSerialized, err := ioutil.ReadFile("../../../../../test/resources/room_with_question_liked.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionLikedSerialized, &expectedLikeQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				var like entities.Like
				err = json.Unmarshal(likeQuestionRequestSerialized, &like)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				input = bytes.NewBuffer(likeQuestionRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().LikeQuestion(roomID, questionID, like).Return(expectedLikeQuestionResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to roomService.UpdateQuestion result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var room entities.Room
				err = json.Unmarshal(body, &room)
				Expect(err).NotTo(HaveOccurred())

				Expect(room).To(Equal(expectedLikeQuestionResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("validation fails while liking a question", func() {
			BeforeEach(func() {
				likeQuestionRequestIncompleteSerialized, err := ioutil.ReadFile("../../../../../test/resources/like_question_request_incomplete.json")
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				input = bytes.NewBuffer(likeQuestionRequestIncompleteSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				input = bytes.NewBuffer(nil)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
			})

			It("response status code should be equal to 400 Bad Request", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusBadRequest))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("a general error occurs while liking a question", func() {
			BeforeEach(func() {
				likeQuestionRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/like_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				var like entities.Like
				err = json.Unmarshal(likeQuestionRequestSerialized, &like)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				input = bytes.NewBuffer(likeQuestionRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				mockRoomService := mocks.NewMockRoomService(mockCtrl)
				mockRoomService.EXPECT().LikeQuestion(roomID, questionID, like).Return(entities.Room{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				roomController = controllers.NewRoomController(mockRoomService, mockAuthenticator, validationProvider)
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
