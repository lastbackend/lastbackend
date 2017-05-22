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
	"github.com/lastbackend/lastbackend/pkg/api/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/service"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

func HookExecuteH(w http.ResponseWriter, r *http.Request) {

	var (
		log       = context.Get().GetLogger()
		storage   = context.Get().GetStorage()
		params    = utils.Vars(r)
		hookParam = params["token"]
	)

	log.Debug("Get hook execute handler")

	hook, err := storage.Hook().Get(r.Context(), hookParam)
	if err != nil || hook == nil {
		log.Error("Error: get hook by token", err.Error())
		errors.HTTP.BadRequest(w)
		return
	}

	if hook.Service != "" {
		ns := namespace.New(r.Context())
		item, err := ns.Get(hook.Namespace)
		if err != nil {
			log.Error("Error: find namespace by name", err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
		if item == nil {
			errors.New("namespace").NotFound().Http(w)
			return
		}

		s := service.New(r.Context(), item.Meta)
		svc, err := s.Get(hook.Service)
		if err != nil {
			log.Error("Error: find service by name", err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
		if svc == nil {
			errors.New("service").NotFound().Http(w)
			return
		}

		if err := s.Redeploy(svc); err != nil {
			log.Error("Error: redeploy service", err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

	} else if hook.Image != "" {
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
