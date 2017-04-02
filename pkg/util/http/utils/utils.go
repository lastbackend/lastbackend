package utils

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"net/http"
)

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func Session(r *http.Request) *types.Session {
	s, ok := context.GetOk(r, `session`)
	if !ok {
		return nil
	}
	return s.(*types.Session)
}
