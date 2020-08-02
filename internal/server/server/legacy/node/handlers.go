//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/server/server/legacy/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:node"
)

// Handler represent the http handler for node
type Handler struct {
}

// NewNodeHandler will initialize the node resources endpoint
func NewNodeHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init node routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/cluster/node", h.Handle(mw.Authenticate(handler.NodeListH))).Methods(http.MethodGet)
	r.Handle("/cluster/node/{node}", h.Handle(mw.Authenticate(handler.NodeInfoH))).Methods(http.MethodGet)
	r.Handle("/cluster/node/{node}/spec", h.Handle(mw.Authenticate(handler.NodeGetSpecH))).Methods(http.MethodGet)
	r.Handle("/cluster/node/{node}", h.Handle(mw.Authenticate(handler.NodeRemoveH))).Methods(http.MethodDelete)
	r.Handle("/cluster/node/{node}", h.Handle(mw.Authenticate(handler.NodeConnectH))).Methods(http.MethodPut)
	r.Handle("/cluster/node/{node}/meta", h.Handle(mw.Authenticate(handler.NodeSetMetaH))).Methods(http.MethodPut)
	r.Handle("/cluster/node/{node}/status", h.Handle(mw.Authenticate(handler.NodeSetStatusH))).Methods(http.MethodPut)
}

