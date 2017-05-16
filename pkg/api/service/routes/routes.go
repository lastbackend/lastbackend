//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

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
	{Path: "/namespace/{namespace}/service/{service}/spec", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceSpecCreateH},
	{Path: "/namespace/{namespace}/service/{service}/spec/{spec}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceSpecUpdateH},
	{Path: "/namespace/{namespace}/service/{service}/spec/{spec}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceSpecRemoveH},
	{Path: "/namespace/{namespace}/service/{service}/activity", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceActivityListH},
	{Path: "/namespace/{namespace}/service/{service}/logs", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceLogsH},

	{Path: "/namespace/{namespace}/watch", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: ServiceWatchH},
}
