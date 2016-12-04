package handler

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/template"
	"net/http"
)

func TemplateListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er             error
		ctx            = context.Get()
		response_empty = func() {
			w.WriteHeader(200)
			_, er = w.Write([]byte("[]"))
			if er != nil {
				ctx.Log.Error("Error: write response", er.Error())
				return
			}
			return
		}
	)

	templates, err := template.List()
	if err != nil {
		ctx.Log.Error(err.Err())
		response_empty()
		return
	}

	if templates == nil {
		response_empty()
		return
	}

	response, err := templates.ToJson()
	if er != nil {
		ctx.Log.Error(err.Err())
		response_empty()
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
