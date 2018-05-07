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

package ingress


import (
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"strings"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:ingress"
)

func IngressInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress/{ingress} ingress ingressInfo
	//
	// Shows an ingress info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:info:> get ingress", logPrefix)

	var (
		im  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["ingress"]
	)

	n, err := im.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:info:> ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	response, err := v1.View().Ingress().New(n).ToJson()
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

func IngressGetSpecH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress/{ingress}/spec ingress ingressGetSpec
	//
	// Shows an ingress spec
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress spec response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_spec"
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:getspec:> list ingress", logPrefix)

	var (
		im  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		cid = utils.Vars(r)["cluster"]
		nid = utils.Vars(r)["ingress"]
		cache = envs.Get().GetCache().Ingress()
	)

	ing, err := im.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:getspec:> get ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ing == nil {
		log.V(logLevel).Warnf("%s:getspec:> ingress `%s` not found", logPrefix, cid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	spec := new(types.IngressSpec)
	spec = cache.Get(ing.Meta.Name)
	if spec == nil {
		sp, err := rm.ListSpec()
		if err != nil {
			log.V(logLevel).Warnf("%s:getspec:> ingress `%s` not found", logPrefix, cid)
			errors.HTTP.InternalServerError(w)
			return
		}
		spec = new(types.IngressSpec)
		spec.Routes = make(map[string]types.RouteSpec)
		for r, rsp := range sp {
			spec.Routes[r]=*rsp
		}
	}

	cache.Flush(ing.Meta.Name)

	response, err := v1.View().Ingress().NewSpec(spec).ToJson()
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

func IngressListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress ingress ingressList
	//
	// Shows an ingress list
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Ingress list response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get ingresss list", logPrefix)

	var (
		im = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
	)

	ingresss, err := im.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Ingress().NewList(ingresss).ToJson()
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

func IngressSetMetaH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/ingress/{ingress}/meta ingress ingressSetMeta
	//
	// Set ingress meta
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_ingress_meta"
	// responses:
	//   '200':
	//     description: Successfully set ingress meta
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["ingress"]

	log.V(logLevel).Debugf("%s:setmeta:> update ingress `%s`", logPrefix, nid)

	var (
		nm = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := new(request.IngressMetaOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update ingress err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:setmeta:> update ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	err = nm.SetMeta(n, opts.Meta)
	if err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update ingress `%s` err: %s", logPrefix, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Ingress().New(n).ToJson()
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

func IngressConnectH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/ingress/{ingress} ingress ingressConnect
	//
	// Connect ingress
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_ingress_connect"
	// responses:
	//   '200':
	//     description: Successfully connect ingress
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:connect:> ingress connect", logPrefix)

	var (
		im  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		iid = utils.Vars(r)["ingress"]
		cache = envs.Get().GetCache().Ingress()
	)

	// request body struct
	opts := new(request.IngressConnectOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:connect:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	ingress, err := im.Get(iid)
	if err != nil {
		log.V(logLevel).Errorf("%s:connect:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ingress == nil {

		nco := types.IngressCreateOptions{}

		nco.Meta.Name = iid
		nco.Status = opts.Status

		ingress, err = im.Create(&nco)
		if err != nil {
			log.V(logLevel).Errorf("%s:connect:> validation incoming data", logPrefix, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
	}

	if err := im.SetStatus(ingress, opts.Status); err != nil {
		log.V(logLevel).Errorf("%s:connect:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	cache.Clear(iid)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func IngressSetStatusH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/ingress/{ingress}/status ingress ingressSetStatus
	//
	// Set ingress status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_ingress_status"
	// responses:
	//   '200':
	//     description: Successfully set ingress status
	//   '400':
	//     description: Bad request
	//   '404':
	//     description: Ingress not found / Route not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:setstatus:> ingress set state", logPrefix)

	var (
		nm  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["ingress"]
	)

	// request body struct
	opts := new(request.IngressStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	ingress, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ingress == nil {
		log.V(logLevel).Warnf("%s:setstatus:> update ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	if err := nm.SetStatus(ingress, types.IngressStatus{
		Ready: opts.Ready,
	}); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	for p, s := range opts.Routes{
		keys := strings.Split(p, ":")
		if len(keys) != 2 {
			log.V(logLevel).Errorf("%s:setroutestatus:> invalid route selflink err: %s", logPrefix, p)
			errors.HTTP.BadRequest(w)
			return
		}

		route, err := rm.Get(keys[0], keys[1])
		if err != nil {
			log.V(logLevel).Errorf("%s:setroutestatus:> route not found selflink err: %s", logPrefix, p)
			errors.HTTP.InternalServerError(w)
			return
		}
		if route == nil {
			log.V(logLevel).Warnf("%s:setroutestatus:> update ingress `%s` not found", logPrefix, nid)
			errors.New("route").NotFound().Http(w)
			return
		}

		if err := rm.SetStatus(route, &types.RouteStatus{
			State:      s.State,
			Message:    s.Message,
		}); err != nil {
			log.V(logLevel).Errorf("%s:setroutestatus:> get ingresss list err: %s", logPrefix, err.Error())
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

func IngressSetRouteStatusH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /cluster/ingress/{ingress}/status/route/{pod} ingress ingressSetRouteStatus
	//
	// Set ingress route status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	//   - name: pod
	//     in: path
	//     description: pod id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_ingress_route_status"
	// responses:
	//   '200':
	//     description: Successfully set ingress route status
	//   '400':
	//     description: Bad request
	//   '404':
	//     description: Ingress not found / Route not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:setroutestatus:> ingress set route state", logPrefix)

	var (
		nm  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["ingress"]
		vid = utils.Vars(r)["route"]
	)

	// request body struct
	opts := new(request.IngressRouteStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:setroutestatus:> update ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
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
		log.V(logLevel).Warnf("%s:setroutestatus:> update ingress `%s` route not found %s", logPrefix, nid, vid)
		errors.New("route").NotFound().Http(w)
		return
	}

	if err := rm.SetStatus(route, &types.RouteStatus{
		State:   opts.State,
		Message: opts.Message,
	}); err != nil {
		log.V(logLevel).Errorf("%s:setroutestatus:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:setroutestatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func IngressRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /cluster/ingress/{ingress} ingress ingressRemove
	//
	// Remove ingress
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress removed
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:remove:>_ create ingress", logPrefix)

	var (
		nm  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["ingress"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if n == nil {
		log.V(logLevel).Warnf("%s:remove:>_ remove ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	if err := nm.Remove(n); err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:>_ write response err: %s", logPrefix, err.Error())
		return
	}
}