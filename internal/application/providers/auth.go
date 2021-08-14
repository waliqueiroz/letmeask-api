package providers

type AuthProvider interface {
	CreateToken(userID string, expiresIn int64) (string, error)
}
