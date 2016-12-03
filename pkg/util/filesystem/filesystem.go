package filesystem

import (
	"io/ioutil"
	"os"
)

// MkDir is used to create directory
func MkDir(path string, mode os.FileMode) (err error) {
	err = os.MkdirAll(path, mode)
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

func WriteStrToFile(path, value string, mode os.FileMode) (err error) {
	err = ioutil.WriteFile(path, []byte(value), mode)
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
