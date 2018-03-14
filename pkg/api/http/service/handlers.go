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

package service

import (
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	v "github.com/lastbackend/lastbackend/pkg/api/views"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
)

const logLevel = 2

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("Handler: Service: list services in %s", nid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		dm = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	items, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service list in namespace `%s` err: %s", ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	dl, err := dm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get pod list by service id `%s` err: %s", ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Service().NewList(items, dl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("Handler: Service: get service `%s` in namespace `%s`", sid, nid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		ns  = r.Context().Value("namespace").(*types.Namespace)
	)

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in namespace `%s` err: %s", sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("Handler: Service: service `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dl, err := dm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get pod list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	pods, err := pdm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get pod list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Service().New(srv, dl, pods).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("Handler: Service: create service in namespace `%s`", nid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	opts := new(types.ServiceCreateOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	// Check memory limit reachable
	if !ns.Quotas.Disabled && opts.Spec.Memory != nil {
		if ns.Quotas.RAM < ns.Resources.RAM+(int64(*opts.Replicas)**opts.Spec.Memory) {
			log.V(logLevel).Warnf("Handler: Service: limit quotes reachable")
			errors.New("service").BadParameter("memory").Http(w)
			return
		}
	}

	srv, err := sm.Get(ns.Meta.Name, *opts.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in namespace `%s` err: %s", opts.Name, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv != nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in namespace `%s` not unique", opts.Name, ns.Meta.Name)
		errors.New("service").NotUnique("name").Http(w)
		return
	}

	srv, err = sm.Create(ns, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: create service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Service().New(srv, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("Handler: Service: update service `%s` in namespace `%s`", sid, nid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	// request body struct
	opts := new(types.ServiceUpdateOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name` err: %s", sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	// Check memory limit reachable
	if !ns.Quotas.Disabled && opts.Spec.Memory != nil {
		if ns.Quotas.RAM < ns.Resources.RAM+(int64(svc.Spec.Replicas)**opts.Spec.Memory) {
			log.V(logLevel).Warnf("Handler: Service: limit quotes reachable")
			errors.New("service").BadParameter("memory").Http(w)
			return
		}
	}

	srv, err := sm.Update(svc, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: update service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v.V1().Service().New(srv, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("Handler: Service: remove service `%s` from app `%s`", sid, nid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		ns = r.Context().Value("namespace").(*types.Namespace)
	)

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in namespace `%s` err: %s", sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	if _, err := sm.Destroy(svc); err != nil {
		log.V(logLevel).Errorf("Handler: Service: remove service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}
