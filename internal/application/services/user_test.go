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

			It("result should be an empty array of users", func() {
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

			It("result should be equal to expected userRepository.Create result", func() {
				Expect(result).To(Equal(expectedCreateResult))
			})

			It("error should be nil", func() {
				Expect(createError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
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

			AfterEach(func() {
				mockCtrl.Finish()
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

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the FindAll function", func() {
		var userID string
		var result entities.User
		var findByIDError error
		var userService services.UserService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, findByIDError = userService.FindByID(userID)
		})

		When("the FindByID function is executed with success", func() {
			var expectedFindByIDResult entities.User

			BeforeEach(func() {
				expectedUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(expectedUserSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(expectedFindByIDResult, nil).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be equal to expected userRepository.FindByID result", func() {
				Expect(result).To(Equal(expectedFindByIDResult))
			})

			It("error should be nil", func() {
				Expect(findByIDError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while executing the FindByID function", func() {
			BeforeEach(func() {
				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(entities.User{}, errors.New("an error")).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be an empty User struct", func() {
				Expect(result).To(Equal(entities.User{}))
			})

			It("error should be the error returned by the userRepository.FindByID function", func() {
				Expect(findByIDError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the Update function", func() {
		var userID string
		var userDTO dtos.UserDTO
		var result entities.User
		var updateError error
		var userService services.UserService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, updateError = userService.Update(userID, userDTO)
		})

		When("the Update function is executed with success", func() {
			var expectedUpdateResult entities.User

			BeforeEach(func() {
				updateUserRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				fullUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updateUserRequestSerialized, &userDTO)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(fullUserSerialized, &expectedUpdateResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"

				user := entities.User{
					Name:   userDTO.Name,
					Email:  userDTO.Email,
					Avatar: userDTO.Avatar,
				}

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().Update(userID, user).Return(expectedUpdateResult, nil).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be equal to expected userRepository.Update result", func() {
				Expect(result).To(Equal(expectedUpdateResult))
			})

			It("error should be nil", func() {
				Expect(updateError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while executing the Update function", func() {
			BeforeEach(func() {
				updateUserRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_user_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updateUserRequestSerialized, &userDTO)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"

				user := entities.User{
					Name:   userDTO.Name,
					Email:  userDTO.Email,
					Avatar: userDTO.Avatar,
				}

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().Update(userID, user).Return(entities.User{}, errors.New("an error")).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("result should be an empty User struct", func() {
				Expect(result).To(Equal(entities.User{}))
			})

			It("error should be the error returned by the userRepository.Update function", func() {
				Expect(updateError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the Delete function", func() {
		var userID string
		var deleteError error
		var userService services.UserService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			deleteError = userService.Delete(userID)
		})

		When("the Delete function is executed with success", func() {
			BeforeEach(func() {
				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().Delete(userID).Return(nil).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be nil", func() {
				Expect(deleteError).Should(BeNil())
			})
		})

		When("an error occurs while executing the Delete function", func() {
			BeforeEach(func() {
				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().Delete(userID).Return(errors.New("an error")).Times(1)

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be the error returned by the userRepository.Delete function", func() {
				Expect(deleteError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the UpdatePassword function", func() {
		var userID string
		var passwordDTO dtos.PasswordDTO
		var updatePasswordError error
		var userService services.UserService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			updatePasswordError = userService.UpdatePassword(userID, passwordDTO)
		})

		When("the UpdatePassword function is executed with success", func() {
			BeforeEach(func() {
				updatePasswordRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_password_request.json")
				Expect(err).NotTo(HaveOccurred())

				fullUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updatePasswordRequestSerialized, &passwordDTO)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.User
				err = json.Unmarshal(fullUserSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"
				hashedPassword := "$2a$10$Chs8KofcRGJxJpjMl.ZS8.bJgD8iDBfyLav/oahSGVaTwBmIUUMMm"

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedFindByIDResult.Password, passwordDTO.Current).Return(nil).Times(1)
				mockSecurityProvider.EXPECT().Hash(passwordDTO.New).Return(hashedPassword, nil).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(expectedFindByIDResult, nil).Times(1)
				mockUserRepository.EXPECT().UpdatePassword(userID, hashedPassword).Return(nil).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be nil", func() {
				Expect(updatePasswordError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding user by ID", func() {
			BeforeEach(func() {
				updatePasswordRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_password_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updatePasswordRequestSerialized, &passwordDTO)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(entities.User{}, errors.New("an error")).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be the error returned by the userRepository.FindByID function", func() {
				Expect(updatePasswordError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while verifying password", func() {
			BeforeEach(func() {
				updatePasswordRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_password_request.json")
				Expect(err).NotTo(HaveOccurred())

				fullUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updatePasswordRequestSerialized, &passwordDTO)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.User
				err = json.Unmarshal(fullUserSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedFindByIDResult.Password, passwordDTO.Current).Return(errors.New("an error")).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(expectedFindByIDResult, nil).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be an unauthorized error", func() {
				Expect(updatePasswordError).To(Equal(application.NewUnauthorizedError("a operação falhou, revise os dados e tente novamente")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while hashing password", func() {
			BeforeEach(func() {
				updatePasswordRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_password_request.json")
				Expect(err).NotTo(HaveOccurred())

				fullUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updatePasswordRequestSerialized, &passwordDTO)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.User
				err = json.Unmarshal(fullUserSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedFindByIDResult.Password, passwordDTO.Current).Return(nil).Times(1)
				mockSecurityProvider.EXPECT().Hash(passwordDTO.New).Return("", errors.New("an error")).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(expectedFindByIDResult, nil).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be the error returned by the securityProvider.Hash function", func() {
				Expect(updatePasswordError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating password in database", func() {
			BeforeEach(func() {
				updatePasswordRequestSerialized, err := ioutil.ReadFile("../../../test/resources/update_password_request.json")
				Expect(err).NotTo(HaveOccurred())

				fullUserSerialized, err := ioutil.ReadFile("../../../test/resources/full_user.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(updatePasswordRequestSerialized, &passwordDTO)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.User
				err = json.Unmarshal(fullUserSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "6117e377b6e7bae09f52c483"
				hashedPassword := "$2a$10$Chs8KofcRGJxJpjMl.ZS8.bJgD8iDBfyLav/oahSGVaTwBmIUUMMm"

				mockCtrl = gomock.NewController(GinkgoT())

				mockSecurityProvider := securityMocks.NewMockSecurityProvider(mockCtrl)
				mockSecurityProvider.EXPECT().Verify(expectedFindByIDResult.Password, passwordDTO.Current).Return(nil).Times(1)
				mockSecurityProvider.EXPECT().Hash(passwordDTO.New).Return(hashedPassword, nil).Times(1)

				mockUserRepository := repositoriesMocks.NewMockUserRepository(mockCtrl)
				mockUserRepository.EXPECT().FindByID(userID).Return(expectedFindByIDResult, nil).Times(1)
				mockUserRepository.EXPECT().UpdatePassword(userID, hashedPassword).Return(errors.New("an error")).Times(1)

				userService = services.NewUserService(mockUserRepository, mockSecurityProvider)
			})

			It("error should be the error returned by the userRepository.UpdatePassword function", func() {
				Expect(updatePasswordError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

})
