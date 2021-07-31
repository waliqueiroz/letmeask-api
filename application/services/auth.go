package services

import (
	"time"

	"github.com/waliqueiroz/letmeask-api/application/dtos"
	"github.com/waliqueiroz/letmeask-api/application/providers"
	"github.com/waliqueiroz/letmeask-api/domain/repositories"
)

type AuthService struct {
	userRepository   repositories.UserRepository
	securityProvider providers.SecurityProvider
	authProvider     providers.AuthProvider
}

func NewAuthService(userRepository repositories.UserRepository, securityProvider providers.SecurityProvider, authProvider providers.AuthProvider) *AuthService {
	return &AuthService{
		userRepository,
		securityProvider,
		authProvider,
	}
}

func (service *AuthService) Login(credentials dtos.CredentialsDTO) (dtos.AuthDTO, error) {
	user, err := service.userRepository.FindByEmail(credentials.Email)
	if err != nil {
		return dtos.AuthDTO{}, err
	}

	if err := service.securityProvider.Verify(user.Password, credentials.Password); err != nil {
		return dtos.AuthDTO{}, err
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
