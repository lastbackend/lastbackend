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
	"github.com/lastbackend/lastbackend/pkg/api/views"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
)

const logLevel = 2

func DeploymentUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]
	did := utils.Vars(r)["deployment"]

	log.V(logLevel).Debugf("Handler: Deployment: update deployment `%s` in service `%s`", sid, nid)

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
	rq := new(v1.RequestDeploymentScaleOptions)
	opts, err := rq.DecodeAndValidate(r.Body)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Deployment: validation incoming data err: %s", err)
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Deployment: get service by name` err: %s", sid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("Handler: Deployment: service name `%s` in namespace `%s` not found", sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dp, err := dm.Get(srv.Meta.Namespace, srv.Meta.Name, did)
	if err != nil {
		log.V(logLevel).Warnf("Handler: Deployment: get deployments by service failed: %s", err.Error())
		errors.New("service").NotFound().Http(w)
		return
	}

	if err := dm.Scale(dp, opts); err != nil {
		log.V(logLevel).Errorf("Handler: Deployment: update service err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	pl, err := pdm.ListByDeployment(srv.Meta.Namespace, srv.Meta.Name, dp.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Deployment: update deployment err: %s", err)
		errors.HTTP.InternalServerError(w)
	}

	response, err := views.V1().Deployment().New(dp, pl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Deployment: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Deployment: write response err: %s", err)
		return
	}
}
