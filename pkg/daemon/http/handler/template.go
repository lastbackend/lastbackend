package handler

import (
	"net/http"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/template"
)

func TemplateListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er             error
		ctx            = c.Get()
		response_empty = func() {
			w.WriteHeader(http.StatusOK)
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
		ctx.Log.Error(err.Error())
		response_empty()
		return
	}

	if templates == nil {
		response_empty()
		return
	}

	response, err := templates.ToJson()
	if er != nil {
		ctx.Log.Error(err.Error())
		response_empty()
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
