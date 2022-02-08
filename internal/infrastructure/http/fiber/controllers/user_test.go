package controllers_test

import (
	"encoding/json"
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
		var mockUserService *mocks.MockUserService

		JustBeforeEach(func() {
			var err error
			validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

			userController := controllers.NewUserController(mockUserService, validationProvider)

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

				mockUserService = mocks.NewMockUserService(mockCtrl)
				mockUserService.EXPECT().FindAll().Return(expectedFindAllResult, nil).Times(1)
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
	})
})
