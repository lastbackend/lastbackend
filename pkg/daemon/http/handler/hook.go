//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package handler

import (
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service"
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
		serviceModel, err := ctx.Storage.Service().GetByID(hookModel.User, hookModel.Service)
		if err != nil && serviceModel == nil {
			ctx.Log.Error("Error: get service by id", err.Error())
			e.HTTP.BadRequest(w)
			return
		}

		projectModel, err := ctx.Storage.Project().GetByID(serviceModel.User, serviceModel.Project)
		if err != nil && serviceModel == nil {
			ctx.Log.Error("Error: get project by id", err.Error())
			e.HTTP.BadRequest(w)
			return
		}

		err = service.UpdateImage(ctx.K8S, projectModel.ID, "lb-"+serviceModel.Name)
		if err != nil && serviceModel == nil {
			ctx.Log.Error("Error: update image for service", err.Error())
			e.HTTP.BadRequest(w)
			return
		}

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
