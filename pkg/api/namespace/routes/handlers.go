//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package routes

import (
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/namespace/routes/request"
	"github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const logLevel = 2

func NamespaceListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
	)

	log.V(logLevel).Debug("Handler: Namespace: list namespace")

	ns := namespace.New(r.Context())
	items, err := ns.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: find namespace list err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewNamespaceList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: ", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err.Error())
		return
	}
}

func NamespaceInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		log  = context.Get().GetLogger()
		name = utils.Vars(r)["namespace"]
	)

	log.V(logLevel).Debugf("Handler: Namespace: get namespace `%s`", name)

	ns := namespace.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: find namespace by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Namespace: namespace `%s` not found", name)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	response, err := v1.NewNamespace(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err.Error())
		return
	}
}

func NamespaceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		log = context.Get().GetLogger()
	)

	log.V(logLevel).Debug("Handler: Namespace: create namespace")

	// request body struct
	rq := new(request.RequestNamespaceCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(rq.Name)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Handler: Namespace: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item != nil {
		log.Warnf("Handler: Namespace: name `%s` not unique", rq.Name)
		errors.New("namespace").NotUnique("name").Http(w)
		return
	}

	n, err := ns.Create(rq)
	response, err := v1.NewNamespace(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err.Error())
		return
	}
}

func NamespaceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		log    = context.Get().GetLogger()
		params = utils.Vars(r)
		name   = params["namespace"]
	)

	log.V(logLevel).Debugf("Handler: Namespace: update namespace `%s`", name)

	// request body struct
	rq := new(request.RequestNamespaceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: validation incoming data", err)
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: check exists by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Namespace: namespace name `%s` not found", name)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	item.Meta.Name = rq.Name
	item.Meta.Description = rq.Description

	item, err = ns.Update(item)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: update namespace `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.NewNamespace(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err.Error())
		return
	}
}

func NamespaceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		err  error
		log  = context.Get().GetLogger()
		name = utils.Vars(r)["namespace"]
	)

	log.V(logLevel).Debugf("Handler: Namespace: remove namespace %s", name)

	ns := namespace.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: find namespace by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Namespace: namespace name `%s` not found", name)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	// Todo: remove all services by namespace name
	// Todo: remove all activity by namespace name

	//err = storage.Service().RemoveByProject(session.Username, name)
	//if err != nil {
	//	log.V(logLevel).Errorf("Handler: Namespace: remove services from db err: %s", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	//err = storage.Activity().RemoveByProject(session.Username, name)
	//if err != nil {
	//	log.V(logLevel).Errorf("Handler: Namespace: remove namespace activity from db err: %s, err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	ns.Remove(item.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: remove namespace `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err.Error())
		return
	}
}

func NamespaceActivityListH(w http.ResponseWriter, r *http.Request) {

	var (
		log  = context.Get().GetLogger()
		name = utils.Vars(r)["namespace"]
	)

	log.V(logLevel).Debugf("Handler: Namespace: list namespace `%s` activity", name)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		log.V(logLevel).Errorf("Handler: Namespace: write response err: %s", err.Error())
		return
	}
}
