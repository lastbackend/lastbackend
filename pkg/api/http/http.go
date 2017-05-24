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

func Listen(host string, port int) error {

	var (
		log = context.Get().GetLogger()
	)

	log.V(logLevel).Debug("Listen HTTP server")

	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(http.Headers)

	for _, route := range namespace.Routes {
		log.V(logLevel).Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range service.Routes {
		log.V(logLevel).Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range vendors.Routes {
		log.V(logLevel).Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range node.Routes {
		log.V(logLevel).Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range events.Routes {
		log.V(logLevel).Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range hook.Routes {
		log.V(logLevel).Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	return http.Listen(host, port, router)
}
