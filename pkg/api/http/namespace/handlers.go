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

const logLevel = 2

func NamespaceListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Namespace: list")

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	items, err := nsm.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: find p list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().NewList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err)
		return
	}
}

func NamespaceInfoH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("Handler: Namespace: get namespace `%s`", nid)

	var nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	response, err := v1.View().Namespace().New(ns).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err)
		return
	}
}

func NamespaceCreateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Namespace: create namespace")

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Namespace().CreateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: validation incoming data err: %s", err)
		err.Http(w)
		return
	}

	item, err := nsm.Get(opts.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: check exists by name err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if item != nil {
		log.V(logLevel).Errorf("Handler: Namespace: name `%s` not unique", opts.Name)
		errors.New("namespace").NotUnique("name").Http(w)
		return
	}

	ns, err := nsm.Create(opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: create namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().New(ns).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err)
		return
	}
}

func NamespaceUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("Handler: Namespace: update namespace `%s`", nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Namespace().UpdateOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: validation incoming data err: %s", err)
		err.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: get namespace err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		log.V(logLevel).Errorf("Handler: Namespace: namespace `%s` not found", nid)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	if err := nsm.Update(ns, opts); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: update namespace `%s` err: %s", nid, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Namespace().New(ns).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err)
		return
	}
}

func NamespaceRemoveH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("Handler: Namespace: remove namespace %s", nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Namespace().RemoveOptions()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: validation incoming data err: %s", err)
		err.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
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
		log.V(logLevel).Errorf("Handler: Namespace: remove namespace err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err)
		return
	}
}
