package context

import (
	"github.com/boltdb/bolt"
	"github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
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
	HTTP    *http.RawReq
	Storage *bolt.DB
	// Other info for HTTP handlers can be here, like user UUID
}
