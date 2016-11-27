package handler

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"io/ioutil"
	"net/http"
)

func TemplateListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er             error
		ctx            = context.Get()
		response_empty = func() {
			w.WriteHeader(404)
			_, er = w.Write([]byte("[]"))
			if er != nil {
				ctx.Log.Error("Error: write response", er.Error())
				return
			}
			return
		}
	)

	_, resp, er := ctx.TemplateRegistry.GET("/template").Do()
	if er != nil {
		ctx.Log.Error(er.Error())
		response_empty()
		return
	}

	buf, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		ctx.Log.Error(er.Error())
		response_empty()
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(buf)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
