package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/waliqueiroz/letmeask-api/application/dtos"
	"github.com/waliqueiroz/letmeask-api/application/security"
	"github.com/waliqueiroz/letmeask-api/domain/repositories"
)

type AuthService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepository,
	}
}

func (service *AuthService) Login(credentials dtos.CredentialsDTO) (dtos.AuthDTO, error) {
	user, err := service.userRepository.FindByEmail(credentials.Email)
	if err != nil {
		return dtos.AuthDTO{}, err
	}

	if err := security.Verify(user.Password, credentials.Password); err != nil {
		return dtos.AuthDTO{}, err
	}

	expiresIn := time.Now().Add(time.Hour * 6).Unix()

	token, err := service.createToken(user.ID, expiresIn)
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

func (service *AuthService) createToken(userID string, expiresIn int64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = expiresIn
	permissions["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte("secret"))
}
