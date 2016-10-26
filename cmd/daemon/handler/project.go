package handler

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"net/http"
)

func ProjectListH(w http.ResponseWriter, _ *http.Request) {
	var ctx = context.Get()
	ctx.Log.Info("get projects list")

	w.WriteHeader(200)
	w.Write([]byte(ctx.Info.Version))
}
