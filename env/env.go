package env

import "github.com/deployithq/deployit/drivers/interfaces"

type Env struct {
	Log         interfaces.ILog
	Storage     interfaces.IStorage
	Host        string
	Path        string
	StoragePath string
}
