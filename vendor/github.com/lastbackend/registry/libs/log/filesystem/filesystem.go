package filesystem

import (
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
)

const _WINDOWS = "windows"

// Check OS Windows
func IsWindows() bool {

	return runtime.GOOS == _WINDOWS

}

// GetHomeDir is used to get user home directory (return: path home directory, error message)
func GetHomeDir() (homePath string, err error) {

	homePath = os.Getenv("HOME")
	if homePath != "" {
		return homePath, err
	}
	user, err := user.Current()
	return user.HomeDir, err

}

// GetPerm is used to get bit permissions
func GetPerm(path string) (perm uint8) {

	return perm

}

// MkDir is used to create directory
func MkDir(path string) (err error) {

	err = os.MkdirAll(path, 0755)
	return err

}

// CreateFile is used to create file
func CreateFile(path string) (err error) {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil

}

func WriteStrToFile(path, value string) (err error) {

	err = ioutil.WriteFile(path, []byte(value), 0755)
	if err != nil {
		if os.IsNotExist(err) {
			CreateFile(path)
		}
		return err
	}
	return nil

}

func ReadFile(path string) (value []byte, err error) {

	value, err = ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return value, nil

}
