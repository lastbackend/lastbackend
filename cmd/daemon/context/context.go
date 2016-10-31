package context

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	lib_storage "github.com/lastbackend/lastbackend/libs/interface/storage"
)

var context Context

func Get() *Context {
	return &context
}

type Context struct {
	Info struct {
		Version    string
		ApiVersion string
	}
	Log     log.ILogger
	Storage adapter.IStorage
	K8S     k8s.IK8S
	Adapter lib_storage.Storage
}
