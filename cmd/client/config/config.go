package config

import (
	"fmt"
	"github.com/lastbackend/lastbackend/libs/log/filesystem"
//	"io/ioutil"
)

var config Config

func Get() *Config {

	homedir, err := filesystem.GetHomeDir()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	config.StoragePath = homedir + "/.lb/"
	config.UserUrl = "http://localhost:3000/user"
	config.AuthUserUrl = "http://localhost:3000/session"
	config.ProjectUrl = "http://localhost:3000/project"



	return &config
}
