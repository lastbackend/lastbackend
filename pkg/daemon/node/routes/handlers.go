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
	"github.com/lastbackend/lastbackend/pkg/daemon/node"
	"github.com/lastbackend/lastbackend/pkg/daemon/node/routes/request"
	"github.com/lastbackend/lastbackend/pkg/daemon/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/daemon/service"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"net/http"
)

func NodeEventH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
	)

	log.Debug("Node event handler")

	// request body struct
	rq := new(request.RequestNodeEventS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	n := node.New()

	log.Debugf("try to find node by hostname: %s", rq.Meta.Hostname)
	item, err := n.Get(r.Context(), rq.Meta.Hostname)
	if err != nil {
		log.Error("Error: find node by hostname", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if item == nil {
		item, err = n.Create(r.Context(), &rq.Meta, &rq.State)
		if err != nil {
			log.Error("Error: can not create node", err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
	} else {
		item.Meta = rq.Meta
		item.State = rq.State
		n.SetMeta(r.Context(), item)
		n.SetState(r.Context(), item)
	}

	s := service.New(r.Context(), types.Meta{})
	if err := s.SetPods(r.Context(), rq.Pods); err != nil {
		log.Errorf("Error: set pods err %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewSpec(item).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}

}
