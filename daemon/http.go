package daemon

import (
	"encoding/json"
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

			output := struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}{}

			output.Error.Code = e.Error()

			response, _ := json.Marshal(output)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(e.Status())
			w.Write(response)

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

	// app logic handler
	route.HandleFunc("/app", Handle(Handler{env, routes.CreateAppHandler})).Methods("PUT")
	route.HandleFunc("/app/{name}/logs", Handle(Handler{env, routes.LogsAppHandler})).Methods("GET")
	route.HandleFunc("/app/{name}/deploy", Handle(Handler{env, routes.DeployAppHandler})).Methods("POST")
	route.HandleFunc("/app/{name}/start", Handle(Handler{env, routes.StartAppHandler})).Methods("POST")
	route.HandleFunc("/app/{name}/stop", Handle(Handler{env, routes.StopAppHandler})).Methods("POST")
	route.HandleFunc("/app/{name}/restart", Handle(Handler{env, routes.RestartAppHandler})).Methods("POST")
	route.HandleFunc("/app/{name}", Handle(Handler{env, routes.RemoveAppHandler})).Methods("DELETE")

	// service logic handler
	route.HandleFunc("/service/{name}", Handle(Handler{env, routes.CreateServiceHandler})).Methods("PUT")
	route.HandleFunc("/service/{name}/logs", Handle(Handler{env, routes.LogsServiceHandler})).Methods("GET")
	route.HandleFunc("/service/{name}/start", Handle(Handler{env, routes.StartServiceHandler})).Methods("POST")
	route.HandleFunc("/service/{name}/stop", Handle(Handler{env, routes.StopServiceHandler})).Methods("POST")
	route.HandleFunc("/service/{name}/restart", Handle(Handler{env, routes.RestartServiceHandler})).Methods("POST")
	route.HandleFunc("/service/{name}", Handle(Handler{env, routes.RemoveServiceHandler})).Methods("DELETE")

	if err := http.ListenAndServe(":"+strconv.Itoa(env.Port), route); err != nil {
		env.Log.Fatal("ListenAndServe: ", err)
	}

	env.Log.Debugf("Listenning... on %v port", strconv.Itoa(env.Port))
}
