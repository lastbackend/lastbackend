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
	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/agent/server/middleware"
	"github.com/lastbackend/lastbackend/internal/agent/server/node"
	"github.com/lastbackend/lastbackend/internal/agent/server/pod"
	"github.com/lastbackend/lastbackend/internal/agent/state"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/cors"
)

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

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(cors.Headers)

	var notFound http.MethodNotAllowedHandler
	r.NotFoundHandler = notFound

	var notAllowed http.MethodNotAllowedHandler
	r.MethodNotAllowedHandler = notAllowed

	mw := middleware.New(stg, cfg.Security.Token)

	node.NewNodeHandler(r, mw)
	pod.NewPodHandler(r, mw, state)

	hs.router = r

	return hs
}

func (hs HttpServer) Run() error {

	if hs.Insecure {
		return http.Listen(hs.host, hs.port, hs.router)
	}

	return http.ListenWithTLS(hs.host, hs.port, hs.CaFile, hs.CertFile, hs.KeyFile, hs.router)
}
