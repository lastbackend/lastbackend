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
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/http/util"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/master/server/middleware"
	"github.com/lastbackend/lastbackend/internal/master/state"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	v1 "github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix  = "api:handler:namespace"
	BufferSize = 512
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
	r.Handle("/namespace/{namespace}", h.Handle(mw.Authenticate(handler.NamespaceInfoH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}", h.Handle(mw.Authenticate(handler.NamespaceUpdateH))).Methods(http.MethodPut)
	r.Handle("/namespace/{namespace}/apply", h.Handle(mw.Authenticate(handler.NamespaceApplyH))).Methods(http.MethodPut)
	r.Handle("/namespace/{namespace}", h.Handle(mw.Authenticate(handler.NamespaceRemoveH))).Methods(http.MethodDelete)
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

	item, err := handler.state.Namespace.Get(ctx, types.NewNamespaceSelfLink(nid))
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

	//log.Debugf("%s:create:> create namespace", logPrefix)
	//
	//var (
	//	nsm  = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	//	opts = v1.Request().Namespace().Manifest()
	//)
	//
	//// request body struct
	//e := opts.DecodeAndValidate(r.Body)
	//if e != nil {
	//	log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//item, err := nsm.Get(*opts.Meta.Name)
	//if err != nil {
	//	log.Errorf("%s:create:> check exists by name err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if item != nil {
	//	log.Errorf("%s:create:> name `%s` not unique", logPrefix, *opts.Meta.Name)
	//	errors.New("namespace").NotUnique("name").Http(w)
	//	return
	//}
	//
	//ns := new(types.Namespace)
	//ns.Meta.SetDefault()
	//opts.SetNamespaceMeta(ns)
	//
	//internal, external := handler.Config.DomainInternal, handler.Config.DomainExternal
	//ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", ns.Meta.Name, internal))
	//
	//if err := opts.SetNamespaceSpec(ns); err != nil {
	//	log.Errorf("%s:create:> set namespace spec err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//

	//
	//ns.Spec.Domain.Internal = internal
	//
	//if opts.Spec.Domain != nil {
	//	if len(*opts.Spec.Domain) == 0 {
	//		ns.Spec.Domain.External = external
	//	} else {
	//		ns.Spec.Domain.External = *opts.Spec.Domain
	//	}
	//}
	//
	//ns, err = nsm.Create(ns)
	//if err != nil {
	//	log.Errorf("%s:create:> create namespace err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Namespace().New(ns).ToJson()
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

	//nid := util.Vars(r)["namespace"]
	//
	//log.Debugf("%s:update:> update namespace `%s`", logPrefix, nid)
	//
	//var (
	//	nsm  = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	//	opts = v1.Request().Namespace().Manifest()
	//)
	//
	//// request body struct
	//e := opts.DecodeAndValidate(r.Body)
	//if e != nil {
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
	//opts.SetNamespaceMeta(ns)
	//if err := opts.SetNamespaceSpec(ns); err != nil {
	//	log.Errorf("%s:create:> set namespace spec err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//internal, external := handler.Config.DomainInternal, handler.Config.DomainExternal
	//ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", ns.Meta.Name, internal))
	//
	//ns.Spec.Domain.Internal = internal
	//
	//if opts.Spec.Domain != nil {
	//	if len(*opts.Spec.Domain) == 0 {
	//		ns.Spec.Domain.External = external
	//	} else {
	//		ns.Spec.Domain.External = *opts.Spec.Domain
	//	}
	//}
	//
	//if err := nsm.Update(ns); err != nil {
	//	log.Errorf("%s:update:> update namespace `%s` err: %s", logPrefix, nid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Namespace().New(ns).ToJson()
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

func (handler Handler) NamespaceApplyH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace} namespace namespaceApply
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
	//       "$ref": "#/definitions/request_namespace_apply"
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

	//nid := util.Vars(r)["namespace"]
	//redeploy := util.QueryBool(r, "redeploy")
	//
	//log.Debugf("%s:apply:> apply namespace %s", logPrefix, nid)
	//
	//var (
	//	opts = v1.Request().Namespace().ApplyManifest()
	//)
	//
	//// request body struct
	//e := opts.DecodeAndValidate(r.Body)
	//if e != nil {
	//	log.Errorf("%s:apply:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//var status = struct {
	//	Configs  map[string]bool
	//	Secrets  map[string]bool
	//	Volumes  map[string]bool
	//	Services map[string]bool
	//	Jobs     map[string]bool
	//	Routes   map[string]bool
	//}{
	//	Secrets:  make(map[string]bool, 0),
	//	Configs:  make(map[string]bool, 0),
	//	Volumes:  make(map[string]bool, 0),
	//	Services: make(map[string]bool, 0),
	//	Routes:   make(map[string]bool, 0),
	//	Jobs:     make(map[string]bool, 0),
	//}
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//for _, m := range opts.Configs {
	//
	//	if m == nil {
	//		errors.New("config").BadParameter("manifest").Http(w)
	//		return
	//	}
	//
	//	if m.Meta.Name == nil {
	//		errors.New("config").BadParameter("meta.name").Http(w)
	//		return
	//	}
	//
	//	status.Configs[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	//}
	//
	//for _, m := range opts.Secrets {
	//	if m == nil {
	//		errors.New("secret").BadParameter("manifest").Http(w)
	//		return
	//	}
	//
	//	if m.Meta.Name == nil {
	//		errors.New("secret").BadParameter("meta.name").Http(w)
	//		return
	//	}
	//
	//	status.Secrets[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	//}
	//
	//for _, m := range opts.Volumes {
	//	if m == nil {
	//		errors.New("volume").BadParameter("manifest").Http(w)
	//		return
	//	}
	//
	//	if m.Meta.Name == nil {
	//		errors.New("volume").BadParameter("meta.name").Http(w)
	//		return
	//	}
	//
	//	status.Volumes[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	//}
	//
	//for _, m := range opts.Services {
	//	if m == nil {
	//		errors.New("service").BadParameter("manifest").Http(w)
	//		return
	//	}
	//
	//	if m.Meta.Name == nil {
	//		errors.New("service").BadParameter("meta.name").Http(w)
	//		return
	//	}
	//	status.Services[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	//}
	//
	//for _, m := range opts.Jobs {
	//	if m == nil {
	//		errors.New("service").BadParameter("manifest").Http(w)
	//		return
	//	}
	//
	//	if m.Meta.Name == nil {
	//		errors.New("job").BadParameter("meta.name").Http(w)
	//		return
	//	}
	//	status.Jobs[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	//}
	//
	//for _, m := range opts.Routes {
	//	if m == nil {
	//		errors.New("route").BadParameter("manifest").Http(w)
	//		return
	//	}
	//
	//	if m.Meta.Name == nil {
	//		errors.New("route").BadParameter("meta.name").Http(w)
	//		return
	//	}
	//	status.Routes[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	//}
	//
	//for _, m := range opts.Configs {
	//	c, e := config.Apply(r.Context(), ns, m)
	//	if e != nil {
	//		e.Http(w)
	//		return
	//	}
	//	status.Configs[c.SelfLink().String()] = true
	//}
	//
	//for _, m := range opts.Secrets {
	//	s, e := secret.Apply(r.Context(), ns, m)
	//	if e != nil {
	//		e.Http(w)
	//		return
	//	}
	//	status.Secrets[s.SelfLink().String()] = true
	//}
	//
	//for _, m := range opts.Volumes {
	//	v, e := volume.Apply(r.Context(), ns, m)
	//	if e != nil {
	//		e.Http(w)
	//		return
	//	}
	//	status.Volumes[v.SelfLink().String()] = true
	//}
	//
	//for _, m := range opts.Services {
	//	s, e := service.Apply(r.Context(), ns, m, &request.ServiceUpdateOptions{Redeploy: redeploy})
	//	if e != nil {
	//		e.Http(w)
	//		return
	//	}
	//	status.Services[s.SelfLink().String()] = true
	//}
	//
	//for _, m := range opts.Routes {
	//	r, e := route.Apply(r.Context(), ns, m)
	//	if e != nil {
	//		e.Http(w)
	//		return
	//	}
	//	status.Routes[r.SelfLink().String()] = true
	//}
	//
	//for _, m := range opts.Jobs {
	//	j, e := job.Apply(r.Context(), ns, m)
	//	if e != nil {
	//		e.Http(w)
	//		return
	//	}
	//	status.Jobs[j.SelfLink().String()] = true
	//}
	//
	//response, err := v1.View().Namespace().NewApplyStatus(status).ToJson()
	//if err != nil {
	//	log.Errorf("%s:apply:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:apply:> write response err: %s", logPrefix, err.Error())
		return
	}
}

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
