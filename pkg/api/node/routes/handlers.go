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
	"github.com/lastbackend/lastbackend/pkg/api/node"
	"github.com/lastbackend/lastbackend/pkg/api/node/routes/request"
	"github.com/lastbackend/lastbackend/pkg/api/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/api/pod"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 2

func NodeEventH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
	)

	log.V(logLevel).Debug("Handler: Node: event handling")

	// request body struct
	rq := new(request.RequestNodeEventS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Node: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	p := pod.New(r.Context())
	for _, item := range rq.Pods {
		if err := p.Set(item); err != nil {
			log.V(logLevel).Errorf("Handler: Node: set pods err: %s", err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
	}

	log.V(logLevel).Debugf("Handler: Node: try to find node by id: %s", rq.Meta.ID)

	n := node.New(r.Context())
	item, err := n.Get(rq.Meta.ID)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node by id  `%s` err: %s", rq.Meta.ID, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if item == nil {
		log.V(logLevel).Debug("Handler: Node: node not found, create a new one")

		item, err = n.Create(&rq.Meta, &rq.State)
		if err != nil {
			log.V(logLevel).Errorf("Handler: Node: create node err: %s", err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
	} else {
		log.V(logLevel).Debug("Handler: Node: update node")

		item.Meta = rq.Meta
		item.State = rq.State
		if err := n.Update(item); err != nil {
			log.V(logLevel).Errorf("Handler: Node: update node err: %s", err.Error())
			return
		}
	}

	response, err := v1.NewSpec(item).ToJson()
	if err != nil {
		log.Error("Handler: Node: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Handler: Node: write response err: %s", err.Error())
		return
	}
}

func NodeListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
	)

	log.V(logLevel).Debug("Handler: Node: list node")

	n := node.New(r.Context())
	nodes, err := n.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewNodeList(nodes).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Handler: Node: write response err: %s", err.Error())
		return
	}
}
