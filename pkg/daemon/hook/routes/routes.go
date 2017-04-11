package routes

import (
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/hook/{token}", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: HookExecuteH},
}
