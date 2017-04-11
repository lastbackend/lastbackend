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

package api

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/daemon/api/routes"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

func Listen(host string, port int) error {

	var (
		log = c.Get().GetLogger()
	)

	log.Debug("Listen HTTP server")

	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(http.Headers)
	for _, route := range Routes {
		log.Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}
	return http.Listen(host, port, router)
}

var Routes = []http.Route{

	{Path: "/status", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.StatusH},

	// Vendor handlers
	{Path: "/oauth/{vendor}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: routes.OAuthDisconnectH},
	{Path: "/oauth/{vendor}/{code}", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: routes.OAuthConnectH},

	// VCS handlers extends
	{Path: "/vcs/{vendor}/repos", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.VCSRepositoryListH},
	{Path: "/vcs/{vendor}/branches", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.VCSBranchListH},

	// Integration routes.
	{Path: "/integrations", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.IntegrationsH},

	// Build handlers
	{Path: "/build", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.BuildListH},

	// Project handlers
	{Path: "/project", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ProjectListH},
	{Path: "/project", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ProjectCreateH},
	{Path: "/project/{project}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ProjectInfoH},
	{Path: "/project/{project}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ProjectUpdateH},
	{Path: "/project/{project}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ProjectRemoveH},
	{Path: "/project/{project}/activity", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ProjectActivityListH},

	// Service handlers
	{Path: "/project/{project}/service", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceCreateH},
	{Path: "/project/{project}/service", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceListH},
	{Path: "/project/{project}/service/{service}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceInfoH},
	{Path: "/project/{project}/service/{service}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceUpdateH},
	{Path: "/project/{project}/service/{service}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceRemoveH},
	{Path: "/project/{project}/service/{service}/activity", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceActivityListH},
	{Path: "/project/{project}/service/{service}/logs", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: routes.ServiceLogsH},

	// Hook routes.
	{Path: "/hook/{token}", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: routes.HookExecuteH},

	// Docker routes.
	{Path: "/docker/repo/search", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.DockerRepositorySearchH},
	{Path: "/docker/repo/tags", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.DockerRepositoryTagListH},

	// Template routes.
	{Path: "/template", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: routes.TemplateListH},
}
