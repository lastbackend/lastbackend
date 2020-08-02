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
// patents in process, and are protected by trade config or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package config

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/server/server/legacy/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:config"
)

// Handler represent the http handler for config
type Handler struct {
}

// NewConfigHandler will initialize the config resources endpoint
func NewConfigHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init config routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/namespace/{namespace}/config", h.Handle(mw.Authenticate(handler.ConfigCreateH))).Methods(http.MethodPost)
	r.Handle("/namespace/{namespace}/config", h.Handle(mw.Authenticate(handler.ConfigListH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/config/{config}", h.Handle(mw.Authenticate(handler.ConfigGetH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/config/{config}", h.Handle(mw.Authenticate(handler.ConfigUpdateH))).Methods(http.MethodPut)
	r.Handle("/namespace/{namespace}/config/{config}", h.Handle(mw.Authenticate(handler.ConfigRemoveH))).Methods(http.MethodDelete)
}

func (handler Handler) ConfigGetH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/config config configList
	//
	// Shows a list of configs
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
	//     description: Config list response
	//     schema:
	//       "$ref": "#/definitions/views_config_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:list:> get config", logPrefix)
	//
	//var (
	//	sid = util.Vars(r)["config"]
	//	nid = util.Vars(r)["namespace"]
	//
	//	rm = model.NewConfigModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//item, err := rm.Get(ns.Meta.Name, sid)
	//if err != nil {
	//	log.Errorf("%s:list:> find config list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if item == nil {
	//	log.Warnf("%s:update:> config `%s` not found", logPrefix, sid)
	//	errors.New("config").NotFound().Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Config().New(item).ToJson()
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

func (handler Handler) ConfigListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/config config configList
	//
	// Shows a list of configs
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
	//     description: Config list response
	//     schema:
	//       "$ref": "#/definitions/views_config_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:list:> get configs list", logPrefix)
	//
	//var (
	//	nid = util.Vars(r)["namespace"]
	//
	//	rm = model.NewConfigModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//items, err := rm.List(ns.Meta.Name)
	//if err != nil {
	//	log.Errorf("%s:list:> find config list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Config().NewList(items).ToJson()
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

func (handler Handler) ConfigCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/config config configCreate
	//
	// Create config
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
	//       "$ref": "#/definitions/request_config_create"
	// responses:
	//   '200':
	//     description: Config was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_config"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:create:> create config", logPrefix)

	//var (
	//	nid  = util.Vars(r)["namespace"]
	//	opts = v1.Request().Config().Manifest()
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//// request body struct
	//e = opts.DecodeAndValidate(r.Body)
	//if e != nil {
	//	log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//cfg, e := config.Create(r.Context(), ns, opts)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Config().New(cfg).ToJson()
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

func (handler Handler) ConfigUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/config/{config} config configUpdate
	//
	// Create config
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
	//   - name: config
	//     in: path
	//     description: config id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_config_update"
	// responses:
	//   '200':
	//     description: Config was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_config"
	//   '404':
	//     description: Namespace not found / Config not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//var (
	//	cid = util.Vars(r)["config"]
	//	nid  = util.Vars(r)["namespace"]
	//	opts = v1.Request().Config().Manifest()
	//)

	//log.Debugf("%s:update:> update config `%s`", logPrefix, cid)

	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//// request body struct
	//e = opts.DecodeAndValidate(r.Body)
	//if e != nil {
	//	log.Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//cfg, e := config.Fetch(r.Context(), ns.Meta.Name, cid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//cfg, e = config.Update(r.Context(), ns, cfg, opts)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Config().New(cfg).ToJson()
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

func (handler Handler) ConfigRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace}/config/{config} config configRemove
	//
	// Remove config
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
	//   - name: config
	//     in: path
	//     description: config id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Config was successfully removed
	//   '404':
	//     description: Namespace not found / Config not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//var (
	//	cid = util.Vars(r)["config"]
	//	nid = util.Vars(r)["namespace"]
	//	nm  = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	//	cm  = model.NewConfigModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//log.Debugf("%s:remove:> remove config %s", logPrefix, cid)
	//
	//ns, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:list:> get namespace", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if ns == nil {
	//	err := errors.New("namespace not found")
	//	log.Errorf("%s:list:> get namespace", logPrefix, err.Error())
	//	errors.New("namespace").NotFound().Http(w)
	//	return
	//}
	//
	//cfg, err := cm.Get(ns.Meta.Name, cid)
	//if err != nil {
	//	log.Errorf("%s:remove:> get config by id `%s` err: %s", logPrefix, cid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if cfg == nil {
	//	log.Warnf("%s:remove:> config `%s` not found", logPrefix, cid)
	//	errors.New("config").NotFound().Http(w)
	//	return
	//}
	//
	//err = cm.Remove(cfg)
	//if err != nil {
	//	log.Errorf("%s:remove:> remove config `%s` err: %s", logPrefix, cid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
