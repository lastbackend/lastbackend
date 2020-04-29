//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package server

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/server/config"
	"github.com/lastbackend/lastbackend/internal/server/server/cluster"
	"github.com/lastbackend/lastbackend/internal/server/server/discovery"
	"github.com/lastbackend/lastbackend/internal/server/server/events"
	"github.com/lastbackend/lastbackend/internal/server/server/exporter"
	"github.com/lastbackend/lastbackend/internal/server/server/ingress"
	"github.com/lastbackend/lastbackend/internal/server/server/middleware"
	"github.com/lastbackend/lastbackend/internal/server/server/namespace"
	"github.com/lastbackend/lastbackend/internal/server/server/node"
	"github.com/lastbackend/lastbackend/internal/server/state"
	"github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/cors"
)

const defaultPort = 2967

type HttpServer struct {
	host string
	port uint

	Insecure bool

	CertFile string
	KeyFile  string
	CaFile   string

	router *mux.Router
}

func NewServer(state *state.State, stg storage.IStorage, cfg config.Config) *HttpServer {

	hs := new(HttpServer)

	if cfg.Server.TLS.Verify {
		hs.CertFile = cfg.Server.TLS.FileCert
		hs.KeyFile = cfg.Server.TLS.FileKey
		hs.CaFile = cfg.Server.TLS.FileCA
	} else {
		hs.Insecure = true
	}

	hs.host = cfg.Server.Host
	hs.port = cfg.Server.Port
	if hs.port == 0 {
		hs.port = defaultPort
	}

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(cors.Headers)

	var notFound http.MethodNotAllowedHandler
	r.NotFoundHandler = notFound

	var notAllowed http.MethodNotAllowedHandler
	r.MethodNotAllowedHandler = notAllowed

	mw := middleware.New(stg, cfg.Security.Token)

	cluster.NewClusterHandler(r, mw, cluster.Config{ClusterName: cfg.ClusterName, ClusterDescription: cfg.ClusterDescription})
	//config.NewConfigHandler(r, mw)
	//deployment.NewDeploymentHandler(r, mw)
	discovery.NewDiscoveryHandler(r, mw)
	events.NewEventHandler(r, mw)
	exporter.NewExporterHandler(r, mw)
	ingress.NewIngressHandler(r, mw)
	//job.NewJobHandler(r, mw, job.Config{SecretToken: v.GetString("security.token")})
	namespace.NewNamespaceHandler(r, mw, state)
	node.NewNodeHandler(r, mw)

	//pod.NewPodHandler(r, mw)
	//route.NewRouteHandler(r, mw)
	//secret.NewSecretHandler(r, mw, secret.Config{Vault: &types.Vault{
	//	Endpoint: v.GetString("vault.endpoint"),
	//	Token:    v.GetString("vault.token"),
	//}})
	//service.NewServiceHandler(r, mw, service.Config{SecretToken: v.GetString("security.token")})
	//task.NewTaskHandler(r, mw)
	//volume.NewVolumeHandler(r, mw)

	hs.router = r

	return hs
}

func (hs HttpServer) Run() error {

	if hs.Insecure {
		return http.Listen(hs.host, hs.port, hs.router)
	}

	return http.ListenWithTLS(hs.host, hs.port, hs.CaFile, hs.CertFile, hs.KeyFile, hs.router)
}
