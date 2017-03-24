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
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/http/handler"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Extends = make(map[string]Handler)

type Handler struct {
	Path    string
	Method  string
	Auth    bool
	Handler func(http.ResponseWriter, *http.Request)
}

func NewRouter() *mux.Router {

	var ctx = c.Get()
	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(headers)

	// Session handlers
	r.HandleFunc("/session", handle(handler.SessionCreateH)).Methods(http.MethodPost)

	// User handlers
	r.HandleFunc("/user", handle(handler.UserGetH, auth)).Methods(http.MethodGet)

	// Build handlers
	r.HandleFunc("/build", handle(handler.BuildListH)).Methods(http.MethodGet)
	r.HandleFunc("/build", handle(handler.BuildCreateH)).Methods(http.MethodPost)

	// Project handlers
	r.HandleFunc("/project", handle(handler.ProjectListH, auth)).Methods(http.MethodGet)
	r.HandleFunc("/project", handle(handler.ProjectCreateH, auth)).Methods(http.MethodPost)
	r.HandleFunc("/project/{project}", handle(handler.ProjectInfoH, auth)).Methods(http.MethodGet)
	//r.HandleFunc("/project/{project}", handle(handler.ProjectUpdateH, auth)).Methods(http.MethodPut)
	r.HandleFunc("/project/{project}", handle(handler.ProjectRemoveH, auth)).Methods(http.MethodDelete)
	//r.HandleFunc("/project/{project}/activity", handle(handler.ProjectActivityListH, auth)).Methods(http.MethodGet)
	r.HandleFunc("/project/{project}/service", handle(handler.ServiceListH, auth)).Methods(http.MethodGet)
	r.HandleFunc("/project/{project}/service/{service}", handle(handler.ServiceInfoH, auth)).Methods(http.MethodGet)
	//r.HandleFunc("/project/{project}/service/{service}", handle(handler.ServiceUpdateH, auth)).Methods(http.MethodPut)
	r.HandleFunc("/project/{project}/service/{service}", handle(handler.ServiceRemoveH, auth)).Methods(http.MethodDelete)
	//r.HandleFunc("/project/{project}/service/{service}/activity", handle(handler.ServiceActivityListH, auth)).Methods(http.MethodGet)
	//r.HandleFunc("/project/{project}/service/{service}/hook", handle(handler.ServiceHookCreateH, auth)).Methods(http.MethodPost)
	//r.HandleFunc("/project/{project}/service/{service}/hook", handle(handler.ServiceHookListH, auth)).Methods(http.MethodGet)
	//r.HandleFunc("/project/{project}/service/{service}/hook/{hook}", handle(handler.ServiceHookRemoveH, auth)).Methods(http.MethodDelete)
	//r.HandleFunc("/project/{project}/service/{service}/logs", handle(handler.ServiceLogsH, auth)).Methods(http.MethodGet)

	// Deploy template/docker/source/repo
	//r.HandleFunc("/deploy", handle(handler.DeployH, auth)).Methods(http.MethodPost)

	// Hook handlers
	r.HandleFunc("/hook/{token}", handle(handler.HookExecuteH)).Methods(http.MethodPost)

	// Docker handlers
	r.HandleFunc("/docker/repo/search", handle(handler.DockerRepositorySearchH)).Methods(http.MethodGet)
	r.HandleFunc("/docker/repo/tags", handle(handler.DockerRepositoryTagListH)).Methods(http.MethodGet)

	ctx.Log.Info("Extends API methods:")

	for name, h := range Extends {
		// TODO: Check path on correctly
		ctx.Log.Info(name)
		if h.Auth {
			r.HandleFunc(h.Path, handle(h.Handler, auth)).Methods(h.Method)
		} else {
			r.HandleFunc(h.Path, handle(h.Handler)).Methods(h.Method)
		}
	}

	return r
}

func RunHttpServer(routes *mux.Router, port int) {

	var ctx = c.Get()

	ctx.Log.Infof("Listen http server on %d port", port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), routes); err != nil {
		ctx.Log.Fatal("ListenAndServe: ", err)
	}
}

func headers(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "X-CSRF-Token, Authorization, Content-Type, x-lastbackend, Origin, X-Requested-With, Content-Name, Accept")
	w.Header().Add("Content-Type", "application/json")
}

func handle(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	headers := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			headers(w, r)
			h.ServeHTTP(w, r)

			fmt.Println(fmt.Sprintf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start)))
		}
	}

	h = headers(h)
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

// Auth - authentication middleware
func auth(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var token string
		var params = mux.Vars(r)

		if _, ok := params["token"]; ok {
			token = params["token"]
		} else if r.Header.Get("Authorization") != "" {
			// Parse authorization header
			var auth = strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			// Check authorization header parts length and authorization header format
			if len(auth) != 2 || auth[0] != "Bearer" {
				e.HTTP.Unauthorized(w)
				return
			}

			token = auth[1]

		} else {
			w.Header().Set("Content-Type", "application/json")
			e.HTTP.Unauthorized(w)
			return
		}

		s := new(model.Session)
		err := s.Decode(token)
		if err != nil {
			e.HTTP.Unauthorized(w)
			return
		}

		// Add session and token to context
		context.Set(r, "token", token)
		context.Set(r, "session", s)

		h.ServeHTTP(w, r)
	}
}
