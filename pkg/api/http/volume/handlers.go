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

package volume

import (
	"github.com/lastbackend/lastbackend/pkg/api/http/namespace/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/http/volume/volume"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"

	"net/http"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:volume"
)

func VolumeListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/volume volume volumeList
	//
	// Shows a list of volumes
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
	//     description: Volume list response
	//     schema:
	//       "$ref": "#/definitions/views_volume_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get volumes list", logPrefix)

	nid := utils.Vars(r)["namespace"]

	var (
		rm = distribution.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	items, err := rm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> find volume list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Volume().NewList(items).ToJson()
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

func VolumeInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/volume/{volume} volume volumeInfo
	//
	// Shows an info about volume
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
	//   - name: volume
	//     in: path
	//     description: volume id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Volume response
	//     schema:
	//       "$ref": "#/definitions/views_volume"
	//   '404':
	//     description: Namespace not found / Volume not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["volume"]

	log.V(logLevel).Debugf("%s:info:> get volume `%s`", logPrefix, rid)

	var (
		rm = distribution.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	item, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> find volume by id `%s` err: %s", rid, logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("%s:info:> volume `%s` not found", logPrefix, rid)
		errors.New("volume").NotFound().Http(w)
		return
	}

	response, err := v1.View().Volume().New(item).ToJson()
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

func VolumeCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/volume volume volumeCreate
	//
	// Creates a volume
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
	//       "$ref": "#/definitions/request_volume_create"
	// responses:
	//   '200':
	//     description: Volume was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_volume"
	//   '400':
	//     description: Bad rules parameter
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:create:> create volume", logPrefix)

	nid := utils.Vars(r)["namespace"]

	var (
		mf = v1.Request().Volume().Manifest()
	)

	// request body struct
	if err := mf.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, err.Err())
		err.Http(w)
		return
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	vol, e := volume.Create(r.Context(), ns, mf)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Volume().New(vol).ToJson()
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

func VolumeUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/volume/{volume} volume volumeUpdate
	//
	// Update volume
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
	//   - name: volume
	//     in: path
	//     description: volume id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_volume_update"
	// responses:
	//   '200':
	//     description: Volume was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_volume"
	//   '400':
	//     description: Bad rules parameter
	//   '404':
	//     description: Namespace not found / Volume not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	vid := utils.Vars(r)["volume"]

	log.V(logLevel).Debugf("%s:update:> update volume `%s`", logPrefix, nid)

	var (
		mf = v1.Request().Volume().Manifest()
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

	vol, e := volume.Fetch(r.Context(), ns.Meta.Name, vid)
	if e != nil {
		e.Http(w)
		return
	}

	vol, e = volume.Update(r.Context(), ns, vol, mf)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Volume().New(vol).ToJson()
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

func VolumeRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace}/volume/{volume} volume volumeRemove
	//
	// Removes volume
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
	//   - name: volume
	//     in: path
	//     description: volume id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Volume was successfully removed
	//   '404':
	//     description: Namespace not found / Volume not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	rid := utils.Vars(r)["volume"]

	log.V(logLevel).Debugf("%s:remove:> remove volume %s", logPrefix, rid)

	var (
		rm = distribution.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	rs, err := rm.Get(ns.Meta.Name, rid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get volume by id `%s` err: %s", logPrefix, rid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if rs == nil {
		log.V(logLevel).Warnf("%s:remove:> volume `%s` not found", logPrefix, rid)
		errors.New("volume").NotFound().Http(w)
		return
	}

	err = rm.Destroy(rs)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove volume `%s` err: %s", logPrefix, rid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
