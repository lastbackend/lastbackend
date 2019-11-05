//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package task

import (
	"github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/namespace/{namespace}/job/{job}/task", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Authenticate}, Handler: TaskCreateH},
	{Path: "/namespace/{namespace}/job/{job}/task", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: TaskListH},
	{Path: "/namespace/{namespace}/job/{job}/task/{task}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: TaskInfoH},
	{Path: "/namespace/{namespace}/job/{job}/task/{task}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: TaskCancelH},
	{Path: "/namespace/{namespace}/job/{job}/task/{task}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Authenticate}, Handler: TaskRemoveH},
	{Path: "/namespace/{namespace}/job/{job}/task/{task}/logs", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: TaskLogsH},
}
