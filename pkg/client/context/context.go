package context

import (
	"github.com/lastbackend/lastbackend/libs/db"
	"github.com/lastbackend/lastbackend/libs/http"
	d "github.com/lastbackend/lastbackend/libs/interface/db"
	"github.com/lastbackend/lastbackend/libs/interface/log"
	l "github.com/lastbackend/lastbackend/libs/log"
)

var context Context

func Get() *Context {
	return &context
}

func Mock() *Context {
	context.mock = true
	context.Log = new(l.Log)
	context.Log.Init()
	context.Log.Disabled()
	context.Storage = new(db.DB)

	return &context
}

type Context struct {
	Token   string
	Log     log.ILogger
	HTTP    *http.RawReq
	Storage d.IDB
	mock    bool
	// Other info for HTTP handlers can be here, like user UUID
}
