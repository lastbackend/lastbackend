package routes

import (
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/oauth/{vendor}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Context}, Handler: OAuthDisconnectH},
	{Path: "/oauth/{vendor}/{code}", Method: http.MethodPost, Middleware: []http.Middleware{middleware.Context}, Handler: OAuthConnectH},

	// VCS handlers extends
	{Path: "/vcs/{vendor}/repos", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: VCSRepositoryListH},
	{Path: "/vcs/{vendor}/branches", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: VCSBranchListH},

	// Integration routes.
	{Path: "/integrations", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: IntegrationsH},

	{Path: "/docker/repo/search", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: DockerRepositorySearchH},
	{Path: "/docker/repo/tags", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Context}, Handler: DockerRepositoryTagListH},
}
