package services_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	repositoriesMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/database/mongodb/repositories/mocks"
	securityMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/security/mocks"
)

var _ = Describe("User", func() {

	Describe("Executing the FindAll function", func() {
		var result []entities.User
		var findAllError error
		var userService services.UserService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, findAllError = userService.FindAll()
		})

		When("the FindAll function is executed with success", func() {
			var expectedFindAllResult []entities.User

			BeforeEach(func() {
				expectedUsersSerialized, err := ioutil.ReadFile("../../../test/resources/full_users.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(expectedUsersSerialized, &expectedFindAllResult)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindAll().Return(expectedFindAllResult, nil).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be equal to expected FindAll result", func() {
				Expect(result).To(Equal(expectedFindAllResult))
			})

			It("error should be nil", func() {
				Expect(findAllError).Should(BeNil())
			})
		})

		When("an error occurs while executing the FindAll function", func() {
			BeforeEach(func() {
				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindAll().Return([]entities.User{}, errors.New("an error")).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be an empty array", func() {
				Expect(result).To(Equal([]entities.User{}))
			})

			It("error should be the error that comes from the repository", func() {
				Expect(findAllError).To(Equal(errors.New("an error")))
			})
		})
	})

})
