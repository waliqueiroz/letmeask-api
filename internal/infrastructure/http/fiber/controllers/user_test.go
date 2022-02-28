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
	domain "github.com/waliqueiroz/letmeask-api/internal/domain/errors"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	infrastructure "github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/errors"
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

		When("a general error occurs while creating user", func() {
			BeforeEach(func() {
				createUserRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/create_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(createUserRequestSerialized)

				var user entities.User
				err = json.Unmarshal(createUserRequestSerialized, &user)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().Create(user).Return(entities.User{}, errors.New("an error")).Times(1)

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

	Describe("Finding users by ID", func() {
		var userID string
		var response *http.Response
		var mockCtrl *gomock.Controller
		var userController *controllers.UserController

		JustBeforeEach(func() {
			var err error

			app := fiber.New(fiber.Config{
				ErrorHandler: infrastructure.Handler,
			})

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)

			route := strings.Replace(routes.FIND_USER_BY_ID_ROUTE, ":userID", userID, 1)

			req := httptest.NewRequest(fiber.MethodGet, route, nil)

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("find user by ID with success", func() {
			var expectedFindByIDResult entities.User

			BeforeEach(func() {
				userSerialized, err := ioutil.ReadFile("../../../../../test/resources/user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(userSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = expectedFindByIDResult.ID

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

				mockUserService.EXPECT().FindByID(userID).Return(expectedFindByIDResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to userService.FindByID result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var user entities.User
				err = json.Unmarshal(body, &user)
				Expect(err).NotTo(HaveOccurred())

				Expect(user).To(Equal(expectedFindByIDResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("user not found while finding by ID", func() {
			BeforeEach(func() {
				userID = "6117e377b6e7bae09f5399983"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

				mockUserService.EXPECT().FindByID(userID).Return(entities.User{}, domain.NewResourceNotFoundError()).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 404 Not Found", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusNotFound))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding user by ID", func() {
			BeforeEach(func() {
				userID = "6117e377b6e7bae09f5399983"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

				mockUserService.EXPECT().FindByID(userID).Return(entities.User{}, errors.New("an error")).Times(1)

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

	Describe("Updating users", func() {
		var userID string
		var input *bytes.Buffer
		var response *http.Response
		var mockCtrl *gomock.Controller
		var userController *controllers.UserController

		JustBeforeEach(func() {
			var err error

			app := fiber.New()

			routes.SetupUserRoutes(app, func(c *fiber.Ctx) error { return c.Next() }, userController)
			route := strings.Replace(routes.UPDATE_USER_ROUTE, ":userID", userID, 1)

			req := httptest.NewRequest(fiber.MethodPut, route, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("update user with success", func() {
			var expectedUpdateResult entities.User

			BeforeEach(func() {
				updateUserRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/update_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				updatedUserSerialized, err := ioutil.ReadFile("../../../../../test/resources/user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updatedUserSerialized, &expectedUpdateResult)
				Expect(err).NotTo(HaveOccurred())

				var updateUserRequest dtos.UserDTO
				err = json.Unmarshal(updateUserRequestSerialized, &updateUserRequest)
				Expect(err).NotTo(HaveOccurred())

				userID = expectedUpdateResult.ID
				input = bytes.NewBuffer(updateUserRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

				mockUserService.EXPECT().Update(userID, updateUserRequest).Return(expectedUpdateResult, nil).Times(1)

				validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

				userController = controllers.NewUserController(mockUserService, validationProvider)
			})

			It("response status code should be 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to userService.Update result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var user entities.User
				err = json.Unmarshal(body, &user)
				Expect(err).NotTo(HaveOccurred())

				Expect(user).To(Equal(expectedUpdateResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("validation fails while updating user", func() {
			BeforeEach(func() {
				updateUserRequestIncompleteSerialized, err := ioutil.ReadFile("../../../../../test/resources/update_user_request_incomplete.json")
				Expect(err).NotTo(HaveOccurred())

				var updateUserRequest dtos.UserDTO
				err = json.Unmarshal(updateUserRequestIncompleteSerialized, &updateUserRequest)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"
				input = bytes.NewBuffer(updateUserRequestIncompleteSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

				mockUserService.EXPECT().Update(userID, updateUserRequest).Return(entities.User{}, nil).Times(0)

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
				userID = "6117e377b6e7bae09f52c483"
				input = bytes.NewBuffer(nil)

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

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

		When("an error occurs while updating user", func() {
			BeforeEach(func() {
				updateUserRequestSerialized, err := ioutil.ReadFile("../../../../../test/resources/update_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				var updateUserRequest dtos.UserDTO
				err = json.Unmarshal(updateUserRequestSerialized, &updateUserRequest)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"
				input = bytes.NewBuffer(updateUserRequestSerialized)

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserService := mocks.NewMockUserService(mockCtrl)

				mockUserService.EXPECT().Update(userID, updateUserRequest).Return(entities.User{}, errors.New("an error")).Times(1)

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
})
