package utils

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
