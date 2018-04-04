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

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"strings"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:node"
)

func NodeInfoH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:info:> get node", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get node err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:info:> node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	response, err := v1.View().Node().New(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeGetSpecH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:getspec:> list node", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		cid = utils.Vars(r)["cluster"]
		nid = utils.Vars(r)["node"]
		cache = envs.Get().GetCache().Node()
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:getspec:> get node err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:getspec:> node `%s` not found", logPrefix, cid)
		errors.New("node").NotFound().Http(w)
		return
	}

	spec := new(types.NodeSpec)
	spec = cache.Get(n.Meta.Name)
	if spec == nil {
		spec, err = nm.GetSpec(n)
		if err != nil {
			log.V(logLevel).Warnf("%s:getspec:> node `%s` not found", logPrefix, cid)
			errors.HTTP.InternalServerError(w)
			return
		}
	}
	cache.Flush(n.Meta.Name)

	response, err := v1.View().Node().NewSpec(spec).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:getspec:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:getspec:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:list:> get nodes list", logPrefix)

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	nodes, err := nm.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Node().NewList(nodes).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeSetMetaH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["node"]

	log.V(logLevel).Debugf("%s:setmeta:> update node `%s`", logPrefix, nid)

	var (
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := new(request.NodeMetaOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update node err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:setmeta:> update node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	err = nm.SetMeta(n, opts.Meta)
	if err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update node `%s` err: %s", logPrefix, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Node().New(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeConnectH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:connect:> node connect", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	// request body struct
	opts := new(request.NodeConnectOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:connect:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	node, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:connect:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if node == nil {

		nco := types.NodeCreateOptions{}

		nco.Meta.Name = opts.Info.Hostname
		nco.Info = opts.Info
		nco.Status = opts.Status
		nco.Network = opts.Network


		node, err = nm.Create(&nco)
		if err != nil {
			log.V(logLevel).Errorf("%s:connect:> validation incoming data", logPrefix, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

		if err := nm.SetOnline(node); err != nil {
			log.V(logLevel).Errorf("%s:connect:> get nodes list err: %s", logPrefix, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte{}); err != nil {
			log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
			return
		}

		return
	}

	if err := nm.SetInfo(node, opts.Info); err != nil {
		log.V(logLevel).Errorf("%s:connect:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.SetStatus(node, opts.Status); err != nil {
		log.V(logLevel).Errorf("%s:connect:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.SetNetwork(node, opts.Network); err != nil {
		log.V(logLevel).Errorf("%s:connect:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.SetOnline(node); err != nil {
		log.V(logLevel).Errorf("%s:connect:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeSetStatusH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:setstatus:> node set state", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		pm  = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	// request body struct
	opts := new(request.NodeStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	node, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if node == nil {
		log.V(logLevel).Warnf("%s:setstatus:> update node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	if err := nm.SetStatus(node, types.NodeStatus{
		Capacity:  opts.Resources.Capacity,
		Allocated: opts.Resources.Allocated,
	}); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := nm.SetOnline(node); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	for p, s := range opts.Pods{
		keys := strings.Split(p, ":")
		if len(keys) != 4 {
			log.V(logLevel).Errorf("%s:setpodstatus:> invalid pod selflink err: %s", logPrefix, p)
			errors.HTTP.BadRequest(w)
			return
		}

		pod, err := pm.Get(keys[0], keys[1], keys[2], keys[3])
		if err != nil {
			log.V(logLevel).Errorf("%s:setpodstatus:> pod not found selflink err: %s", logPrefix, p)
			errors.HTTP.InternalServerError(w)
			return
		}
		if pod == nil {
			log.V(logLevel).Warnf("%s:setpodstatus:> update node `%s` not found", logPrefix, nid)
			errors.New("pod").NotFound().Http(w)
			return
		}

		if err := pm.SetStatus(pod, &types.PodStatus{
			Stage:      s.State,
			Message:    s.Message,
			Steps:      s.Steps,
			Network:    s.Network,
			Containers: s.Containers,
		}); err != nil {
			log.V(logLevel).Errorf("%s:setpodstatus:> get nodes list err: %s", logPrefix, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:setstatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeSetPodStatusH(w http.ResponseWriter, r *http.Request) {

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		pm  = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
		pid = utils.Vars(r)["pod"]
	)

	// request body struct
	opts := new(request.NodePodStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setpodstatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setpodstatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:setpodstatus:> update node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	keys := strings.Split(pid, ":")
	if len(keys) != 4 {
		log.V(logLevel).Errorf("%s:setpodstatus:> invalid pod selflink err: %s", logPrefix, pid)
		errors.HTTP.BadRequest(w)
		return
	}

	pod, err := pm.Get(keys[0], keys[1], keys[2], keys[3])
	if err != nil {
		log.V(logLevel).Errorf("%s:setpodstatus:> pod not found selflink err: %s", logPrefix, pid)
		errors.HTTP.InternalServerError(w)
		return
	}
	if pod == nil {
		log.V(logLevel).Warnf("%s:setpodstatus:> update node `%s` not found", logPrefix, nid)
		errors.New("pod").NotFound().Http(w)
		return
	}

	if err := pm.SetStatus(pod, &types.PodStatus{
		Stage:      opts.State,
		Message:    opts.Message,
		Steps:      opts.Steps,
		Network:    opts.Network,
		Containers: opts.Containers,
	}); err != nil {
		log.V(logLevel).Errorf("%s:setpodstatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:setpodstatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeSetVolumeStatusH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:setvolumestatus:> node set volume state", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		vm  = distribution.NewVolumeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
		vid = utils.Vars(r)["volume"]
	)

	// request body struct
	opts := new(request.NodeVolumeStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setvolumestatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setvolumestatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:setvolumestatus:> update node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	keys := strings.Split(vid, ":")
	if len(keys) != 2 {
		log.V(logLevel).Errorf("%s:setvolumestatus:> invalid volume selflink err: %s", logPrefix, vid)
		errors.HTTP.BadRequest(w)
		return
	}

	volume, err := vm.Get(keys[0], keys[1])
	if err != nil {
		log.V(logLevel).Errorf("%s:setvolumestatus:> pod not found selflink err: %s", logPrefix, vid)
		errors.HTTP.NotFound(w)
		return
	}
	if volume == nil {
		log.V(logLevel).Warnf("%s:setvolumestatus:> update node `%s` volume not found %s", logPrefix, nid, vid)
		errors.New("volume").NotFound().Http(w)
		return
	}

	if err := vm.SetStatus(volume, &types.VolumeStatus{
		State:   opts.State,
		Message: opts.Message,
	}); err != nil {
		log.V(logLevel).Errorf("%s:setvolumestatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:setvolumestatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeSetRouteStatusH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:setroutestatus:> node set route state", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
		vid = utils.Vars(r)["route"]
	)

	// request body struct
	opts := new(request.NodeRouteStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:setroutestatus:> update node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	keys := strings.Split(vid, ":")
	if len(keys) != 2 {
		log.V(logLevel).Errorf("%s:setroutestatus:> invalid route selflink err: %s", logPrefix, vid)
		errors.HTTP.BadRequest(w)
		return
	}

	route, err := rm.Get(keys[0], keys[1])
	if err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> pod not found selflink err: %s", logPrefix, vid)
		errors.HTTP.NotFound(w)
		return
	}
	if route == nil {
		log.V(logLevel).Warnf("%s:setroutestatus:> update node `%s` route not found %s", logPrefix, nid, vid)
		errors.New("route").NotFound().Http(w)
		return
	}

	if err := rm.SetStatus(route, &types.RouteStatus{
		State:   opts.State,
		Message: opts.Message,
	}); err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> get nodes list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:setroutestatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NodeRemoveH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:remove:>_ create node", logPrefix)

	var (
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["node"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove node err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if n == nil {
		log.V(logLevel).Warnf("%s:remove:>_ remove node `%s` not found", logPrefix, nid)
		errors.New("node").NotFound().Http(w)
		return
	}

	if err := nm.Remove(n); err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove node err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:>_ write response err: %s", logPrefix, err.Error())
		return
	}
}
