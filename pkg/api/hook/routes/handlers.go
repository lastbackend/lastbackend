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
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

func HookExecuteH(w http.ResponseWriter, r *http.Request) {

	var (
		log       = context.Get().GetLogger()
		storage   = context.Get().GetStorage()
		hookModel *types.Hook
		params    = utils.Vars(r)
		hookParam = params["token"]
	)

	log.Debug("Get hook execute handler")

	hookModel, err := storage.Hook().GetByToken(r.Context(), hookParam)
	if err != nil || hookModel == nil {
		log.Error("Error: get hook by token", err.Error())
		errors.HTTP.BadRequest(w)
		return
	}

	if hookModel.Service != "" {
		serviceModel, err := storage.Service().GetByName(r.Context(), hookModel.Project, hookModel.Service)
		if err != nil && serviceModel == nil {
			log.Error("Error: get service by name", err.Error())
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
	if _, err = w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}
