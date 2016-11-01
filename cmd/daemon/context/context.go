package context

import (
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	r "gopkg.in/dancannon/gorethink.v2"
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
	Storage struct {
		Session *r.Session
	}
	K8S k8s.IK8S
}
