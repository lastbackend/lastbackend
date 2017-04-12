package routes

import (
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/namespace/{namespace}/service", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceCreateH},
	{Path: "/namespace/{namespace}/service", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceListH},
	{Path: "/namespace/{namespace}/service/{service}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceInfoH},
	{Path: "/namespace/{namespace}/service/{service}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceUpdateH},
	{Path: "/namespace/{namespace}/service/{service}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceRemoveH},
	{Path: "/namespace/{namespace}/service/{service}/activity", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceActivityListH},
	{Path: "/namespace/{namespace}/service/{service}/logs", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceLogsH},
}
