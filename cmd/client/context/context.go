package context

import (
	"github.com/lastbackend/api/libs/interface/http"
	"github.com/lastbackend/lastbackend/libs/interface/log"
<<<<<<< HEAD
	"github.com/lastbackend/lastbackend/libs/interface/sdb"
=======
	"github.com/lastbackend/lastbackend/libs/interface/HTTP"
	//"github.com/lastbackend/lastbackend/libs/http"
>>>>>>> upstream/master
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
<<<<<<< HEAD
	Log     log.ILogger
	HTTP    http.IHTTP
	Storage sdb.ISDB
=======
	HTTP http.IHTTP
	Log log.ILogger
>>>>>>> upstream/master
	// Other info for HTTP handlers can be here, like user UUID
}
