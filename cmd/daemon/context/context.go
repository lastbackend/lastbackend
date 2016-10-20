package context

import (
	"github.com/deployithq/deployit/libs/interface/log"
	"github.com/deployithq/deployit/libs/interface/k8s"
)

var context Context

func Get() *Context {
	return &context
}

type Context struct {
	Version string
	Log     log.ILogger
	K8S     k8s.IK8S
	// Other info for HTTP handlers can be here, like user UUID
}
