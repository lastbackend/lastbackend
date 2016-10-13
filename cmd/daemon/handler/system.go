package handler

import (
	"github.com/deployithq/deployit/cmd/daemon/context"
	"net/http"
)

func SystemVersionH(w http.ResponseWriter, _ *http.Request) {
	var ctx = context.Get()
	w.WriteHeader(200)
	w.Write([]byte(ctx.Version))
}
