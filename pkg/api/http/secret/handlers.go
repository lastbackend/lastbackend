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

package secret

import (
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:secret"
)

func SecretGetH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/secret secret secretList
	//
	// Shows a list of secrets
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
	//     description: Secret list response
	//     schema:
	//       "$ref": "#/definitions/views_secret_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get secret", logPrefix)

	var (
		sid = utils.Vars(r)["secret"]
		nid = utils.Vars(r)["namespace"]
		nm  = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		rm  = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	item, err := rm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> find secret list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if item == nil {
		log.V(logLevel).Warnf("%s:update:> secret `%s` not found", logPrefix, sid)
		errors.New("secret").NotFound().Http(w)
		return
	}

	response, err := v1.View().Secret().New(item).ToJson()
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

func SecretListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/secret secret secretList
	//
	// Shows a list of secrets
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
	//     description: Secret list response
	//     schema:
	//       "$ref": "#/definitions/views_secret_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:list:> get secrets list", logPrefix)

	var (
		nid = utils.Vars(r)["namespace"]
		nm  = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		rm = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	items, err := rm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> find secret list err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Secret().NewList(items).ToJson()
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

func SecretCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/secret secret secretCreate
	//
	// Create secret
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
	//       "$ref": "#/definitions/request_secret_create"
	// responses:
	//   '200':
	//     description: Secret was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_secret"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	log.V(logLevel).Debugf("%s:create:> create secret", logPrefix)

	var (
		nid  = utils.Vars(r)["namespace"]
		nm   = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm   = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
		opts = v1.Request().Secret().Manifest()
	)

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	// request body struct
	e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ss := new(types.Secret)
	opts.SetSecretMeta(ss)
	opts.SetSecretSpec(ss)

	rs, err := sm.Create(ns, ss)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create secret err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Secret().New(rs).ToJson()
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

func SecretUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/secret/{secret} secret secretUpdate
	//
	// Create secret
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
	//   - name: secret
	//     in: path
	//     description: secret id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_secret_update"
	// responses:
	//   '200':
	//     description: Secret was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_secret"
	//   '404':
	//     description: Namespace not found / Secret not found
	//   '500':
	//     description: Internal server error

	sid := utils.Vars(r)["secret"]

	log.V(logLevel).Debugf("%s:update:> update secret `%s`", logPrefix, sid)

	var (
		nid = utils.Vars(r)["namespace"]
		nm  = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		rm = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
		opts = v1.Request().Secret().Manifest()
	)

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	// request body struct
	e := opts.DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ss, err := rm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> check secret exists by selflink err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ss == nil {
		log.V(logLevel).Warnf("%s:update:> secret `%s` not found", logPrefix, sid)
		errors.New("secret").NotFound().Http(w)
		return
	}

	opts.SetSecretMeta(ss)
	opts.SetSecretSpec(ss)

	ss, err = rm.Update(ss)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> update secret `%s` err: %s", logPrefix, ss.Meta.SelfLink, err.Error())
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.View().Secret().New(ss).ToJson()
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

func SecretRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace}/secret/{secret} secret secretRemove
	//
	// Remove secret
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
	//   - name: secret
	//     in: path
	//     description: secret id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Secret was successfully removed
	//   '404':
	//     description: Namespace not found / Secret not found
	//   '500':
	//     description: Internal server error

	sid := utils.Vars(r)["secret"]

	log.V(logLevel).Debugf("%s:remove:> remove secret %s", logPrefix, sid)

	var (
		nid = utils.Vars(r)["namespace"]
		nm  = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	ss, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get secret by id `%s` err: %s", logPrefix, sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ss == nil {
		log.V(logLevel).Warnf("%s:remove:> secret `%s` not found", logPrefix, sid)
		errors.New("secret").NotFound().Http(w)
		return
	}

	err = sm.Remove(ss)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove secret `%s` err: %s", logPrefix, sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
