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
	v "github.com/lastbackend/lastbackend/pkg/api/views"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const logLevel = 2

func RouteListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Route: list")

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	items, err := rm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: find route list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Route().NewList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Route: write response err: %s", err)
		return
	}
}

func RouteInfoH(w http.ResponseWriter, r *http.Request) {

	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("Handler: Route: get route `%s`", rid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	item, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: find route by id `%s` err: %s", rid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Route: route `%s` not found", rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	response, err := v.V1().Route().New(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Route: write response err: %s", err)
		return
	}
}

func RouteCreateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Route: create route")

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	// request body struct
	opts := new(types.RouteOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Route: validation incoming data err: %s", err.Err())
		err.Http(w)
		return
	}

	// Check routes limit reachable
	if !ns.Quotas.Disabled && (ns.Resources.Routes+1) > ns.Quotas.Routes {
		log.V(logLevel).Warnf("Handler: Route: limit quotes reachable")
		errors.BadParameter("router").Http(w)
		return
	}

	ss, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get list services err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	rs, err := rm.Create(ns, ss, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: create route err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Route().New(rs).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Route: write response err: %s", err)
		return
	}
}

func RouteUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("Handler: Route: update route `%s`", nid)

	var (
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	// request body struct
	opts := new(types.RouteOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Route: validation incoming data err: %s", err.Err())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ss, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: check service exists by name `%s` err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ss == nil {
		log.Warnf("Handler: Route: service `%s` not found", ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: check route exists by name `%s` err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil && rs.Meta.Namespace != ns.Meta.Name {
		log.V(logLevel).Warnf("Handler: Route: route `%s` not found", rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	rs, err = rm.Update(rs, ns, ss, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: update route `%s` err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
	}

	response, err := v.V1().Route().New(rs).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Route: write response err: %s", err)
		return
	}
}

func RouteRemoveH(w http.ResponseWriter, r *http.Request) {

	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("Handler: Route: remove route %s", rid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		rm = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get route by id `%s` err: %s", rid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil && rs.Meta.Namespace != ns.Meta.Name {
		log.V(logLevel).Warnf("Handler: Route: route `%s` not found", rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	err = rm.SetState(rs, &types.RouteState{Provision: false, Destroy: true})
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: remove route `%s` err: %s", rid, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Route: write response err: %s", err)
		return
	}
}
