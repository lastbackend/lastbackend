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

package discovery

import (

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:discovery"
)

func DiscoveryInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/discovery/{discovery} discovery discoveryInfo
	//
	// Shows an discovery info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: discovery
	//     in: path
	//     description: discovery id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Discovery response
	//     schema:
	//       "$ref": "#/definitions/views_discovery_list"
	//   '404':
	//     description: Discovery not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:info:> get discovery", logPrefix)

	var (
		im  = distribution.NewDiscoveryModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["discovery"]
	)

	n, err := im.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get discovery err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:info:> discovery `%s` not found", logPrefix, nid)
		errors.New("discovery").NotFound().Http(w)
		return
	}

	response, err := v1.View().Discovery().New(n).ToJson()
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

func DiscoveryListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/discovery discovery discoveryList
	//
	// Shows an discovery list
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Discovery list response
	//     schema:
	//       "$ref": "#/definitions/views_discovery_list"
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get discoverys list", logPrefix)

	var (
		im = distribution.NewDiscoveryModel(r.Context(), envs.Get().GetStorage())
	)

	discoverys, err := im.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get discoverys list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Discovery().NewList(discoverys).ToJson()
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

func DiscoveryConnectH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/discovery/{discovery} discovery discoveryInfo
	//
	// Shows an discovery info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: discovery
	//     in: path
	//     description: discovery id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Discovery response
	//     schema:
	//       "$ref": "#/definitions/views_discovery_list"
	//   '404':
	//     description: Discovery not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:info:> discovery connect", logPrefix)

	var (
		dm  = distribution.NewDiscoveryModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["discovery"]
	)


	// request body struct
	opts := new(request.DiscoveryConnectOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	dvc, err := dm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get discovery err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if dvc == nil {
		log.V(logLevel).Debugf("%s:info:> create new discovery `%s`", logPrefix, nid)

		discovery := new(types.Discovery)
		discovery.Meta.SetDefault()
		discovery.Meta.Name = opts.Info.Hostname

		discovery.Status.Port = opts.Status.Port
		discovery.Status.IP = opts.Status.IP
		discovery.Status.Ready = opts.Status.Ready

		dm.Put(discovery)

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte{}); err != nil {
			log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
			return
		}

		return
	}

	dvc.Status.Ready = opts.Status.Ready
	dvc.Status.Port = opts.Status.Port
	dvc.Status.IP = opts.Status.IP

	if err := dm.Set(dvc); err != nil {
		log.V(logLevel).Errorf("%s:connect:> get discovery set err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}


	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func DiscoverySetStatusH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /discovery/{discovery}/status discovery discoverySetStatus
	//
	// Set discovery status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: discovery
	//     in: path
	//     description: discovery id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_discovery_status"
	// responses:
	//   '200':
	//     description: Successfully set discovery status
	//   '400':
	//     description: Bad request
	//   '404':
	//     description: Node not found / Pod not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:setstatus:> discovery set state", logPrefix)

	var (
		dm  = distribution.NewDiscoveryModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["discovery"]
	)

	// request body struct
	opts := new(request.DiscoveryStatusOptions)
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
		err.Http(w)
		return
	}

	discovery, err := dm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> get discoverys list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if discovery == nil {
		log.V(logLevel).Warnf("%s:setstatus:> update discovery `%s` not found", logPrefix, nid)
		errors.New("discovery").NotFound().Http(w)
		return
	}

	discovery.Status.Ready = opts.Ready
	discovery.Status.Port = opts.Port
	discovery.Status.IP = opts.IP

	discovery.Status.Online = true

	if err := dm.Set(discovery); err != nil {
		log.V(logLevel).Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:setstatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func DiscoveryRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /cluster/discovery/{discovery} discovery discoveryRemove
	//
	// Remove discovery
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: discovery
	//     in: path
	//     description: discovery id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Discovery removed
	//   '404':
	//     description: Discovery not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:remove:>_ create discovery", logPrefix)

	var (
		nm  = distribution.NewDiscoveryModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["discovery"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove discovery err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if n == nil {
		log.V(logLevel).Warnf("%s:remove:>_ remove discovery `%s` not found", logPrefix, nid)
		errors.New("discovery").NotFound().Http(w)
		return
	}

	if err := nm.Remove(n); err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove discovery err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:>_ write response err: %s", logPrefix, err.Error())
		return
	}
}