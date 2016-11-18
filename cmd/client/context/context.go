package context

import (
	"github.com/lastbackend/lastbackend/libs/interface/log"
)

var context Context

func Get() *Context {
	return &context
}

func Mock() *Context {
	return &context
}

type Context struct {
	Info struct {
		Version    string
		ApiVersion string
	}
	Log log.ILogger
	// Other info for HTTP handlers can be here, like user UUID
}
