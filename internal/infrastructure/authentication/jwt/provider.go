package jwt

import (
	"github.com/golang-jwt/jwt"
	application "github.com/waliqueiroz/letmeask-api/internal/application/errors"
	"github.com/waliqueiroz/letmeask-api/internal/infrastructure/configurations"
)

type JwtProvider struct {
	configuration configurations.Configuration
}

func NewJwtProvider(configuration configurations.Configuration) *JwtProvider {
	return &JwtProvider{
		configuration,
	}
}

func (provider *JwtProvider) CreateToken(userID string, expiresIn int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["exp"] = expiresIn
	claims["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(provider.configuration.Auth.SecretKey))
}

func (provider *JwtProvider) ExtractUserID(token interface{}) (string, error) {
	err := application.NewUnauthorizedError("token inválido")

	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", err
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}
