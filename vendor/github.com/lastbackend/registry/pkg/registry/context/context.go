package context

import (
	"github.com/lastbackend/registry/libs/interface/log"
)

var context Context

func Get() *Context {
	return &context
}

type Context struct {
	Log log.ILogger
}
