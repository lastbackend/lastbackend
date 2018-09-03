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

package ingress

import (
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
	logPrefix = "api:handler:ingress"
)

func IngressInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress/{ingress} ingress ingressInfo
	//
	// Shows an ingress info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:info:> get ingress", logPrefix)

	var (
		im  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["ingress"]
	)

	n, err := im.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if n == nil {
		log.V(logLevel).Warnf("%s:info:> ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	response, err := v1.View().Ingress().New(n).ToJson()
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

func IngressListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress ingress ingressList
	//
	// Shows an ingress list
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Ingress list response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get ingresss list", logPrefix)

	var (
		im = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
	)

	ingresss, err := im.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get ingresss list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Ingress().NewList(ingresss).ToJson()
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

func IngressRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /cluster/ingress/{ingress} ingress ingressRemove
	//
	// Remove ingress
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress removed
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:remove:>_ create ingress", logPrefix)

	var (
		nm  = distribution.NewIngressModel(r.Context(), envs.Get().GetStorage())
		nid = utils.Vars(r)["ingress"]
	)

	n, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if n == nil {
		log.V(logLevel).Warnf("%s:remove:>_ remove ingress `%s` not found", logPrefix, nid)
		errors.New("ingress").NotFound().Http(w)
		return
	}

	if err := nm.Remove(n); err != nil {
		log.V(logLevel).Errorf("%s:remove:>_ remove ingress err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:>_ write response err: %s", logPrefix, err.Error())
		return
	}
}
