package providers

import (
	"github.com/golang-jwt/jwt"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
)

type AuthProvider struct {
	configuration configurations.Configuration
}

func NewAuthProvider(configuration configurations.Configuration) *AuthProvider {
	return &AuthProvider{
		configuration,
	}
}

func (provider *AuthProvider) CreateToken(userID string, expiresIn int64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = expiresIn
	permissions["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte(provider.configuration.Auth.SecretKey))
}
