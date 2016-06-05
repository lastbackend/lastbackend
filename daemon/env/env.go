package env

import (
	"github.com/deployithq/deployit/drivers/interfaces"
)

const (
	Default_root_path string = "/var/lib/deployit"
	Default_hub string = "hub.deployit.io"
)

type Env struct {
	Log        interfaces.ILog
	LDB        interfaces.ILDB
	Containers interfaces.IContainers
	Host       string
}