func (handler Handler) NodeInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/node/{node} node nodeInfo
	//
	// Shows an info about node
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: node
	//     in: path
	//     description: node id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Node response
	//     schema:
	//       "$ref": "#/definitions/views_node"
	//   '404':
	//     description: Node not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:info:> get node", logPrefix)
	//
	//var (
	//	nm  = model.NewNodeModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["node"]
	//)
	//
	//n, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:info:> get node err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if n == nil {
	//	log.Warnf("%s:info:> node `%s` not found", logPrefix, nid)
	//	errors.New("node").NotFound().Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Node().New(n).ToJson()
	//if err != nil {
	//	log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NodeGetSpecH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/node/{node}/spec node nodeGetSpec
	//
	// Shows an info about node spec
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: node
	//     in: path
	//     description: node id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Node spec response
	//     schema:
	//       "$ref": "#/definitions/views_node_spec"
	//   '404':
	//     description: Node not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:getspec:> list node", logPrefix)
	//
	//var (
	//	stg = envs.Get().GetStorage()
	//	nm  = model.NewNodeModel(r.Context(), stg)
	//
	//	cid = util.Vars(r)["cluster"]
	//	nid = util.Vars(r)["node"]
	//)
	//
	//n, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:getspec:> get node err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if n == nil {
	//	log.Warnf("%s:getspec:> node `%s` not found", logPrefix, cid)
	//	errors.New("node").NotFound().Http(w)
	//	return
	//}
	//
	//spec, err := getNodeSpec(r.Context(), n)
	//if err != nil {
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Node().NewManifest(spec).ToJson()
	//if err != nil {
	//	log.Errorf("%s:getspec:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:getspec:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NodeListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/node node nodeList
	//
	// Shows a list of nodes
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Node list response
	//     schema:
	//       "$ref": "#/definitions/views_node_list"
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:list:> get nodes list", logPrefix)
	//
	//var (
	//	nm = model.NewNodeModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//nodes, err := nm.List()
	//if err != nil {
	//	log.Errorf("%s:list:> get nodes list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Node().NewList(nodes).ToJson()
	//if err != nil {
	//	log.Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NodeSetMetaH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/node/{node}/meta node nodeSetMeta
	//
	// Set node meta
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: node
	//     in: path
	//     description: node id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_node_meta"
	// responses:
	//   '200':
	//     description: Successfully set node meta
	//     schema:
	//       "$ref": "#/definitions/views_node"
	//   '404':
	//     description: Node not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//nid := util.Vars(r)["node"]
	//
	//log.Debugf("%s:setmeta:> update node `%s`", logPrefix, nid)
	//
	//var (
	//	nm = model.NewNodeModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//// request body struct
	//opts := new(request.NodeMetaOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:setmeta:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//n, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:setmeta:> update node err: %s", err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if n == nil {
	//	log.Warnf("%s:setmeta:> update node `%s` not found", logPrefix, nid)
	//	errors.New("node").NotFound().Http(w)
	//	return
	//}
	//
	//n.Meta.Set(opts.Meta)
	//
	//err = nm.Set(n)
	//if err != nil {
	//	log.Errorf("%s:setmeta:> update node `%s` err: %s", logPrefix, nid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Node().New(n).ToJson()
	//if err != nil {
	//	log.Errorf("%s:setmeta:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:setmeta:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NodeConnectH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/node/{node} node nodeConnect
	//
	// Connect node
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: node
	//     in: path
	//     description: node id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_node_connect"
	// responses:
	//   '200':
	//     description: Node connected
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:connect:> node connect", logPrefix)
	//
	//var (
	//	stg   = envs.Get().GetStorage()
	//	nm    = model.NewNodeModel(r.Context(), stg)
	//	sn    = model.NewNetworkModel(r.Context(), stg)
	//	nid   = util.Vars(r)["node"]
	//	cache = envs.Get().GetCache().Node()
	//)
	//
	//// request body struct
	//opts := new(request.NodeConnectOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:connect:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//node, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:connect:> get node err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//snet, err := sn.SubnetGet(opts.Network.SubnetSpec.CIDR)
	//if err != nil {
	//	log.Errorf("%s:connect:> get subnet err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if node == nil {
	//
	//	nco := types.NodeCreateOptions{}
	//	nco.Meta.Name = opts.Info.Hostname
	//
	//	nco.Info = opts.Info
	//	nco.Status = opts.Status
	//	nco.Meta.Subnet = types.SubnetGetNameFromCIDR(opts.Network.CIDR)
	//	if snet != nil {
	//		nco.Status.State.CNI.State = types.StateWarning
	//		nco.Status.State.CNI.Message = errors.ErrEntityExists
	//	}
	//
	//	nco.Security.TLS = opts.TLS
	//
	//	if opts.SSL != nil {
	//		nco.Security.SSL = new(types.NodeSSL)
	//		nco.Security.SSL.CA = opts.SSL.CA
	//		nco.Security.SSL.Cert = opts.SSL.Cert
	//		nco.Security.SSL.Key = opts.SSL.Key
	//	}
	//
	//	nco.Status.Capacity = opts.Status.Capacity
	//
	//	node, err = nm.Put(&nco)
	//	if err != nil {
	//		log.Errorf("%s:connect:> node put error: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//
	//	if snet == nil {
	//		if _, err := sn.SubnetPut(node.SelfLink().String(), opts.Network.SubnetSpec); err != nil {
	//			log.Errorf("%s:connect:> snet put error: %s", logPrefix, err.Error())
	//			errors.HTTP.InternalServerError(w)
	//			return
	//		}
	//	}
	//
	//	w.WriteHeader(http.StatusOK)
	//	if _, err := w.Write([]byte{}); err != nil {
	//		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
	//		return
	//	}
	//
	//	return
	//}
	//
	//ou := new(types.NodeUpdateMetaOptions)
	//ou.NodeUpdateInfoOptions.Set(opts.Info)
	//node.Meta.Set(ou)
	//node.Status.Capacity = opts.Status.Capacity
	//node.Spec.Security.TLS = opts.TLS
	//
	//if opts.SSL != nil {
	//	node.Spec.Security.SSL = new(types.NodeSSL)
	//	node.Spec.Security.SSL.CA = opts.SSL.CA
	//	node.Spec.Security.SSL.Cert = opts.SSL.Cert
	//	node.Spec.Security.SSL.Key = opts.SSL.Key
	//}
	//
	//if err := nm.Set(node); err != nil {
	//	log.Errorf("%s:connect:> get node set err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if snet == nil {
	//	if _, err := sn.SubnetPut(node.SelfLink().String(), opts.Network.SubnetSpec); err != nil {
	//		log.Errorf("%s:connect:> snet put error: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//	}
	//} else {
	//	if !sn.SubnetEqual(snet, opts.Network.SubnetSpec) {
	//		snet.Spec = opts.Network.SubnetSpec
	//		if err := sn.SubnetSet(snet); err != nil {
	//			log.Errorf("%s:connect:> get subnet set err: %s", logPrefix, err.Error())
	//			errors.HTTP.InternalServerError(w)
	//			return
	//		}
	//	}
	//}
	//
	//cache.Clear(nid)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NodeSetStatusH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/node/{node}/status node nodeSetStatus
	//
	// Set node status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: node
	//     in: path
	//     description: node id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_node_status"
	// responses:
	//   '200':
	//     description: Successfully set node status
	//   '400':
	//     description: Bad request
	//   '404':
	//     description: Node not found / Pod not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:setstatus:> node set state", logPrefix)
	//
	//var (
	//	nm = model.NewNodeModel(r.Context(), envs.Get().GetStorage())
	//	pm = model.NewPodModel(r.Context(), envs.Get().GetStorage())
	//	vm = model.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	//
	//	nid = util.Vars(r)["node"]
	//)
	//
	//// request body struct
	//opts := new(request.NodeStatusOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//node, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:setstatus:> get nodes list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if node == nil {
	//	log.Warnf("%s:setstatus:> update node `%s` not found", logPrefix, nid)
	//	errors.New("node").NotFound().Http(w)
	//	return
	//}
	//
	//node.Status.State = opts.State
	//node.Status.Online = true
	//node.Status.Capacity = opts.Resources.Capacity
	//
	//if err := nm.Set(node); err != nil {
	//	log.Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//for p, s := range opts.Pods {
	//
	//	log.Debugf("set pod status: %s", p)
	//
	//	keys := strings.Split(p, ":")
	//	if len(keys) != 4 {
	//		log.Errorf("%s:setpodstatus:> invalid pod selflink err: %s", logPrefix, p)
	//		errors.HTTP.BadRequest(w)
	//		return
	//	}
	//
	//	sl := types.PodSelfLink{}
	//	if err := sl.Parse(p); err != nil {
	//		continue
	//	}
	//
	//	pod, err := pm.Get(sl.String())
	//	if err != nil {
	//		log.Errorf("%s:setpodstatus:> pod not found selflink err: %s", logPrefix, p)
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//	if pod == nil {
	//		log.Warnf("%s:setpodstatus:>pod not found `%s` not found", logPrefix, p)
	//		if err := pm.ManifestDel(nid, p); err != nil {
	//			if !errors.Storage().IsErrEntityNotFound(err) {
	//				log.Warnf("%s:setpodstatus:>pod manifest del err `%s` ", logPrefix, err.Error())
	//				continue
	//			}
	//		}
	//		continue
	//	}
	//
	//	pod.Status.State = s.State
	//	pod.Status.Status = s.Status
	//	pod.Status.Running = s.Running
	//	pod.Status.Message = s.Message
	//	pod.Status.Runtime = s.Runtime
	//	pod.Status.Network = s.Network
	//	pod.Status.Steps = s.Steps
	//
	//	if err := pm.Update(pod); err != nil {
	//		log.Errorf("%s:setpodstatus:> update pod err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//}
	//
	//for v, s := range opts.Volumes {
	//	log.Debugf("set volume status: %s", v)
	//
	//	keys := strings.Split(v, ":")
	//	if len(keys) != 2 {
	//		log.Errorf("%s:set volume status:> invalid volume selflink err: %s", logPrefix, v)
	//		errors.HTTP.BadRequest(w)
	//		return
	//	}
	//
	//	volume, err := vm.Get(keys[0], keys[1])
	//	if err != nil {
	//		log.Errorf("%s:set volume status:> volume not found by selflink err: %s", logPrefix, v)
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//	if volume == nil {
	//		log.Warnf("%s:set volume status:>volume not found `%s` not found", logPrefix, v)
	//		if err := vm.ManifestDel(nid, v); err != nil {
	//			if !errors.Storage().IsErrEntityNotFound(err) {
	//				log.Warnf("%s:set volume status:>volume manifest del err `%s` ", logPrefix, err.Error())
	//				continue
	//			}
	//		}
	//		continue
	//	}
	//
	//	volume.Status.State = s.State
	//	volume.Status.Message = s.Message
	//
	//	if err := vm.Update(volume); err != nil {
	//		log.Errorf("%s:set volume status:> update pod err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//}
	//
	//spec, err := getNodeSpec(r.Context(), node)
	//if err != nil {
	//	errors.HTTP.InternalServerError(w)
	//}
	//
	//response, err := v1.View().Node().NewManifest(spec).ToJson()
	//if err != nil {
	//	log.Errorf("%s:getspec:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:setstatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NodeRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /cluster/node/{node} node nodeRemove
	//
	// Remove node
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: node
	//     in: path
	//     description: node id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Node removed
	//   '404':
	//     description: Node not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:remove:> remove node", logPrefix)
	//
	//var (
	//	stg = envs.Get().GetStorage()
	//	nm  = model.NewNodeModel(r.Context(), stg)
	//	sm  = model.NewNetworkModel(r.Context(), stg)
	//	nid = util.Vars(r)["node"]
	//)
	//
	//n, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:remove:>_ remove node err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if n == nil {
	//	log.Warnf("%s:remove:>_ remove node `%s` not found", logPrefix, nid)
	//	errors.New("node").NotFound().Http(w)
	//	return
	//}
	//
	//if err := nm.Remove(n); err != nil {
	//	log.Errorf("%s:remove:> remove node err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if err := sm.SubnetDel(n.Meta.Subnet); err != nil {
	//	log.Errorf("%s:remove:> remove subnet err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:>_ write response err: %s", logPrefix, err.Error())
		return
	}
}

