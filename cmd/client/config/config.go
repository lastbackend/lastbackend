package config

import (
	"github.com/lastbackend/lastbackend/libs/log/filesystem"
	"fmt"
)

var config Config

func Get() *Config {

	homedir, err := filesystem.GetHomeDir()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	config.StoragePath = homedir + "/token.txt"
	config.CreateUserUrl = "http://localhost:3000/user"
	return &config
}
