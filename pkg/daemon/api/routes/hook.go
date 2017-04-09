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

package routes

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

func HookExecuteH(w http.ResponseWriter, r *http.Request) {

	var (
		hookModel *types.Hook
		ctx       = context.Get()
		params    = utils.Vars(r)
		hookParam = params["token"]
	)

	ctx.Log.Debug("Get hook execute handler")

	hookModel, err := ctx.Storage.Hook().GetByToken(r.Context(), hookParam)
	if err != nil || hookModel == nil {
		ctx.Log.Error("Error: get hook by token", err.Error())
		errors.HTTP.BadRequest(w)
		return
	}

	if hookModel.Service.String() != "" {
		serviceModel, err := ctx.Storage.Service().GetByID(r.Context(), hookModel.Project, hookModel.Service)
		if err != nil && serviceModel == nil {
			ctx.Log.Error("Error: get service by name", err.Error())
			errors.HTTP.BadRequest(w)
			return
		}

		// TODO: REDEPLOY

	} else if hookModel.Image != "" {
		// TODO: Run rebuild
	} else {
		errors.HTTP.BadRequest(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte{})
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
