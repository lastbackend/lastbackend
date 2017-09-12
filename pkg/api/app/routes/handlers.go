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
	"github.com/lastbackend/lastbackend/pkg/api/app"
	"github.com/lastbackend/lastbackend/pkg/api/app/routes/request"
	"github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const logLevel = 2

func AppListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
	)

	log.V(logLevel).Debug("Handler: App: list app")

	ns := app.New(r.Context())
	items, err := ns.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: find app list err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewAppList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: convert struct to json err: ", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: App: write response err: %s", err.Error())
		return
	}
}

func AppInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		name = utils.Vars(r)["app"]
	)

	log.V(logLevel).Debugf("Handler: App: get app `%s`", name)

	ns := app.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: find app by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: App: app `%s` not found", name)
		errors.New("app").NotFound().Http(w)
		return
	}

	response, err := v1.NewApp(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: App: write response err: %s", err.Error())
		return
	}
}

func AppCreateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: App: create app")

	// request body struct
	rq := new(request.RequestAppCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: App: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := app.New(r.Context())
	item, err := ns.Get(rq.Name)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Handler: App: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item != nil {
		log.Warnf("Handler: App: name `%s` not unique", rq.Name)
		errors.New("app").NotUnique("name").Http(w)
		return
	}

	n, err := ns.Create(rq)
	response, err := v1.NewApp(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: App: write response err: %s", err.Error())
		return
	}
}

func AppUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		params = utils.Vars(r)
		name   = params["app"]
	)

	log.V(logLevel).Debugf("Handler: App: update app `%s`", name)

	// request body struct
	rq := new(request.RequestAppUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: App: validation incoming data", err)
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := app.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: check exists by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: App: app name `%s` not found", name)
		errors.New("app").NotFound().Http(w)
		return
	}

	item.Meta.Name = rq.Name
	item.Meta.Description = rq.Description

	item, err = ns.Update(item)
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: update app `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.NewApp(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: App: write response err: %s", err.Error())
		return
	}
}

func AppRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		err  error
		name = utils.Vars(r)["app"]
	)

	log.V(logLevel).Debugf("Handler: App: remove app %s", name)

	ns := app.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: find app by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: App: app name `%s` not found", name)
		errors.New("app").NotFound().Http(w)
		return
	}

	exist, err := ns.CheckExistServices(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: check exists services for app `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if exist {
		err := errors.New("have any services")
		log.V(logLevel).Errorf("Handler: App: remove app `%s` err: %s", name, err.Error())
		errors.Forbidden().Http(w)
		return
	}

	// TODO: remove all services by app name
	// TODO: remove all activity by app name

	ns.Remove(item.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: App: remove app `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: App: write response err: %s", err.Error())
		return
	}
}

func AppActivityListH(w http.ResponseWriter, r *http.Request) {

	var (
		name = utils.Vars(r)["app"]
	)

	log.V(logLevel).Debugf("Handler: App: list app `%s` activity", name)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		log.V(logLevel).Errorf("Handler: App: write response err: %s", err.Error())
		return
	}
}
