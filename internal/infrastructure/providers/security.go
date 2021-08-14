package providers

import "golang.org/x/crypto/bcrypt"

type SecurityProvider struct{}

func NewSecurityProvider() *SecurityProvider {
	return &SecurityProvider{}
}

func (provider *SecurityProvider) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (provider *SecurityProvider) Verify(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
