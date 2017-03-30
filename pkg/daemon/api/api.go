package api

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/daemon/api/routes"
	"github.com/lastbackend/lastbackend/pkg/util/http"
)

func Listen(host string, port int) error {
	router := mux.NewRouter()
	for _, route := range Routes {
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}
	return http.Listen(host, port, router)
}

var Routes = []http.Route{
	{Path: "/session", Handler: routes.SessionCreateH, Method: http.MethodPost},
	// User handlers
	{Path: "/user", Handler: routes.UserGetH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},

	// Vendor handlers
	{Path: "/oauth/{vendor}", Handler: routes.OAuthDisconnectH, Method: http.MethodDelete, Middleware: []http.Middleware{http.Authenticate}},
	{Path: "/oauth/{vendor}/{code}", Handler: routes.OAuthConnectH, Method: http.MethodPost, Middleware: []http.Middleware{http.Authenticate}},

	// VCS handlers extends
	{Path: "/vcs/{vendor}/repos", Handler: routes.VCSRepositoriesListH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},
	{Path: "/vcs/{vendor}/branches", Handler: routes.VCSBranchesListH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},

	// Build handlers
	{Path: "/build", Handler: routes.BuildListH, Method: http.MethodGet},
	{Path: "/build", Handler: routes.BuildCreateH, Method: http.MethodPost},

	// Project handlers
	{Path: "/project", Handler: routes.ProjectListH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},
	{Path: "/project", Handler: routes.ProjectCreateH, Method: http.MethodPost, Middleware: []http.Middleware{http.Authenticate}},
	{Path: "/project/{project}", Handler: routes.ProjectInfoH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},
	//{ Path: "/project/{project}", Handler: routes.routes.ProjectUpdateH, auth, Method: http.MethodPut},
	{Path: "/project/{project}", Handler: routes.ProjectRemoveH, Method: http.MethodDelete, Middleware: []http.Middleware{http.Authenticate}},
	//{ Path: "/project/{project}/activity", Handler: routes.routes.ProjectActivityListH, auth, Method: http.MethodGet},
	{Path: "/project/{project}/service", Handler: routes.ServiceListH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},
	{Path: "/project/{project}/service/{service}", Handler: routes.ServiceInfoH, Method: http.MethodGet, Middleware: []http.Middleware{http.Authenticate}},
	//{ Path: "/project/{project}/service/{service}", Handler: routes.routes.ServiceUpdateH, auth, Method: http.MethodPut},
	{Path: "/project/{project}/service/{service}", Handler: routes.ServiceRemoveH, Method: http.MethodDelete, Middleware: []http.Middleware{http.Authenticate}},

	//{ Path: "/project/{project}/service/{service}/activity", Handler: routes.routes.ServiceActivityListH, auth, Method: http.MethodGet},
	//{ Path: "/project/{project}/service/{service}/hook", Handler: routes.routes.ServiceHookCreateH, auth, Method: http.MethodPost},
	//{ Path: "/project/{project}/service/{service}/hook", Handler: routes.routes.ServiceHookListH, auth, Method: http.MethodGet},
	//{ Path: "/project/{project}/service/{service}/hook/{hook}", Handler: routes.routes.ServiceHookRemoveH, auth, Method: http.MethodDelete},
	//{ Path: "/project/{project}/service/{service}/logs", Handler: routes.routes.ServiceLogsH, auth, Method: http.MethodGet},

	// Deploy template/docker/source/repo
	//{ Path: "/deploy", Handler: routes.routes.DeployH, auth, Method: http.MethodPost},

	// Hook routes.
	{Path: "/hook/{token}", Handler: routes.HookExecuteH, Method: http.MethodPost},

	// Docker routes.
	{Path: "/docker/repo/search", Handler: routes.DockerRepositorySearchH, Method: http.MethodGet},
	{Path: "/docker/repo/tags", Handler: routes.DockerRepositoryTagListH, Method: http.MethodGet},
}
