//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/master/http/namespace/namespace"
	"github.com/lastbackend/lastbackend/internal/master/http/route/route"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"net/http"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/http/utils"
	"github.com/lastbackend/lastbackend/tools/log"
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
		rm = model.NewRouteModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
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
		rm = model.NewRouteModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
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
		mf = v1.Request().Route().Manifest()
	)

	// request body struct
	if e := mf.DecodeAndValidate(r.Body); e != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	rt, e := route.Create(r.Context(), ns, mf)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Route().New(rt).ToJson()
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
		mf = v1.Request().Route().Manifest()
	)

	// request body struct
	if e := mf.DecodeAndValidate(r.Body); e != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	rt, e := route.Fetch(r.Context(), ns.Meta.Name, rid)
	if e != nil {
		e.Http(w)
		return
	}

	rt, e = route.Update(r.Context(), ns, rt, mf)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Route().New(rt).ToJson()
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
		rm = model.NewRouteModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
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

	rs.Status.State = types.StateDestroy
	_, err = rm.Set(rs)
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
