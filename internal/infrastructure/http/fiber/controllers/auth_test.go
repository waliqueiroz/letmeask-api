package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		var mockAuthService *mocks.MockAuthService
		var expectedLoginResult dtos.AuthDTO
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			var err error
			validationProvider := goplayground.NewGoPlaygroundValidatorProvider()

			authController := controllers.NewAuthController(mockAuthService, validationProvider)

			app := fiber.New()

			routes.SetupAuthRoutes(app, authController)

			req := httptest.NewRequest(fiber.MethodPost, routes.LOGIN_ROUTE, input)
			req.Header.Set("Content-Type", "application/json")

			response, err = app.Test(req)
			Expect(err).NotTo(HaveOccurred())
		})

		When("login is performed with success", func() {
			BeforeEach(func() {
				credentialsSerialized, err := ioutil.ReadFile("../../../../../test/resources/credentials.json")
				Expect(err).NotTo(HaveOccurred())

				authSerialized, err := ioutil.ReadFile("../../../../../test/resources/auth.json")
				Expect(err).NotTo(HaveOccurred())

				// Entrada
				input = bytes.NewBuffer(credentialsSerialized)

				// Mocks
				err = json.Unmarshal(authSerialized, &expectedLoginResult)
				Expect(err).NotTo(HaveOccurred())

				var credentialsDTO dtos.CredentialsDTO
				err = json.Unmarshal(credentialsSerialized, &credentialsDTO)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockAuthService = mocks.NewMockAuthService(mockCtrl)
				mockAuthService.EXPECT().Login(credentialsDTO).Return(expectedLoginResult, nil).Times(1)
			})

			It("response status code should be 200 OK", func() {
				Expect(response.StatusCode).To(Equal(fiber.StatusOK))
			})

			It("response body should be equal to authService.Login result", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				var auth dtos.AuthDTO
				err = json.Unmarshal(body, &auth)
				Expect(err).NotTo(HaveOccurred())

				Expect(auth).To(Equal(expectedLoginResult))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

})
