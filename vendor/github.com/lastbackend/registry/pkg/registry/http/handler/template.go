package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/registry/libs/errors"
	"github.com/lastbackend/registry/pkg/registry/context"
	"github.com/lastbackend/registry/pkg/template"
	"net/http"
)

func TemplateGetH(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		ctx     = context.Get()
		params  = mux.Vars(r)
		name    = params["name"]
		version = params["version"]
	)

	t := template.Get(name, version)
	if t == nil {
		e.Template.NotFound().Http(w)
		return
	}

	response, err := t.ToJson()
	if err != nil {
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func TemplateListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		ctx = context.Get()
	)

	t := template.List()

	buf, err := json.Marshal(t)
	if err != nil {
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(buf)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
