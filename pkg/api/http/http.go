//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/util/http"

	"github.com/lastbackend/lastbackend/pkg/api/http/cluster"
	"github.com/lastbackend/lastbackend/pkg/api/http/deployment"
	"github.com/lastbackend/lastbackend/pkg/api/http/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/http/node"
	"github.com/lastbackend/lastbackend/pkg/api/http/route"
	"github.com/lastbackend/lastbackend/pkg/api/http/secret"
	"github.com/lastbackend/lastbackend/pkg/api/http/service"
	"github.com/lastbackend/lastbackend/pkg/api/http/trigger"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logLevel  = 2
	logPrefix = "api:http"
)

// Extends routes variable
var Routes = make([]http.Route, 0)

func AddRoutes(r ...[]http.Route) {
	for i := range r {
		Routes = append(Routes, r[i]...)
	}
}

func init() {

	// Cluster
	AddRoutes(cluster.Routes)
	AddRoutes(node.Routes)

	// Namespace
	AddRoutes(namespace.Routes)
	AddRoutes(service.Routes)
	AddRoutes(deployment.Routes)
	AddRoutes(route.Routes)
	AddRoutes(secret.Routes)

	// Hooks
	AddRoutes(trigger.Routes)
}

func Listen(host string, port int) error {

	log.V(logLevel).Debugf("%s:> listen HTTP server on %s:%d", logPrefix, host, port)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(http.Headers)

	for _, route := range Routes {
		log.V(logLevel).Debugf("%s:> init route: %s", logPrefix, route.Path)
		r.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	return http.Listen(host, port, r)
}
