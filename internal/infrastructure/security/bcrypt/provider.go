package bcrypt

import "golang.org/x/crypto/bcrypt"

type BcryptProvider struct{}

func NewBcryptProvider() *BcryptProvider {
	return &BcryptProvider{}
}

func (provider *BcryptProvider) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (provider *BcryptProvider) Verify(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
