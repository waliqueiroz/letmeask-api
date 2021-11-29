package providers

type Authenticator interface {
	CreateToken(userID string, expiresIn int64) (string, error)
	ExtractUserID(token interface{}) (string, error)
}
