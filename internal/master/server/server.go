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
	"github.com/lastbackend/lastbackend/internal/master/server/cluster"
	"github.com/lastbackend/lastbackend/internal/master/server/config"
	"github.com/lastbackend/lastbackend/internal/master/server/deployment"
	"github.com/lastbackend/lastbackend/internal/master/server/discovery"
	"github.com/lastbackend/lastbackend/internal/master/server/events"
	"github.com/lastbackend/lastbackend/internal/master/server/exporter"
	"github.com/lastbackend/lastbackend/internal/master/server/ingress"
	"github.com/lastbackend/lastbackend/internal/master/server/job"
	"github.com/lastbackend/lastbackend/internal/master/server/middleware"
	"github.com/lastbackend/lastbackend/internal/master/server/namespace"
	"github.com/lastbackend/lastbackend/internal/master/server/node"
	"github.com/lastbackend/lastbackend/internal/master/server/pod"
	"github.com/lastbackend/lastbackend/internal/master/server/route"
	"github.com/lastbackend/lastbackend/internal/master/server/secret"
	"github.com/lastbackend/lastbackend/internal/master/server/service"
	"github.com/lastbackend/lastbackend/internal/master/server/task"
	"github.com/lastbackend/lastbackend/internal/master/server/volume"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/cors"
	"github.com/spf13/viper"
)

type HttpServer struct {
	host string
	port int

	Insecure bool

	CertFile string
	KeyFile  string
	CaFile   string

	router *mux.Router
}

func NewServer(stg storage.Storage, v *viper.Viper) *HttpServer {

	hs := new(HttpServer)

	if v.GetBool("server.tls") {
		hs.CertFile = v.GetString("server.tls.cert")
		hs.KeyFile = v.GetString("server.tls.key")
		hs.CaFile = v.GetString("server.tls.ca")
	} else {
		hs.Insecure = true
	}

	hs.host = v.GetString("server.host")
	hs.port = v.GetInt("server.port")

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(cors.Headers)

	var notFound http.MethodNotAllowedHandler
	r.NotFoundHandler = notFound

	var notAllowed http.MethodNotAllowedHandler
	r.MethodNotAllowedHandler = notAllowed

	mw := middleware.New(stg, v)

	cluster.NewClusterHandler(r, mw)
	config.NewConfigHandler(r, mw)
	deployment.NewDeploymentHandler(r, mw)
	discovery.NewDiscoveryHandler(r, mw)
	events.NewEventHandler(r, mw)
	exporter.NewExporterHandler(r, mw)
	ingress.NewIngressHandler(r, mw)
	job.NewJobHandler(r, mw)
	namespace.NewNamespaceHandler(r, mw)
	node.NewNodeHandler(r, mw)
	pod.NewPodHandler(r, mw)
	route.NewRouteHandler(r, mw)
	secret.NewSecretHandler(r, mw)
	service.NewServiceHandler(r, mw)
	task.NewTaskHandler(r, mw)
	volume.NewVolumeHandler(r, mw)

	hs.router = r

	return hs
}

func (hs HttpServer) Run() error {

	if hs.Insecure {
		return http.Listen(hs.host, hs.port, hs.router)
	}

	return http.ListenWithTLS(hs.host, hs.port, hs.CaFile, hs.CertFile, hs.KeyFile, hs.router)
}
