package env

import "github.com/deployithq/deployit/drivers/interfaces"

type Env struct {
	Log         interfaces.ILog
	Storage     interfaces.IStorage
	HostUrl     string
	Host        string
	Path        string
	StoragePath string
	LogMode     bool
	Port        int
}
