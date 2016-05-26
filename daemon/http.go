package daemon

import (
	"github.com/deployithq/deployit/daemon/env"
	"net/http"
)

type Handler struct {
	*env.Env
	H func(e *env.Env, w http.ResponseWriter, r *http.Request) error
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
