package configurations

type Configuration struct {
	Database Database `mapstructure:",squash"`
	Auth     Auth     `mapstructure:",squash"`
}
