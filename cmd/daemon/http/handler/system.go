package handler

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"net/http"
)

func SystemVersionH(w http.ResponseWriter, _ *http.Request) {
	var ctx = context.Get()
	w.WriteHeader(200)
	w.Write([]byte(ctx.Info.Version))
}
