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

const logLevel = 2

func RouteListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Route: list")

	nid := utils.Vars(r)["namespace"]

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	items, err := rm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: find route list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Route().NewList(items).ToJson()
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

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("Handler: Route: get route `%s`", rid)

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

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

	response, err := v1.View().Route().New(item).ToJson()
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

	nid := utils.Vars(r)["namespace"]

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Route().CreateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Route: validation incoming data err: %s", err.Err())
		err.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	rs, err := rm.Create(ns, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: create route err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Route().New(rs).ToJson()
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
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Route().UpdateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Route: validation incoming data err: %s", err.Err())
		err.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: check route exists by selflink `%s` err: %s", ns.Meta.SelfLink, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil {
		log.V(logLevel).Warnf("Handler: Route: route `%s` not found", rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	rs, err = rm.Update(rs, ns, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: update route `%s` err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.View().Route().New(rs).ToJson()
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

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["route"]

	log.V(logLevel).Debugf("Handler: Route: remove route %s", rid)

	var (
		rm  = distribution.NewRouteModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Route().RemoveOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Route: validation incoming data err: %s", err.Err())
		err.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Route: get route by id `%s` err: %s", rid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil {
		log.V(logLevel).Warnf("Handler: Route: route `%s` not found", rid)
		errors.New("route").NotFound().Http(w)
		return
	}

	err = rm.SetStatus(rs, &types.RouteStatus{Stage:types.StageReady})
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