//
//func (handler Handler) getNodeSpec(ctx context.Context, n *types.Node) (*types.NodeManifest, error) {
//
//	var (
//		cache = envs.Get().GetCache().Node()
//		spec  = cache.Get(n.Meta.Name)
//		stg   = envs.Get().GetStorage()
//		pm    = model.NewPodModel(ctx, stg)
//		vm    = model.NewVolumeModel(ctx, stg)
//		em    = model.NewEndpointModel(ctx, stg)
//		ns    = model.NewNetworkModel(ctx, stg)
//	)
//
//	if spec == nil {
//
//		spec = new(types.NodeManifest)
//		spec.Meta.Initial = true
//		spec.Resolvers = cache.GetResolvers()
//		spec.Exporter = cache.GetExporterEndpoint()
//		spec.Configs = cache.GetConfigs()
//
//		pods, err := pm.ManifestMap(n.Meta.Name)
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get pod manifests for node err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//		spec.Pods = pods.Items
//
//		volumes, err := vm.ManifestMap(n.Meta.Name)
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get volume manifests for node err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//		spec.Volumes = volumes.Items
//
//		endpoints, err := em.ManifestMap()
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get endpoint manifests for node err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//		spec.Endpoints = endpoints.Items
//
//		subnets, err := ns.SubnetManifestMap()
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get endpoint manifests for node err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//
//		spec.Network = subnets.Items
//	}
//	cache.Flush(n.Meta.Name)
//
//	return spec, nil
//
//}
