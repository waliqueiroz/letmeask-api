package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
)

type AuthServiceMock struct {
	mock.Mock
}

func NewAuthServiceMock() *AuthServiceMock {
	return &AuthServiceMock{}
}

func (service *AuthServiceMock) Login(credentials dtos.CredentialsDTO) (dtos.AuthDTO, error) {
	args := service.Called(credentials)

	return args.Get(0).(dtos.AuthDTO), args.Error(1)
}
