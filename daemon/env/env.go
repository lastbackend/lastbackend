package env

import "github.com/deployithq/deployit/drivers/interfaces"

const (
	Default_root_path string = "/var/lib/deployit"
)

type Env struct {
	Log        interfaces.ILog
	Containers interfaces.IContainers
	Host       string
}
