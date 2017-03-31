package utils

import (
	"net/http"
	"github.com/gorilla/mux"
)

func GetVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
