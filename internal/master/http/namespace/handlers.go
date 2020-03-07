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

package namespace

import (
	"fmt"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/master/http/config/config"
	"github.com/lastbackend/lastbackend/internal/master/http/job/job"
	"github.com/lastbackend/lastbackend/internal/master/http/namespace/namespace"
	"github.com/lastbackend/lastbackend/internal/master/http/route/route"
	"github.com/lastbackend/lastbackend/internal/master/http/secret/secret"
	"github.com/lastbackend/lastbackend/internal/master/http/service/service"
	"github.com/lastbackend/lastbackend/internal/master/http/volume/volume"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"net/http"
	"strings"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/utils"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel   = 2
	logPrefix  = "api:handler:namespace"
	BufferSize = 512
)

func NamespaceListH(w http.ResponseWriter, r *http.Request) {

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

	log.V(logLevel).Debugf("%s:list:> get namespace list", logPrefix)

	var (
		nsm = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	items, err := nsm.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> find p list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().NewList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NamespaceInfoH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:info:> get namespace `%s`", logPrefix, nid)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Namespace().New(ns).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NamespaceCreateH(w http.ResponseWriter, r *http.Request) {

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

	log.V(logLevel).Debugf("%s:create:> create namespace", logPrefix)

	var (
		nsm  = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		opts = v1.Request().Namespace().Manifest()
	)

	// request body struct
	e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	item, err := nsm.Get(*opts.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> check exists by name err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item != nil {
		log.V(logLevel).Errorf("%s:create:> name `%s` not unique", logPrefix, *opts.Meta.Name)
		errors.New("namespace").NotUnique("name").Http(w)
		return
	}

	ns := new(types.Namespace)
	ns.Meta.SetDefault()
	opts.SetNamespaceMeta(ns)

	internal, _ := envs.Get().GetDomain()
	ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", ns.Meta.Name, internal))

	if err := opts.SetNamespaceSpec(ns); err != nil {
		log.V(logLevel).Errorf("%s:create:> set namespace spec err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	internal, external := envs.Get().GetDomain()

	ns.Spec.Domain.Internal = internal

	if opts.Spec.Domain != nil {
		if len(*opts.Spec.Domain) == 0 {
			ns.Spec.Domain.External = external
		} else {
			ns.Spec.Domain.External = *opts.Spec.Domain
		}
	}

	ns, err = nsm.Create(ns)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create namespace err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().New(ns).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NamespaceUpdateH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:update:> update namespace `%s`", logPrefix, nid)

	var (
		nsm  = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		opts = v1.Request().Namespace().Manifest()
	)

	// request body struct
	e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	opts.SetNamespaceMeta(ns)
	if err := opts.SetNamespaceSpec(ns); err != nil {
		log.V(logLevel).Errorf("%s:create:> set namespace spec err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	internal, _ := envs.Get().GetDomain()
	ns.Meta.Endpoint = strings.ToLower(fmt.Sprintf("%s.%s", ns.Meta.Name, internal))

	internal, external := envs.Get().GetDomain()

	ns.Spec.Domain.Internal = internal

	if opts.Spec.Domain != nil {
		if len(*opts.Spec.Domain) == 0 {
			ns.Spec.Domain.External = external
		} else {
			ns.Spec.Domain.External = *opts.Spec.Domain
		}
	}

	if err := nsm.Update(ns); err != nil {
		log.V(logLevel).Errorf("%s:update:> update namespace `%s` err: %s", logPrefix, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().New(ns).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NamespaceApplyH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]
	redeploy := utils.QueryBool(r, "redeploy")

	log.V(logLevel).Debugf("%s:apply:> apply namespace %s", logPrefix, nid)

	var (
		opts = v1.Request().Namespace().ApplyManifest()
	)

	// request body struct
	e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:apply:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	var status = struct {
		Configs  map[string]bool
		Secrets  map[string]bool
		Volumes  map[string]bool
		Services map[string]bool
		Jobs     map[string]bool
		Routes   map[string]bool
	}{
		Secrets:  make(map[string]bool, 0),
		Configs:  make(map[string]bool, 0),
		Volumes:  make(map[string]bool, 0),
		Services: make(map[string]bool, 0),
		Routes:   make(map[string]bool, 0),
		Jobs:     make(map[string]bool, 0),
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	for _, m := range opts.Configs {

		if m == nil {
			errors.New("config").BadParameter("manifest").Http(w)
			return
		}

		if m.Meta.Name == nil {
			errors.New("config").BadParameter("meta.name").Http(w)
			return
		}

		status.Configs[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	}

	for _, m := range opts.Secrets {
		if m == nil {
			errors.New("secret").BadParameter("manifest").Http(w)
			return
		}

		if m.Meta.Name == nil {
			errors.New("secret").BadParameter("meta.name").Http(w)
			return
		}

		status.Secrets[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	}

	for _, m := range opts.Volumes {
		if m == nil {
			errors.New("volume").BadParameter("manifest").Http(w)
			return
		}

		if m.Meta.Name == nil {
			errors.New("volume").BadParameter("meta.name").Http(w)
			return
		}

		status.Volumes[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	}

	for _, m := range opts.Services {
		if m == nil {
			errors.New("service").BadParameter("manifest").Http(w)
			return
		}

		if m.Meta.Name == nil {
			errors.New("service").BadParameter("meta.name").Http(w)
			return
		}
		status.Services[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	}

	for _, m := range opts.Jobs {
		if m == nil {
			errors.New("service").BadParameter("manifest").Http(w)
			return
		}

		if m.Meta.Name == nil {
			errors.New("job").BadParameter("meta.name").Http(w)
			return
		}
		status.Jobs[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	}

	for _, m := range opts.Routes {
		if m == nil {
			errors.New("route").BadParameter("manifest").Http(w)
			return
		}

		if m.Meta.Name == nil {
			errors.New("route").BadParameter("meta.name").Http(w)
			return
		}
		status.Routes[fmt.Sprintf("%s:%s", ns.SelfLink(), *m.Meta.Name)] = false
	}

	for _, m := range opts.Configs {
		c, e := config.Apply(r.Context(), ns, m)
		if e != nil {
			e.Http(w)
			return
		}
		status.Configs[c.SelfLink().String()] = true
	}

	for _, m := range opts.Secrets {
		s, e := secret.Apply(r.Context(), ns, m)
		if e != nil {
			e.Http(w)
			return
		}
		status.Secrets[s.SelfLink().String()] = true
	}

	for _, m := range opts.Volumes {
		v, e := volume.Apply(r.Context(), ns, m)
		if e != nil {
			e.Http(w)
			return
		}
		status.Volumes[v.SelfLink().String()] = true
	}

	for _, m := range opts.Services {
		s, e := service.Apply(r.Context(), ns, m, &request.ServiceUpdateOptions{Redeploy: redeploy})
		if e != nil {
			e.Http(w)
			return
		}
		status.Services[s.SelfLink().String()] = true
	}

	for _, m := range opts.Routes {
		r, e := route.Apply(r.Context(), ns, m)
		if e != nil {
			e.Http(w)
			return
		}
		status.Routes[r.SelfLink().String()] = true
	}

	for _, m := range opts.Jobs {
		j, e := job.Apply(r.Context(), ns, m)
		if e != nil {
			e.Http(w)
			return
		}
		status.Jobs[j.SelfLink().String()] = true
	}

	response, err := v1.View().Namespace().NewApplyStatus(status).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:apply:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:apply:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func NamespaceRemoveH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:remove:> remove namespace %s", logPrefix, nid)

	var (
		nsm = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = model.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get namespace err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:remove:> get namespace err: %s", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	exists, err := sm.List(ns.Meta.Name)
	if len(exists.Items) > 0 {
		errors.New("namespace").Forbidden().Http(w)
		return
	}

	err = nsm.Remove(ns)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove namespace err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
