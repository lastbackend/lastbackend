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

package namespace

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/server/server/middleware"
	"github.com/lastbackend/lastbackend/internal/server/state"

	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/internal/util/http/util"
	v1 "github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:namespace"
)

// Handler represent the http handler for namespace
type Handler struct {
	state *state.State
}

// NewNamespaceHandler will initialize the namespace resources endpoint
func NewNamespaceHandler(r *mux.Router, mw middleware.Middleware, state *state.State) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init namespace routes", logPrefix)

	handler := &Handler{
		state: state,
	}

	r.Handle("/namespace", h.Handle(mw.Authenticate(handler.NamespaceListH))).Methods(http.MethodGet)
	r.Handle("/namespace", h.Handle(mw.Authenticate(handler.NamespaceCreateH))).Methods(http.MethodPost)
	r.Handle("/namespace", h.Handle(mw.Authenticate(handler.NamespaceUpdateH))).Methods(http.MethodPut)

	r.Handle("/namespace/{namespace}", h.Handle(mw.Authenticate(handler.NamespaceInfoH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}", h.Handle(mw.Authenticate(handler.NamespaceRemoveH))).Methods(http.MethodDelete)

	r.Handle("/resource", h.Handle(mw.Authenticate(handler.NamespaceResourceListH))).Methods(http.MethodGet)
	r.Handle("/resource", h.Handle(mw.Authenticate(handler.NamespaceResourceCreateH))).Methods(http.MethodPost)
	r.Handle("/resource", h.Handle(mw.Authenticate(handler.NamespaceResourceUpdateH))).Methods(http.MethodPut)

	r.Handle("/resource/{resource}/{name}", h.Handle(mw.Authenticate(handler.NamespaceResourceGetH))).Methods(http.MethodGet)
	r.Handle("/resource/{resource}/{name}", h.Handle(mw.Authenticate(handler.NamespaceResourceRemoveH))).Methods(http.MethodDelete)
}

// NamespaceListH handler returns namespaces from state
func (handler Handler) NamespaceListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace namespace namespaceList
	//
	// Shows a list of namespaces
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Environment list response
	//     schema:
	//       "$ref": "#/definitions/views_namespace_list"
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:list:> get namespace list", logPrefix)

	items := handler.state.Namespace.List(ctx)

	response, err := v1.View().Namespace().NewList(items).ToJson()
	if err != nil {
		log.Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.NotFound(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

// NamespaceInfoH handler returns particular namespace info
func (handler Handler) NamespaceInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace} namespace namespaceInfo
	//
	// Shows an info about namespace
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
	//     description: Environment response
	//     schema:
	//       "$ref": "#/definitions/views_namespace"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	nid := util.Vars(r)["namespace"]

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, nid)

	item, err := handler.state.Namespace.Get(ctx, models.NewNamespaceSelfLink(nid))
	if err != nil {
		errors.HTTP.NotFound(w)
		return
	}

	response, err := v1.View().Namespace().New(item).ToJson()
	if err != nil {
		log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

// NamespaceCreateH
func (handler Handler) NamespaceCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace namespace namespaceCreate
	//
	// Create new namespace
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_namespace_create"
	// responses:
	//   '200':
	//     description: Environment was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_namespace"
	//   '400':
	//     description: Name is already in use
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:create:> create namespace", logPrefix)

	var (
		opts = v1.Request().Namespace().Manifest()
	)

	mf, e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	item, err := handler.state.Namespace.Put(ctx, mf)
	if err != nil {
		log.Errorf("%s:create:> create namespace err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().New(item).ToJson()
	if err != nil {
		log.Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
		return
	}
}

// NamespaceUpdateH
func (handler Handler) NamespaceUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace} namespace namespaceUpdate
	//
	// Update namespace parameters
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
	//       "$ref": "#/definitions/request_namespace_update"
	// responses:
	//   '200':
	//     description: Environment was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_namespace"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	var (
		opts = v1.Request().Namespace().Manifest()
	)

	// request body struct
	mf, e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	item, err := handler.state.Namespace.Set(ctx, mf)
	if err != nil {
		log.Errorf("%s:update:> update namespace `%s` err: %s", logPrefix, mf.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().New(item).ToJson()
	if err != nil {
		log.Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

// NamespaceRemoveH
func (handler Handler) NamespaceRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace} namespace namespaceRemove
	//
	// Remove namespace
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
	//       "$ref": "#/definitions/request_namespace_remove"
	// responses:
	//   '200':
	//     description: Environment was successfully removed
	//   '403':
	//     description: Forbidden
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//nid := util.Vars(r)["namespace"]
	//
	//log.Debugf("%s:remove:> remove namespace %s", logPrefix, nid)
	//
	//var (
	//	nsm = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	//	sm  = model.NewServiceModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, err := nsm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:remove:> get namespace err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if ns == nil {
	//	err := errors.New("namespace not found")
	//	log.Errorf("%s:remove:> get namespace err: %s", logPrefix, err.Error())
	//	errors.New("namespace").NotFound().Http(w)
	//	return
	//}
	//
	//exists, err := sm.List(ns.Meta.Name)
	//if len(exists.Items) > 0 {
	//	errors.New("namespace").Forbidden().Http(w)
	//	return
	//}
	//
	//err = nsm.Remove(ns)
	//if err != nil {
	//	log.Errorf("%s:remove:> remove namespace err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) NamespaceResourceListH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")
	resource := util.QueryString(r, "resource")

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	items, err := handler.state.Resource().List(ctx, state.NewResourceFilter().WithNamespace(namespace).WithKind(resource))
	if err != nil {
		errors.HTTP.NotFound(w)
		return
	}

	response, err := v1.View().Namespace().NewResourceList(items).ToJson()
	if err != nil {
		log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}

}

func (handler Handler) NamespaceResourceGetH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")
	resource := util.Vars(r)["resource"]
	name := util.Vars(r)["name"]

	if len(namespace) == 0 {
		namespace = models.DefaultNamespace
	}

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	item, err := handler.state.Resource().Get(ctx, models.NewResourceSelfLink(namespace, resource, name))
	if err != nil {
		errors.HTTP.NotFound(w)
		return
	}

	response, err := v1.View().Namespace().NewResource(item).ToJson()
	if err != nil {
		log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}

}

func (handler Handler) NamespaceResourceCreateH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	mf, e := v1.Request().Namespace().ReadManifest(r.Body)
	if e != nil {
		log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	if len(namespace) != 0 {
		mf.SetNamespace(namespace)
	}

	item, err := handler.state.Resource().Put(ctx, mf)
	if err != nil {
		errors.HTTP.NotFound(w)
		return
	}

	response, err := v1.View().Namespace().NewResource(item).ToJson()
	if err != nil {
		log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}

}

func (handler Handler) NamespaceResourceUpdateH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")
	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	mf, e := v1.Request().Namespace().ReadManifest(r.Body)
	if e != nil {
		log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	if len(namespace) != 0 {
		mf.SetNamespace(namespace)
	}

	item, err := handler.state.Resource().Set(ctx, mf)
	if err != nil {
		errors.HTTP.NotFound(w)
		return
	}

	response, err := v1.View().Namespace().NewResource(item).ToJson()
	if err != nil {
		log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}

}

func (handler Handler) NamespaceResourceRemoveH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")
	if len(namespace) == 0 {
		namespace = models.DefaultNamespace
	}

	resource := util.Vars(r)["resource"]
	name := util.Vars(r)["name"]

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	item, err := handler.state.Resource().Del(ctx, models.NewResourceSelfLink(namespace, resource, name))
	if err != nil {
		errors.HTTP.NotFound(w)
		return
	}

	response, err := v1.View().Namespace().NewResource(item).ToJson()
	if err != nil {
		log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}

}
