package config

import (
	"fmt"
	"github.com/lastbackend/lastbackend/libs/log/filesystem"
	"io/ioutil"
	"os"
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
	config.Token, _ = getToken()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &config
}

func getToken() (string, error) {
	tokenFile, err := os.Open(config.StoragePath + "token")
	if err != nil {
		return "", err
	}
	defer tokenFile.Close()

	fileContent, err := ioutil.ReadAll(tokenFile)
	if err != nil {
		return "", err
	}

	return string(fileContent), err
}
