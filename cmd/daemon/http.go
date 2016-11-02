package daemon

import (
	"fmt"
	c "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"github.com/lastbackend/lastbackend/cmd/daemon/handler"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RunHttpServer(port int) {

	var ctx = context.Get()

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(Headers)

	r.HandleFunc("/version", Handler(handler.SystemVersionH)).Methods("GET")

	r.HandleFunc("/session", Handler(handler.SessionCreateH)).Methods("POST")

	// Account handlers
	r.HandleFunc("/account", Handler(handler.AccountCreateH)).Methods("POST")
	r.HandleFunc("/account", Handler(handler.AccountGetH, Auth)).Methods("GET")

	// User handlers
	r.HandleFunc("/user", Handler(handler.UserGetH, Auth)).Methods("GET")

	// Build handlers
	r.HandleFunc("/build", Handler(handler.BuildListH)).Methods("GET")
	r.HandleFunc("/build", Handler(handler.BuildCreateH)).Methods("POST")

	// Project handlers
	r.HandleFunc("/project", Handler(handler.ProjectListH, Auth)).Methods("GET")
	r.HandleFunc("/project", Handler(handler.ProjectCreateH, Auth)).Methods("POST")
	r.HandleFunc("/project/:id", Handler(handler.ProjectInfoH, Auth)).Methods("GET")
	r.HandleFunc("/project/:id", Handler(handler.ProjectDeleteH, Auth)).Methods("DELETE")
	r.HandleFunc("/project/:id", Handler(handler.ProjectDeleteH, Auth)).Methods("PUT")

	ctx.Log.Infof("Listen server on %d port", port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), r); err != nil {
		ctx.Log.Fatal("ListenAndServe: ", err)
	}
}

func Headers(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "X-CSRF-Token, Authorization, Content-Type, x-lastbackend, Origin, X-Requested-With, Content-Name, Accept")
	w.Header().Add("Content-Type", "application/json")
}

func Handler(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	headers := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			Headers(w, r)
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
func Auth(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var token string
		var params = mux.Vars(r)

		if _, ok := params["token"]; ok {
			token = params["token"]
		} else if r.Header.Get("Authorization") != "" {
			// Parse authorization header
			var auth = strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			// Check authorization header parts length and authorization header format
			if len(auth) != 2 && auth[0] != "Bearer" {
				e.HTTP.AccessDenied(w)
				return
			}

			token = auth[1]

		} else {
			w.Header().Set("Content-Type", "application/json")
			e.HTTP.AccessDenied(w)
			return
		}

		s := new(model.Session)
		err := s.Decode(token)
		if err != nil {
			e.HTTP.AccessDenied(w)
			return
		}

		// Add session and token to context
		c.Set(r, "token", token)
		c.Set(r, "session", s)

		h.ServeHTTP(w, r)
	}
}
