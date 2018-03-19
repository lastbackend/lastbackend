//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package node

import (
	"net/http"

	v "github.com/lastbackend/lastbackend/pkg/api/views"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

const logLevel = 2

func NodeGetH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: list node")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		cid = utils.Vars(r)["cluster"]
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("Handler: Node: node `%s` not found", cid)
		errors.New("node").NotFound().Http(w)
		return
	}

	response, err := v.V1().Node().New(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeGetSpecH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: list node")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		cid = utils.Vars(r)["cluster"]
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("Handler: Node: node `%s` not found", cid)
		errors.New("node").NotFound().Http(w)
		return
	}

	response, err := v.V1().Node().NewSpec(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: list node")

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	nodes, err := nm.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Node().NewList(nodes).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["node"]

	log.V(logLevel).Debugf("Handler: Node: update node `%s`", nid)

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := new(types.NodeUpdateOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.V(logLevel).Errorf("Handler: Node: validation incoming data", err)
	//	errors.New("Invalid incoming data").Unknown().Http(w)
	//	return
	//}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	err = nm.Update(n, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: update node `%s` err: %s", nid, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Node().New(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetInfoH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set info")

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	nodes, err := nm.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: get nodes list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Node().NewList(nodes).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetStateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set state")

	var response []byte

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetPodStateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set pod state")

	var response []byte

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetVolumeStateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set volume state")

	var response []byte

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeSetRouteStateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: node set route state")

	var response []byte

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}

func NodeRemoveH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Node: create node")

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Node: remove node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.Remove(n); err != nil {
		log.V(logLevel).Errorf("Handler: Node: remove node err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("Handler: Node: write response err: %s", err)
		return
	}
}