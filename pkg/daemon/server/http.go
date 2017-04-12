package server

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	namespace "github.com/lastbackend/lastbackend/pkg/daemon/namespace/routes"
	service "github.com/lastbackend/lastbackend/pkg/daemon/service/routes"
	vendors "github.com/lastbackend/lastbackend/pkg/daemon/vendors/routes"
	"github.com/lastbackend/lastbackend/pkg/util/http"
)

func Listen(host string, port int) error {

	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Listen HTTP server")

	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(http.Headers)

	for _, route := range namespace.Routes {
		log.Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range service.Routes {
		log.Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}

	for _, route := range vendors.Routes {
		log.Debugf("Init route: %s", route.Path)
		router.Handle(route.Path, http.Handle(route.Handler, route.Middleware...)).Methods(route.Method)
	}
	return http.Listen(host, port, router)
}
