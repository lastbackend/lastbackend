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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const logLevel = 2

func ClusterInfoH(w http.ResponseWriter, r *http.Request) {

	name := utils.Vars(r)["cluster"]

	log.V(logLevel).Debugf("Handler: Cluster: get cluster `%s`", name)

	response, err := v1.View().Cluster().New(nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: write response err: %s", err)
		return
	}
}

func ClusterUpdateH(w http.ResponseWriter, r *http.Request) {

	name := utils.Vars(r)["cluster"]

	log.V(logLevel).Debugf("Handler: Cluster: update cluster `%s`", name)

	opts := v1.Request().Cluster().UpdateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: validation incoming data", err)
		err.Http(w)
		return
	}

	response, err := v1.View().Cluster().New(nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: write response err: %s", err)
		return
	}
}

func ClusterRemoveH(w http.ResponseWriter, r *http.Request) {

	name := utils.Vars(r)["cluster"]

	log.V(logLevel).Debugf("Handler: Cluster: remove cluster %s", name)

	opts := v1.Request().Cluster().UpdateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: validation incoming data", err)
		err.Http(w)
		return
	}

	response, err := v1.View().Cluster().New(nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Cluster: write response err: %s", err)
		return
	}
}
