package filesystem

import (
	"os"
	"os/user"
)


func GetHomeDir() (homePath string, err error) {
	homePath = os.Getenv("HOME")
	if homePath != "" {
		return homePath, err
	}
	user, err := user.Current()
	return user.HomeDir, err
}

func MkDir(path string) (err error) {

	err = os.MkdirAll(path, 0755)
	return err

}
