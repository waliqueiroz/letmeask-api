package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (service *UserServiceMock) FindAll() ([]entities.User, error) {
	args := service.Called()

	return args.Get(0).([]entities.User), args.Error(1)
}

func (service *UserServiceMock) Create(user entities.User) (entities.User, error) {
	args := service.Called(user)

	return args.Get(0).(entities.User), args.Error(1)
}

func (service *UserServiceMock) FindByID(userID string) (entities.User, error) {
	args := service.Called(userID)

	return args.Get(0).(entities.User), args.Error(1)
}

func (service *UserServiceMock) Update(userID string, userDTO dtos.UserDTO) (entities.User, error) {
	args := service.Called(userID, userDTO)

	return args.Get(0).(entities.User), args.Error(1)
}

func (service *UserServiceMock) Delete(userID string) error {
	args := service.Called(userID)

	return args.Error(0)
}

func (service *UserServiceMock) UpdatePassword(userID string, password dtos.PasswordDTO) error {
	args := service.Called(userID, password)

	return args.Error(0)
}
