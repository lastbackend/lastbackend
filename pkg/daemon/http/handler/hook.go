package handler

import (
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"net/http"
)

func HookExecuteH(w http.ResponseWriter, r *http.Request) {

	var (
		err       error
		hookModel *model.Hook
		ctx       = c.Get()
		params    = mux.Vars(r)
		hookParam = params["token"]
	)

	ctx.Log.Debug("Get project handler")

	hookModel, err = ctx.Storage.Hook().GetByToken(hookParam)
	if err != nil || hookModel == nil {
		ctx.Log.Error("Error: get hook by token", err.Error())
		e.HTTP.BadRequest(w)
		return
	}

	if hookModel.Service != "" {
		// TODO: Run redeploy
	} else if hookModel.Image != "" {
		// TODO: Run rebuild
	} else {
		e.HTTP.BadRequest(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte{})
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
