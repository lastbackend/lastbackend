package handler

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"net/http"
)

func InstanceListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	ctx.Log.Info("get projects list")

	w.WriteHeader(200)
	_, er = w.Write([]byte(ctx.Info.Version))
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
