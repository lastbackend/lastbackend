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
	"github.com/lastbackend/lastbackend/internal/minion/server/middleware"
	"github.com/lastbackend/lastbackend/internal/minion/server/node"
	"github.com/lastbackend/lastbackend/internal/minion/server/pod"
	"github.com/lastbackend/lastbackend/internal/minion/state"
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

func NewServer(state *state.State, stg storage.Storage, v *viper.Viper) *HttpServer {

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

	node.NewNodeHandler(r, mw)
	pod.NewPodHandler(r, mw)

	hs.router = r

	return hs
}

func (hs HttpServer) Run() error {

	if hs.Insecure {
		return http.Listen(hs.host, hs.port, hs.router)
	}

	return http.ListenWithTLS(hs.host, hs.port, hs.CaFile, hs.CertFile, hs.KeyFile, hs.router)
}
