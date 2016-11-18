package config

var config Config

func Get() *Config {
	config.StoragePath = "token.txt"
	config.CreateUserUrl = "http://localhost:3000/user"
	return &config
}
