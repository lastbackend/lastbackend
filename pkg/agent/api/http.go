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
	"github.com/lastbackend/lastbackend/pkg/agent/api/routes"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

func Listen(host string, port int) error {

	log := context.Get().GetLogger()
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
	// Vendor handlers
	{Path: "/patch", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Context},
		Handler: routes.SetPods},
}