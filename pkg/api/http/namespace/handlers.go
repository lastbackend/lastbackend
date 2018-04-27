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

package namespace

import (
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:namespace"
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
	//     description: Namespace list response
	//     schema:
	//       "$ref": "#/definitions/views_namespace_list"
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get namespace list", logPrefix)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
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
	//     description: Namespace response
	//     schema:
	//       "$ref": "#/definitions/views_namespace"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:info:> get namespace `%s`", logPrefix, nid)

	var nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get namespace err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:info:> get namespace err: %s", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
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
	//     description: Namespace was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_namespace"
	//   '400':
	//     description: Name is already in use
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:create:> create namespace", logPrefix)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Namespace().CreateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	item, err := nsm.Get(opts.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> check exists by name err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item != nil {
		log.V(logLevel).Errorf("%s:create:> name `%s` not unique", logPrefix, opts.Name)
		errors.New("namespace").NotUnique("name").Http(w)
		return
	}

	ns, err := nsm.Create(opts)
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
	//     description: Namespace was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_namespace"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:update:> update namespace `%s`", logPrefix, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Namespace().UpdateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> get namespace err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		log.V(logLevel).Errorf("%s:update:> namespace `%s` not found", logPrefix, nid)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	if err := nsm.Update(ns, opts); err != nil {
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
	//     description: Namespace was successfully removed
	//   '403':
	//     description: Forbidden
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:remove:> remove namespace %s", logPrefix, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
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
	if len(exists) > 0 {
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
