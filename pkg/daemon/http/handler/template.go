package handler

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"io/ioutil"
	"net/http"
)

func TemplateListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	var err = new(e.Http)

	_, resp, er := ctx.TemplateRegistry.GET("/template").Do()
	if er != nil {
		ctx.Log.Error(err.Message)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(buf)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
