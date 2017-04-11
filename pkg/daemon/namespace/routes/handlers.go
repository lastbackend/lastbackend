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
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace/routes/request"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

func NamespaceListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
	)

	log.Debug("List project handler")

	ns := namespace.New()
	items, err := ns.List(r.Context())
	if err != nil {
		log.Error("Error: find namespcaes", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewNamespaceList(items).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func NamespaceInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		id  = utils.Vars(r)["namespace"]
	)

	log.Info("Get namespace handler")
	ns := namespace.New()
	item, err := ns.Get(r.Context(), id)
	if err != nil {
		log.Error("Error: find namespace by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	response, err := v1.NewNamespace(item).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}

}

func NamespaceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		log = context.Get().GetLogger()
	)

	log.Debug("Create project handler")

	// request body struct
	rq := new(request.RequestNamespaceCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New()
	item, err := ns.Get(r.Context(), rq.Name)
	if err != nil {
		log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if item != nil {
		errors.New("project").NotUnique("name").Http(w)
		return
	}

	n, err := ns.Create(r.Context(), rq)
	response, err := v1.NewNamespace(n).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func NamespaceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		log    = context.Get().GetLogger()
		params = utils.Vars(r)
		id     = params["namespace"]
	)

	log.Debug("Update project handler")

	// request body struct
	rq := new(request.RequestNamespaceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New()
	item, err := ns.Get(r.Context(), id)
	if err != nil {
		log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	item.Meta.Name = rq.Name
	item.Meta.Description = rq.Description

	item, err = ns.Update(r.Context(), item)
	if err != nil {
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.NewNamespace(item).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func NamespaceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		id  = utils.Vars(r)["namespace"]
	)

	log.Info("Remove namespace")
	ns := namespace.New()
	item, err := ns.Get(r.Context(), id)
	if err != nil {
		log.Error("Error: find project by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if item == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	// Todo: remove all services by project id
	// Todo: remove all activity by project id

	//err = storage.Service().RemoveByProject(session.Username, id)
	//if err != nil {
	//	log.Error("Error: remove services from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	//err = storage.Activity().RemoveByProject(session.Username, id)
	//if err != nil {
	//	log.Error("Error: remove activity from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	ns.Remove(r.Context(), item.Meta.ID)
	if err != nil {
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func NamespaceActivityListH(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		context.Get().GetLogger().Error("Error: write response", err.Error())
		return
	}
}
