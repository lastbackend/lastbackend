package daemon

import (
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/daemon/routes"
	"github.com/deployithq/deployit/errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	*env.Env
	H func(env *env.Env, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case errors.Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		return err
	}

	return nil
}

func Handle(handlers ...Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")

		for _, handler := range handlers {
			if err := handler.ServeHTTP(w, r); err != nil {
				return
			}
		}
	}
}

func SetHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "X-CSRF-Token, Authorization, Content-Type, x-lastbackend, Origin, X-Requested-With, Content-Name, Accept")
	w.Header().Add("Content-Type", "application/json")
}

type Route struct {
}

func (r Route) Init(env *env.Env) {
	env.Log.Info("Init routes")

	route := mux.NewRouter()

	route.Methods("OPTIONS").HandlerFunc(SetHeaders)

	route.HandleFunc("/app/deploy", Handle(Handler{env, routes.DeployAppHandler})).Methods("POST")
	route.HandleFunc("/app/{id}/start", Handle(Handler{env, routes.StartAppHandler})).Methods("GET")
	route.HandleFunc("/app/{id}/stop", Handle(Handler{env, routes.StopAppHandler})).Methods("GET")
	route.HandleFunc("/app/{id}/restart", Handle(Handler{env, routes.RestartAppHandler})).Methods("GET")
	route.HandleFunc("/app/{id}", Handle(Handler{env, routes.RemoveAppHandler})).Methods("DELETE")

	if err := http.ListenAndServe(":"+strconv.Itoa(Port), route); err != nil {
		env.Log.Fatal("ListenAndServe: ", err)
	}

	env.Log.Debugf("Listenning... on %v port", strconv.Itoa(Port))
}
