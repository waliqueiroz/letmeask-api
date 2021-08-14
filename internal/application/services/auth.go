package services

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/application/errors"
	"github.com/waliqueiroz/letmeask-api/internal/application/providers"
	"github.com/waliqueiroz/letmeask-api/internal/domain/repositories"
)

type AuthService interface {
	Login(credentials dtos.CredentialsDTO) (dtos.AuthDTO, error)
}

type authService struct {
	userRepository   repositories.UserRepository
	securityProvider providers.SecurityProvider
	authProvider     providers.AuthProvider
}

func NewAuthService(userRepository repositories.UserRepository, securityProvider providers.SecurityProvider, authProvider providers.AuthProvider) *authService {
	return &authService{
		userRepository,
		securityProvider,
		authProvider,
	}
}

func (service *authService) Login(credentials dtos.CredentialsDTO) (dtos.AuthDTO, error) {
	user, err := service.userRepository.FindByEmail(credentials.Email)
	if err != nil {
		return dtos.AuthDTO{}, errors.NewUnauthorizedError("credenciais inválidas")
	}

	if err := service.securityProvider.Verify(user.Password, credentials.Password); err != nil {
		return dtos.AuthDTO{}, errors.NewUnauthorizedError("credenciais inválidas")
	}

	expiresIn := time.Now().Add(time.Hour * 6).Unix()

	token, err := service.authProvider.CreateToken(user.ID, expiresIn)
	if err != nil {
		return dtos.AuthDTO{}, err
	}

	return dtos.AuthDTO{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}
