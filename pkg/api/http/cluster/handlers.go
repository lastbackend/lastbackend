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

package cluster

import (
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:cluster"
)

func ClusterInfoH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:info get cluster", logPrefix)

	var clm = distribution.NewClusterModel(r.Context(), envs.Get().GetStorage())

	cl, err := clm.Get()
	if err != nil {
		log.V(logLevel).Errorf("%s:info get cluster err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if cl == nil {
		log.V(logLevel).Errorf("%s:info cluster not found", logPrefix)
		errors.New("cluster").NotFound().Http(w)
		return
	}

	response, err := v1.View().Cluster().New(cl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:info write response err: %s", logPrefix, err.Error())
		return
	}
}

func ClusterUpdateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debugf("%s:update update cluster")

	opts := v1.Request().Cluster().UpdateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:update validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	var clm = distribution.NewClusterModel(r.Context(), envs.Get().GetStorage())

	cl, err := clm.Get()
	if err != nil {
		log.V(logLevel).Errorf("%s:update get cluster err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if cl == nil {
		err := errors.New("cluster not found")
		log.V(logLevel).Errorf("%s:update get cluster err: %s", logPrefix, err.Error())
		errors.New("cluster").NotFound().Http(w)
		return
	}

	err = clm.Update(cl, opts)
	if err != nil {
		log.V(logLevel).Errorf("%s:update get cluster err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if cl == nil {
		log.V(logLevel).Errorf("%s:update cluster not found", logPrefix)
		errors.New("cluster").NotFound().Http(w)
		return
	}

	response, err := v1.View().Cluster().New(cl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update write response err: %s", logPrefix, err.Error())
		return
	}
}
