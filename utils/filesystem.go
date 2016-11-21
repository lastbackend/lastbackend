package utils

import (
	"os"
	"os/user"
)

func GetHomeDir() string {

	var dir string

	usr, err := user.Current()
	if err == nil {
		dir = usr.HomeDir
	} else {
		// Maybe it's cross compilation without cgo support. (darwin, unix)
		dir = os.Getenv("HOME")
	}

	return dir
}
