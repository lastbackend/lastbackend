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

package deployment

import (
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
)

const logLevel = 2

func DeploymentListH(w http.ResponseWriter, r *http.Request) {

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("api:handler:deployment:list get deployments list for `%s/%s`", sid, nid)

	var (
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:list get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:deployment:list get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:list get service by name `%s` in namespace `%s` err: %s", sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("api:handler:deployment:list service `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dl, err := dm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:list get deployment list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	pods, err := pdm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:list get pod list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Deployment().NewList(dl, pods).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:list convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:list write response err: %s", err)
		return
	}
}

func DeploymentInfoH(w http.ResponseWriter, r *http.Request) {

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]
	did := utils.Vars(r)["deployment"]

	log.V(logLevel).Debugf("api:handler:deployment:info get deployments list for `%s/%s`", sid, nid)

	var (
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:info get namespace %s err: %s", nid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err = errors.New("namespace not found")
		log.V(logLevel).Errorf("api:handler:deployment:info namespace %s not found", nid)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:info get service `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("api:handler:deployment:info service `%s` not found", sid)
		errors.New("service").NotFound().Http(w)
		return
	}

	d, err := dm.Get(srv.Meta.Namespace, srv.Meta.Name, did)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:info get deployment by name `%s` err: %s", did, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if d == nil {
		log.V(logLevel).Warnf("api:handler:deployment:info deployment `%s` not found", did)
		errors.New("deployment").NotFound().Http(w)
		return
	}

	pods, err := pdm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:info get pod list by service id `%s` err: %s", srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Deployment().New(d, pods).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:info convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:info write response err: %s", err)
		return
	}
}

func DeploymentUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]
	did := utils.Vars(r)["deployment"]

	log.V(logLevel).Debugf("api:handler:deployment:update update deployment `%s` in service `%s`", sid, nid)

	if r.Context().Value("namespace") == nil {
		errors.HTTP.Forbidden(w)
		return
	}

	var (
		err error
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		ns  = r.Context().Value("namespace").(*types.Namespace)
	)

	// request body struct
	opts := v1.Request().Deployment().UpdateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:update validation incoming data err: %s", err)
		err.Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:update get service by name` err: %s", sid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("api:handler:deployment:update service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dp, err := dm.Get(srv.Meta.Namespace, srv.Meta.Name, did)
	if err != nil {
		log.V(logLevel).Warnf("api:handler:deployment:update get deployments by service failed: %s", err.Error())
		errors.New("service").NotFound().Http(w)
		return
	}

	if err := dm.SetSpec(dp, opts); err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:update update service err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	pl, err := pdm.ListByDeployment(srv.Meta.Namespace, srv.Meta.Name, dp.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:update update deployment err: %s", err)
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.View().Deployment().New(dp, pl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:update convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("api:handler:deployment:update write response err: %s", err)
		return
	}
}
