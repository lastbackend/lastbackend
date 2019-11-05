//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package browser

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var Os = runtime.GOOS

var CommandWrapper = func(name string, parameters ...string) error {
	return exec.Command(name, parameters...).Start()
}

func Open(url string) error {
	var err error

	switch Os {
	case "linux":
		err = CommandWrapper("xdg-open", url)
	case "windows":
		cmd := "url.dll,FileProtocolHandler"
		runDll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
		err = CommandWrapper(runDll32, cmd, url)
	case "darwin":
		err = CommandWrapper("open", url)
	default:
		err = errors.New("unsupported platform")
	}

	return err
}
