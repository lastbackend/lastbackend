package handler

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"net/http"
	"time"
)

func BuildListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	ctx.Log.Info("get builds list")
	builds, err := ctx.Storage.Build().GetByImage("", "")
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(builds)
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

func BuildCreateH(w http.ResponseWriter, _ *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	ctx.Log.Info("create build")

	b := new(model.Build)
	b.Created = time.Now()
	b.Updated = time.Now()

	build, err := ctx.Storage.Build().Insert(b)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(build)
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
