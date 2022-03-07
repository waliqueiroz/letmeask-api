package services_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	application "github.com/waliqueiroz/letmeask-api/internal/application/errors"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	authMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/authentication/mocks"
	repositoriesMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/database/mongodb/repositories/mocks"
	securityMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/security/mocks"
)

var _ = Describe("Auth", func() {

	Describe("Executing the Login funcition", func() {
		var result dtos.AuthDTO
		var authError error
		var credentials dtos.CredentialsDTO
		var authService services.AuthService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, authError = authService.Login(credentials)
		})

		When("the function flow is executed with success", func() {
			var expectedCreateTokenResult string
			var expectedFindByEmailResult entities.User

			BeforeEach(func() {
				credentialsSerialized, err := ioutil.ReadFile("../../../test/resources/credentials.json")
				Expect(err).NotTo(HaveOccurred())

				expectedUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(credentialsSerialized, &credentials)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(expectedUserSerialized, &expectedFindByEmailResult)
				Expect(err).NotTo(HaveOccurred())

				expectedCreateTokenResult = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MzUxODIyMTYsInVzZXJJRCI6IjYxNjQxYjA5NmJlYjg1YWRiMGU1ZDI5NyJ9.AiB8xDBV5-kNwFwZuf_gLv239IVqfPKYMXiMff_OzDU"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByEmail(credentials.Email).Return(expectedFindByEmailResult, nil).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedFindByEmailResult.Password, credentials.Password).Return(nil).Times(1)

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().CreateToken(expectedFindByEmailResult.ID, gomock.Any()).Return(expectedCreateTokenResult, nil).Times(1)

				authService = services.NewAuthService(mockUserRepository, mockSecurityProvider, mockAuthenticator)
			})

			It("result user should be equal to expected user", func() {
				Expect(result.User).To(Equal(expectedFindByEmailResult))
			})

			It("result access token should be equal to expected token", func() {
				Expect(result.AccessToken).To(Equal(expectedCreateTokenResult))
			})

			It("result token type should be equal to Bearer", func() {
				Expect(result.TokenType).To(Equal("Bearer"))
			})

			It("result expiresIn should be greater than zero", func() {
				Expect(result.ExpiresIn > 0).Should(BeTrue())
			})

			It("error should be nil", func() {
				Expect(authError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding user by email", func() {
			BeforeEach(func() {
				credentialsSerialized, err := ioutil.ReadFile("../../../test/resources/credentials.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(credentialsSerialized, &credentials)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByEmail(credentials.Email).Return(entities.User{}, errors.New("an error")).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				authService = services.NewAuthService(mockUserRepository, mockSecurityProvider, mockAuthenticator)
			})

			It("result should be an empty struct", func() {
				Expect(result).To(Equal(dtos.AuthDTO{}))
			})

			It("error should be an unauthorized error", func() {
				Expect(authError).To(Equal(application.NewUnauthorizedError("credenciais inválidas")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while verifying password", func() {
			BeforeEach(func() {
				credentialsSerialized, err := ioutil.ReadFile("../../../test/resources/credentials.json")
				Expect(err).NotTo(HaveOccurred())

				expectedUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(credentialsSerialized, &credentials)
				Expect(err).NotTo(HaveOccurred())

				var expectedUser entities.User
				err = json.Unmarshal(expectedUserSerialized, &expectedUser)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByEmail(credentials.Email).Return(expectedUser, nil).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedUser.Password, credentials.Password).Return(errors.New("an error")).Times(1)

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)

				authService = services.NewAuthService(mockUserRepository, mockSecurityProvider, mockAuthenticator)
			})

			It("result should be an empty struct", func() {
				Expect(result).To(Equal(dtos.AuthDTO{}))
			})

			It("error should an unauthorized error", func() {
				Expect(authError).To(Equal(application.NewUnauthorizedError("credenciais inválidas")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while creating token", func() {
			BeforeEach(func() {
				credentialsSerialized, err := ioutil.ReadFile("../../../test/resources/credentials.json")
				Expect(err).NotTo(HaveOccurred())

				expectedUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(credentialsSerialized, &credentials)
				Expect(err).NotTo(HaveOccurred())

				var expectedUser entities.User
				err = json.Unmarshal(expectedUserSerialized, &expectedUser)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByEmail(credentials.Email).Return(expectedUser, nil).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedUser.Password, credentials.Password).Return(nil).Times(1)

				mockAuthenticator := authMocks.NewMockAuthenticator(mockCtrl)
				mockAuthenticator.EXPECT().CreateToken(expectedUser.ID, gomock.Any()).Return("", errors.New("an error")).Times(1)

				authService = services.NewAuthService(mockUserRepository, mockSecurityProvider, mockAuthenticator)
			})

			It("result should be an empty struct", func() {
				Expect(result).To(Equal(dtos.AuthDTO{}))
			})

			It("error should an unauthorized error", func() {
				Expect(authError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

})
