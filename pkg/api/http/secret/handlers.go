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
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const logLevel = 2

func SecretListH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Secret: list")

	nid := utils.Vars(r)["namespace"]

	var (
		rm  = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	items, err := rm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: find secret list err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Secret().NewList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Secret: write response err: %s", err)
		return
	}
}

func SecretCreateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Secret: create secret")

	nid := utils.Vars(r)["namespace"]

	var (
		rm  = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Secret().CreateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("Handler: Secret: validation incoming data err: %s", e.Err())
		e.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	rs, err := rm.Create(ns, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: create secret err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Secret().New(rs).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Secret: write response err: %s", err)
		return
	}
}

func SecretUpdateH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["secret"]

	log.V(logLevel).Debugf("Handler: Secret: update secret `%s`", nid)

	var (
		rm  = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts, e := v1.Request().Secret().UpdateOptions().DecodeAndValidate(r.Body)
	if e != nil {
		log.V(logLevel).Errorf("Handler: Secret: validation incoming data err: %s", e.Err())
		e.Http(w)
		return
	}

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	ss, err := rm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: check secret exists by selflink `%s` err: %s", ns.Meta.SelfLink, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ss == nil {
		log.V(logLevel).Warnf("Handler: Secret: secret `%s` not found", sid)
		errors.New("secret").NotFound().Http(w)
		return
	}

	ss, err = rm.Update(ss, ns, opts)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: update secret `%s` err: %s", ns.Meta.Name, err)
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.View().Secret().New(ss).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: convert struct to json err: %s", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Secret: write response err: %s", err)
		return
	}
}

func SecretRemoveH(w http.ResponseWriter, r *http.Request) {

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["secret"]

	log.V(logLevel).Debugf("Handler: Secret: remove secret %s", sid)

	var (
		sm  = distribution.NewSecretModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: get namespace", err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("Handler: Namespace: get namespace", err)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	ss, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: get secret by id `%s` err: %s", sid, err)
		errors.HTTP.InternalServerError(w)
		return
	}
	if ss == nil {
		log.V(logLevel).Warnf("Handler: Secret: secret `%s` not found", sid)
		errors.New("secret").NotFound().Http(w)
		return
	}

	err = sm.Remove(ss)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Secret: remove secret `%s` err: %s", sid, err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Secret: write response err: %s", err)
		return
	}
}
