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
	{Path: "/vendor", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: IntegrationsH},
	{Path: "/vendor/{vendor}/oauth", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: OAuthDisconnectH},
	{Path: "/vendor/{vendor}/oauth/{code}", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: OAuthConnectH},
	{Path: "/vendor/{vendor}/repos", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: VCSRepositoryListH},
	{Path: "/vendor/{vendor}/branches", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: VCSBranchListH},
	{Path: "/vendor/docker/search", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: DockerRepositorySearchH},
	{Path: "/vendor/docker/tags", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: DockerRepositoryTagListH},
}
