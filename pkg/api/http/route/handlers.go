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

package route

import (
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:route"
)

func RouteListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/route route routeList
	//
	// Shows a list of routes
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Route list response
	//     schema:
	//       "$ref": "#/definitions/views_route_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get routes list", logPrefix)

	nid := utils.Vars(r)["namespace"]

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	items, err := rm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> find route list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Route().NewList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func RouteInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/route/{route} route routeInfo
	//
	// Shows an info about route
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: route
	//     in: path
	//     description: route id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Route response
	//     schema:
	//       "$ref": "#/definitions/views_route"
	//   '404':
	//     description: Namespace not found / Route not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("%s:info:> get route `%s`", logPrefix, rid)

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:info:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	item, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> find route by id `%s` err: %s", rid, logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("%s:info:> route `%s` not found", logPrefix, rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	response, err := v1.View().Route().New(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func RouteCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/route route routeCreate
	//
	// Creates a route
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_route_create"
	// responses:
	//   '200':
	//     description: Route was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_route"
	//   '400':
	//     description: Bad rules parameter
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:create:> create route", logPrefix)

	nid := utils.Vars(r)["namespace"]

	var (
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Route().CreateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:create:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.List(ns.Meta.SelfLink)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> get services", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	rcopts := types.RouteCreateOptions{
		Name:     opts.Name,
		Domain:   ns.Meta.Endpoint,
		Security: opts.Security,
		Rules:    make([]types.RuleOption, 0),
	}

	var links = make(map[string]string)

	for _, s := range svc {
		links[s.Meta.Name] = s.Meta.SelfLink
	}

	for _, r := range opts.Rules {

		if r.Service == "" || r.Port == 0 {
			continue
		}

		if r.Path == "" {
			r.Path = "/"
		}

		if _, ok := links[r.Service]; ok {
			rcopts.Rules = append(rcopts.Rules, types.RuleOption{
				Service:  r.Service,
				Endpoint: svc[links[r.Service]].Meta.Endpoint,
				Path:     r.Path,
				Port:     r.Port,
			})
		}
	}

	if len(rcopts.Rules) == 0 {
		err := errors.New("route rules are incorrect")
		log.V(logLevel).Errorf("%s:create:> route rules empty", logPrefix, err.Error())
		errors.New("route").BadParameter("rules", err).Http(w)
		return
	}

	rs, err := rm.Create(ns, &rcopts)

	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create route err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Route().New(rs).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func RouteUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/route/{route} route routeUpdate
	//
	// Update route
	//
	// ---
	// deprecated: true
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: route
	//     in: path
	//     description: route id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_route_update"
	// responses:
	//   '200':
	//     description: Route was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_route"
	//   '400':
	//     description: Bad rules parameter
	//   '404':
	//     description: Namespace not found / Route not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("%s:update:> update route `%s`", logPrefix, nid)

	var (
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Route().UpdateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:update:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> check route exists by selflink `%s` err: %s", logPrefix, ns.Meta.SelfLink, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil {
		log.V(logLevel).Warnf("%s:update:> route `%s` not found", logPrefix, rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	svc, err := sm.List(ns.Meta.SelfLink)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> get services", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	ruopts := types.RouteUpdateOptions{
		Security: opts.Security,
		Rules:    make([]types.RuleOption, 0),
	}

	var links = make(map[string]string)

	for _, s := range svc {
		links[s.Meta.Name] = s.Meta.SelfLink
	}

	for _, r := range opts.Rules {

		if r.Service == "" || r.Port == 0 {
			continue
		}

		if r.Path == "" {
			r.Path = "/"
		}

		if _, ok := links[r.Service]; ok {
			ruopts.Rules = append(ruopts.Rules, types.RuleOption{
				Service:  r.Service,
				Endpoint: svc[links[r.Service]].Meta.Endpoint,
				Path:     r.Path,
				Port:     r.Port,
			})
		}
	}

	if len(ruopts.Rules) == 0 {
		err := errors.New("route rules are incorrect")
		log.V(logLevel).Errorf("%s:create:> route rules empty", logPrefix, err.Error())
		errors.New("route").BadParameter("rules", err).Http(w)
		return
	}

	rs, err = rm.Update(rs, &ruopts)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> update route `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.View().Route().New(rs).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func RouteRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace}/route/{route} route routeRemove
	//
	// Removes route
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: route
	//     in: path
	//     description: route id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Route was successfully removed
	//   '404':
	//     description: Namespace not found / Route not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("%s:remove:> remove route %s", logPrefix, rid)

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:remove:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get route by id `%s` err: %s", logPrefix, rid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil {
		log.V(logLevel).Warnf("%s:remove:> route `%s` not found", logPrefix, rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	err = rm.SetStatus(rs, &types.RouteStatus{State: types.StateDestroy})
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove route `%s` err: %s", logPrefix, rid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
