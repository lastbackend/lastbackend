package daemon

import (
	"github.com/deployithq/deployit/daemon/env"
	"net/http"
	"github.com/deployithq/deployit/daemon/routes"
	"strconv"
	"github.com/gorilla/mux"
)

type Handler struct {
	*env.Env
	H func(env *env.Env, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch err.(type) {
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

type Route struct {
}

func (r Route) Init(env *env.Env) {
	env.Log.Info("Init routes")

	route := mux.NewRouter()

	route.HandleFunc("/app/deploy", Handle(Handler{env, routes.DeployAppHandler})).Methods("POST")

	if err := http.ListenAndServe(":"+strconv.Itoa(Port), route); err != nil {
		env.Log.Fatal("ListenAndServe: ", err)
	}

	env.Log.Debugf("Listenning... on %v port", strconv.Itoa(Port))
}