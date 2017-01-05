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
