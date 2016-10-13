package daemon

import (
	"fmt"
	"github.com/deployithq/deployit/cmd/daemon/context"
	"github.com/deployithq/deployit/cmd/daemon/handler"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func RunHttpServer(port int) {

	var ctx = context.Get()

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(Headers)

	r.HandleFunc("/version", Handler(handler.SystemVersionH)).Methods("GET")

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
