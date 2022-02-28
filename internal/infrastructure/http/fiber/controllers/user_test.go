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
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

var _ = Describe("User", func() {

	Describe("Finding all users", func() {
		var response *http.Response
		var mockCtrl *gomock.Controller
		var userController *controllers.UserController

		JustBeforeEach(func() {
			var err error

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			req := httptest.NewRequest(fiber.MethodGet, routes.FIND_ALL_USERS_ROUTE, nil)

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("find all users with success", func() {
			var expectedFindAllResult []entities.User

			BeforeEach(func() {
				usersSerialized, err := ioutil.ReadFile("../../../../../test/resources/users.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(usersSerialized, &expectedFindAllResult)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().FindAll().Return(expectedFindAllResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to userService.FindAll result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var users []entities.User
				err = json.Unmarshal(body, &users)
				Expect(err).NotTo(HaveOccurred())

				Expect(users).To(Equal(expectedFindAllResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding all users", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().FindAll().Return([]entities.User{}, errors.New("an error")).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusInternalServerError))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Creating users", func() {
		var input *bytes.Buffer
		var response *http.Response
		var mockCtrl *gomock.Controller
		var userController *controllers.UserController

		JustBeforeEach(func() {
			var err error

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			req := httptest.NewRequest(fiber.MethodPost, routes.CREATE_USER_ROUTE, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("create user with success", func() {
			var expectedCreateResult entities.User

			BeforeEach(func() {
				createUserRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				createdUserSerialized, err := ioutil.ReadFile("../../../../../test/resources/user.json")
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(createUserRequestSerialized)

				err = json.Unmarshal(createdUserSerialized, &expectedCreateResult)
				Expect(err).NotTo(HaveOccurred())

				var user entities.User
				err = json.Unmarshal(createUserRequestSerialized, &user)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().Create(user).Return(expectedCreateResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 201 Created", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusCreated))
			})

			It("response body should be equal to userService.Create result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var user entities.User
				err = json.Unmarshal(body, &user)
				Expect(err).NotTo(HaveOccurred())

				Expect(user).To(Equal(expectedCreateResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("validation fails while creating user", func() {
			BeforeEach(func() {
				createUserRequestIncompleteSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_user_request_incomplete.json")
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(createUserRequestIncompleteSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().Create(gomock.Any()).Return(entities.User{}, nil).Times(0)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 422 Unprocessable Entity", func() {
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

				mockUserService := mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().Create(gomock.Any()).Return(entities.User{}, nil).Times(0)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 400 Bad Request", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusBadRequest))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})
})
