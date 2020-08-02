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

package resource

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/server/server/legacy/middleware"
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
func NewResourceHandler(r *mux.Router, mw middleware.Middleware, state *state.State) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init namespace routes", logPrefix)

	handler := &Handler{
		state: state,
	}

	r.Handle("/resource", h.Handle(mw.Authenticate(handler.ResourceListH))).Methods(http.MethodGet)
	r.Handle("/resource", h.Handle(mw.Authenticate(handler.ResourceCreateH))).Methods(http.MethodPost)
	r.Handle("/resource", h.Handle(mw.Authenticate(handler.ResourceUpdateH))).Methods(http.MethodPut)

	r.Handle("/resource/{resource}/{name}", h.Handle(mw.Authenticate(handler.ResourceGetH))).Methods(http.MethodGet)
	r.Handle("/resource/{resource}/{name}", h.Handle(mw.Authenticate(handler.ResourceRemoveH))).Methods(http.MethodDelete)
}

// NamespaceResource management

// NamespaceResourceListH handler returns namespace resources from state by kind
// swagger:operation GET /resource namespace namespaceList
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
func (handler Handler) ResourceListH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")
	resource := util.QueryString(r, "resource")

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	items, err := handler.state.Resource(resource).List(ctx, state.NewResourceFilter().WithNamespace(namespace))
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

func (handler Handler) ResourceGetH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")
	resource := util.Vars(r)["resource"]
	name := util.Vars(r)["name"]

	if len(namespace) == 0 {
		namespace = models.DefaultNamespace
	}

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	item, err := handler.state.Resource(resource).Get(ctx, models.NewResourceSelfLink(namespace, resource, name))
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

func (handler Handler) ResourceCreateH(w http.ResponseWriter, r *http.Request) {

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

	item, err := handler.state.Resource(mf.Kind()).Put(ctx, mf)
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

func (handler Handler) ResourceUpdateH(w http.ResponseWriter, r *http.Request) {

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

	item, err := handler.state.Resource(mf.Kind()).Set(ctx, mf)
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

func (handler Handler) ResourceRemoveH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	namespace := util.QueryString(r, "namespace")

	if len(namespace) == 0 {
		namespace = models.DefaultNamespace
	}

	resource := util.Vars(r)["resource"]
	name := util.Vars(r)["name"]

	log.Debugf("%s:info:> get namespace `%s`", logPrefix, namespace)

	item, err := handler.state.Resource(resource).Del(ctx, models.NewResourceSelfLink(namespace, resource, name))
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
