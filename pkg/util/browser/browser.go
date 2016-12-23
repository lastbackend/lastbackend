package browser

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Open(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		cmd := "url.dll,FileProtocolHandler"
		runDll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")

    err = exec.Command(runDll32, cmd, url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errors.New("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil
}
