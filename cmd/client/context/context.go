package context

import (
	"github.com/lastbackend/lastbackend/libs/interface/log"
	"github.com/lastbackend/lastbackend/libs/interface/HTTP"
	//"github.com/lastbackend/lastbackend/libs/http"
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
	HTTP http.IHTTP
	Log log.ILogger
	// Other info for HTTP handlers can be here, like user UUID
}
