package config

// The structure of the config to run the client
type Config struct {
	Debug         bool   `yaml:"debug"`
	CreateUserUrl string `yaml:"createUserUrl"`
	StoragePath   string `yaml:"storagePah"`
}
