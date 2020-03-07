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

package http

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/master/http/cluster"
	"github.com/lastbackend/lastbackend/internal/master/http/config"
	"github.com/lastbackend/lastbackend/internal/master/http/deployment"
	"github.com/lastbackend/lastbackend/internal/master/http/discovery"
	"github.com/lastbackend/lastbackend/internal/master/http/events"
	"github.com/lastbackend/lastbackend/internal/master/http/exporter"
	"github.com/lastbackend/lastbackend/internal/master/http/ingress"
	"github.com/lastbackend/lastbackend/internal/master/http/job"
	"github.com/lastbackend/lastbackend/internal/master/http/namespace"
	"github.com/lastbackend/lastbackend/internal/master/http/node"
	"github.com/lastbackend/lastbackend/internal/master/http/pod"
	"github.com/lastbackend/lastbackend/internal/master/http/route"
	"github.com/lastbackend/lastbackend/internal/master/http/secret"
	"github.com/lastbackend/lastbackend/internal/master/http/service"
	"github.com/lastbackend/lastbackend/internal/master/http/task"
	"github.com/lastbackend/lastbackend/internal/master/http/volume"
	"github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/cors"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 2
	logPrefix = "api:http"
)

// Extends routes variable
var Routes = make([]http.Route, 0)

type HttpOpts struct {
	Insecure bool

	CertFile string
	KeyFile  string
	CaFile   string

	BearerToken string
}

func AddRoutes(r ...[]http.Route) {
	for i := range r {
		Routes = append(Routes, r[i]...)
	}
}

func init() {

	// Cluster
	AddRoutes(cluster.Routes)
	AddRoutes(node.Routes)
	AddRoutes(ingress.Routes)
	AddRoutes(exporter.Routes)
	AddRoutes(discovery.Routes)

	// Environment
	AddRoutes(namespace.Routes)
	AddRoutes(secret.Routes)
	AddRoutes(config.Routes)
	AddRoutes(route.Routes)
	AddRoutes(service.Routes)
	AddRoutes(deployment.Routes)
	AddRoutes(pod.Routes)
	AddRoutes(volume.Routes)
	AddRoutes(ingress.Routes)
	AddRoutes(job.Routes)
	AddRoutes(task.Routes)

	// events
	AddRoutes(events.Routes)
}

func Listen(host string, port int, opts *HttpOpts) error {

	if opts == nil {
		opts = new(HttpOpts)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "access_token", opts.BearerToken)

	log.V(logLevel).Debugf("%s:> listen HTTP server on %s:%d", logPrefix, host, port)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(cors.Headers)

	var notFound http.MethodNotAllowedHandler
	r.NotFoundHandler = notFound

	var notAllowed http.MethodNotAllowedHandler
	r.MethodNotAllowedHandler = notAllowed

	for _, rt := range Routes {
		log.V(logLevel).Debugf("%s:> init route: %s", logPrefix, rt.Path)
		r.Handle(rt.Path, http.Handle(ctx, rt.Handler, rt.Middleware...)).Methods(rt.Method)
	}

	if len(opts.CaFile) == 0 || len(opts.CertFile) == 0 || len(opts.KeyFile) == 0 {
		log.V(logLevel).Debugf("%s:> run insecure http server on %d port", logPrefix, port)
		return http.Listen(host, port, r)
	}

	log.V(logLevel).Debugf("%s:> run http server with tls on %d port", logPrefix, port)
	return http.ListenWithTLS(host, port, opts.CaFile, opts.CertFile, opts.KeyFile, r)
}
