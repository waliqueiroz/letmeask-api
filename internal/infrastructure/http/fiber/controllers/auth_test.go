package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/services/mocks"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/controllers"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/http/fiber/routes"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/validation/goplayground"
)

var _ = Describe("Auth", func() {

	Describe("Performing login", func() {
		var input *bytes.Buffer
		var response *http.Response
		var authServicemock *mocks.AuthServiceMock
		var expectedLoginResult dtos.AuthDTO

		JustBeforeEach(func() {
			var err error
			validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

			authController := controllers.NewAuthController(authServicemock, validationProvider)

			app := fiber.New()

			routes.SetupAuthRoutes(app, authController)

			req := httptest.NewRequest(fiber.MethodPost, routes.LOGIN_ROUTE, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("login is performed with success", func() {
			BeforeEach(func() {
				// Entrada
				credentialsSerialized, err := ioutil.ReadFile("../../../../../test/resources/credentials.json")
				Expect(err).NotTo(HaveOccurred())

				input = bytes.NewBuffer(credentialsSerialized)

				// Mocks
				authSerialized, err := ioutil.ReadFile("../../../../../test/resources/auth.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(authSerialized, &expectedLoginResult)
				Expect(err).NotTo(HaveOccurred())

				authServicemock = mocks.NewAuthServiceMock()
				authServicemock.On("Login", mock.AnythingOfType("dtos.CredentialsDTO")).Return(expectedLoginResult, nil)
			})

			It("response status code should be 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("authServicemock.Login should be called once", func() {
				Expect(authServicemock.AssertNumberOfCalls(GinkgoT(), "Login", 1)).To(BeTrue())
			})

			It("response body should be equal to authServicemock.Login result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var auth dtos.AuthDTO
				err = json.Unmarshal(body, &auth)
				Expect(err).NotTo(HaveOccurred())

				Expect(auth).To(Equal(expectedLoginResult))
			})
		})
	})

})
