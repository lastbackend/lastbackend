package context

import (
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
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
	K8S     k8s.IK8S
	Storage storage.IStorage
}
