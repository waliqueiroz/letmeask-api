package configurations

type Auth struct {
	SecretKey string `env:"SECRET_KEY"`
}
