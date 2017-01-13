package http

import (
	"fmt"
	c "github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/http/handler"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewRouter() *mux.Router {

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(headers)

	// Session handlers
	r.HandleFunc("/session", handle(handler.SessionCreateH)).Methods("POST")

	// User handlers
	r.HandleFunc("/user", handle(handler.UserCreateH)).Methods("POST")
	r.HandleFunc("/user", handle(handler.UserGetH, auth)).Methods("GET")

	// Build handlers
	r.HandleFunc("/build", handle(handler.BuildListH)).Methods("GET")
	r.HandleFunc("/build", handle(handler.BuildCreateH)).Methods("POST")

	// Project handlers
	r.HandleFunc("/project", handle(handler.ProjectListH, auth)).Methods("GET")
	r.HandleFunc("/project", handle(handler.ProjectCreateH, auth)).Methods("POST")
	r.HandleFunc("/project/{project}", handle(handler.ProjectInfoH, auth)).Methods("GET")
	r.HandleFunc("/project/{project}", handle(handler.ProjectUpdateH, auth)).Methods("PUT")
	r.HandleFunc("/project/{project}", handle(handler.ProjectRemoveH, auth)).Methods("DELETE")
	r.HandleFunc("/project/{project}/service", handle(handler.ServiceListH, auth)).Methods("GET")
	r.HandleFunc("/project/{project}/service/{service}", handle(handler.ServiceInfoH, auth)).Methods("GET")
	r.HandleFunc("/project/{project}/service/{service}", handle(handler.ServiceUpdateH, auth)).Methods("PUT")
	r.HandleFunc("/project/{project}/service/{service}", handle(handler.ServiceRemoveH, auth)).Methods("DELETE")
	r.HandleFunc("/project/{project}/service/{service}/logs", handle(handler.ServiceLogsH, auth)).Methods("GET")

	r.HandleFunc("/proxy", handle(handler.ProxyToken)).Methods("POST")

	// Deploy template/docker/source/repo
	r.HandleFunc("/deploy", handle(handler.DeployH, auth)).Methods("POST")

	// Template handlers
	r.HandleFunc("/template", handle(handler.TemplateListH)).Methods("GET")

	return r
}

func RunHttpServer(routes *mux.Router, port int) {

	var ctx = context.Get()

	ctx.Log.Infof("Listen server on %d port", port)

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
		c.Set(r, "token", token)
		c.Set(r, "session", s)

		h.ServeHTTP(w, r)
	}
}
