package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/registry/pkg/registry/context"
	"github.com/lastbackend/registry/pkg/registry/http/handler"
	"net/http"
	"strconv"
	"time"
)

func NewRouter() *mux.Router {

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(Headers)

	// Session handlers
	r.HandleFunc("/template/{name}/{version}", Handler(handler.TemplateGetH)).Methods("GET")
	r.HandleFunc("/template", Handler(handler.TemplateListH)).Methods("GET")

	return r
}

func RunHttpServer(routes *mux.Router, port int) {

	var ctx = context.Get()

	ctx.Log.Infof("Listen server on %d port", port)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), routes); err != nil {
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
