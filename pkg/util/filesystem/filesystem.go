//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package filesystem

import (
	"io/ioutil"
	"os"
	"runtime"
)

const _WINDOWS = "windows"

// Check OS Windows
func IsWindows() bool {
	return runtime.GOOS == _WINDOWS
}

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
