package config

// The structure of the config to run the client
type Config struct {
	Debug       bool   `yaml:"debug"`
	UserUrl     string `yaml:"createUserUrl"`
	StoragePath string `yaml:"storagePah"`
	AuthUserUrl string `yaml:"authUserUrl"`
}
