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

package http

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/util/http"

	events "github.com/lastbackend/lastbackend/pkg/api/events/routes"
	hook "github.com/lastbackend/lastbackend/pkg/api/hook/routes"
	namespace "github.com/lastbackend/lastbackend/pkg/api/namespace/routes"
	node "github.com/lastbackend/lastbackend/pkg/api/node/routes"
	service "github.com/lastbackend/lastbackend/pkg/api/service/routes"
	vendors "github.com/lastbackend/lastbackend/pkg/api/vendors/routes"
)

const logLevel = 2

// Extends routes variable
var Routes = make([]http.Route, 0)

func AddRoutes(r ...[]http.Route) {
	for i := range r {
		Routes = append(Routes, r[i]...)
	}
}

func init() {
	AddRoutes(events.Routes)
	AddRoutes(hook.Routes)
	AddRoutes(namespace.Routes)
	AddRoutes(node.Routes)
	AddRoutes(service.Routes)
	AddRoutes(vendors.Routes)
}

func Listen(host string, port int) error {

	log := context.Get().GetLogger()
	log.V(logLevel).Debug("HTTP: listen HTTP server")

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(http.Headers)

	for _, route := range Routes {
		log.V(logLevel).Debugf("HTTP: init route: %s", route.Path)
		r.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	return http.Listen(host, port, r)
}
