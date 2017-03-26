package utils

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetVars(r *http.Request) map[string]string  {
	return mux.Vars(r)
}
