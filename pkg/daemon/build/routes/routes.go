package routes

import (
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/build", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: BuildListH},
}
