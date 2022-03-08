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

	Describe("Executing the Create function", func() {
		var user entities.User
		var result entities.User
		var createError error
		var userService services.UserService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, createError = userService.Create(user)
		})

		When("the Create function is executed with success", func() {
			var expectedCreateResult entities.User

			BeforeEach(func() {
				createUserRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				fullUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createUserRequestSerialized, &user)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(fullUserSerialized, &expectedCreateResult)
				Expect(err).NotTo(HaveOccurred())

				hashedPassword := "$2a$10$Chs8KofcRGJxJpjMl.ZS8.bJgD8iDBfyLav/oahSGVaTwBmIUUMMm"

				userWithHashedPassword := user
				userWithHashedPassword.Password = hashedPassword

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Hash(user.Password).Return(hashedPassword, nil).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().Create(userWithHashedPassword).Return(expectedCreateResult, nil).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be equal to expected Create result", func() {
				Expect(result).To(Equal(expectedCreateResult))
			})

			It("error should be nil", func() {
				Expect(createError).Should(BeNil())
			})
		})

		When("an error occurs while generating password hash", func() {
			BeforeEach(func() {
				createUserRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createUserRequestSerialized, &user)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Hash(user.Password).Return("", errors.New("an error")).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be an empty User struct", func() {
				Expect(result).To(Equal(entities.User{}))
			})

			It("error should be the error returned by the Hash function", func() {
				Expect(createError).To(Equal(errors.New("an error")))
			})
		})

		When("an error occurs while saving user in database", func() {
			BeforeEach(func() {
				createUserRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createUserRequestSerialized, &user)
				Expect(err).NotTo(HaveOccurred())

				hashedPassword := "$2a$10$Chs8KofcRGJxJpjMl.ZS8.bJgD8iDBfyLav/oahSGVaTwBmIUUMMm"

				userWithHashedPassword := user
				userWithHashedPassword.Password = hashedPassword

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Hash(user.Password).Return(hashedPassword, nil).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().Create(userWithHashedPassword).Return(entities.User{}, errors.New("an error")).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be an empty User struct", func() {
				Expect(result).To(Equal(entities.User{}))
			})

			It("error should be the error returned by the userRepository.Create function", func() {
				Expect(createError).To(Equal(errors.New("an error")))
			})
		})
	})

})
