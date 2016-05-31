package env

import "github.com/deployithq/deployit/drivers/interfaces"

type Env struct {
	Log        interfaces.ILog
	Containers interfaces.IContainers
	DB         interfaces.IDB
	Host       string
}
