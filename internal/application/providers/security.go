package providers

type SecurityProvider interface {
	Hash(password string) (string, error)
	Verify(hashedPassword string, password string) error
}
