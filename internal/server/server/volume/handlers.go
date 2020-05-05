//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/master/server/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:volume"
)

// Handler represent the http handler for volume
type Handler struct {
}

// NewVolumeHandler will initialize the volume resources endpoint
func NewVolumeHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init volume routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/namespace/{namespace}/volume", h.Handle(mw.Authenticate(handler.VolumeCreateH))).Methods(http.MethodPost)
	r.Handle("/namespace/{namespace}/volume", h.Handle(mw.Authenticate(handler.VolumeListH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/volume/{volume}", h.Handle(mw.Authenticate(handler.VolumeInfoH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/volume/{volume}", h.Handle(mw.Authenticate(handler.VolumeUpdateH))).Methods(http.MethodPut)
	r.Handle("/namespace/{namespace}/volume/{volume}", h.Handle(mw.Authenticate(handler.VolumeRemoveH))).Methods(http.MethodDelete)
}

func (handler Handler) VolumeListH(w http.ResponseWriter, r *http.Request) {

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

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:list:> get volumes list", logPrefix)
	//
	//nid := util.Vars(r)["namespace"]
	//
	//var (
	//	rm = model.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//items, err := rm.ListByNamespace(ns.Meta.Name)
	//if err != nil {
	//	log.Errorf("%s:list:> find volume list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Volume().NewList(items).ToJson()
	//if err != nil {
	//	log.Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) VolumeInfoH(w http.ResponseWriter, r *http.Request) {

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

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//nid := util.Vars(r)["namespace"]
	//rid := util.Vars(r)["volume"]
	//
	//log.Debugf("%s:info:> get volume `%s`", logPrefix, rid)
	//
	//var (
	//	rm = model.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//item, err := rm.Get(ns.Meta.Name, rid)
	//if err != nil {
	//	log.Errorf("%s:info:> find volume by id `%s` err: %s", rid, logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if item == nil {
	//	log.Warnf("%s:info:> volume `%s` not found", logPrefix, rid)
	//	errors.New("volume").NotFound().Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Volume().New(item).ToJson()
	//if err != nil {
	//	log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) VolumeCreateH(w http.ResponseWriter, r *http.Request) {

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

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:create:> create volume", logPrefix)
	//
	//nid := util.Vars(r)["namespace"]
	//
	//var (
	//	mf = v1.Request().Volume().Manifest()
	//)
	//
	//// request body struct
	//if err := mf.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//vol, e := volume.Create(r.Context(), ns, mf)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Volume().New(vol).ToJson()
	//if err != nil {
	//	log.Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) VolumeUpdateH(w http.ResponseWriter, r *http.Request) {

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

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//nid := util.Vars(r)["namespace"]
	//vid := util.Vars(r)["volume"]
	//
	//log.Debugf("%s:update:> update volume `%s`", logPrefix, nid)
	//
	//var (
	//	mf = v1.Request().Volume().Manifest()
	//)
	//
	//// request body struct
	//if e := mf.DecodeAndValidate(r.Body); e != nil {
	//	log.Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//vol, e := volume.Fetch(r.Context(), ns.Meta.Name, vid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//vol, e = volume.Update(r.Context(), ns, vol, mf)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Volume().New(vol).ToJson()
	//if err != nil {
	//	log.Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) VolumeRemoveH(w http.ResponseWriter, r *http.Request) {

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

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//nid := util.Vars(r)["namespace"]
	//rid := util.Vars(r)["volume"]
	//
	//log.Debugf("%s:remove:> remove volume %s", logPrefix, rid)
	//
	//var (
	//	rm = model.NewVolumeModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//rs, err := rm.Get(ns.Meta.Name, rid)
	//if err != nil {
	//	log.Errorf("%s:remove:> get volume by id `%s` err: %s", logPrefix, rid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if rs == nil {
	//	log.Warnf("%s:remove:> volume `%s` not found", logPrefix, rid)
	//	errors.New("volume").NotFound().Http(w)
	//	return
	//}
	//
	//err = rm.Destroy(rs)
	//if err != nil {
	//	log.Errorf("%s:remove:> remove volume `%s` err: %s", logPrefix, rid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
