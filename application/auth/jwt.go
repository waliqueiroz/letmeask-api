package auth

import "github.com/golang-jwt/jwt"

func CreateToken(userID string, expiresIn int64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = expiresIn
	permissions["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte("secret"))
}
