package context

import (
	"github.com/lastbackend/api/libs/interface/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	"github.com/lastbackend/lastbackend/libs/interface/sdb"
)

var context Context
var mock Context

func Get() *Context {
	return &context
}

func Mock() *Context {
	return &mock
}

type Context struct {
	Info struct {
		Version    string
		ApiVersion string
	}
	Log     log.ILogger
	HTTP    http.IHTTP
	Storage sdb.ISDB
	// Other info for HTTP handlers can be here, like user UUID
}
