package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *model.Session
		projectModel *model.Project
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("List service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*model.Session)

	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceModel, err := ctx.Storage.Service().ListByProject(session.Uid, projectModel.ID)
	if err != nil {
		ctx.Log.Error("Error: find services by user", err)
		e.HTTP.InternalServerError(w)
		return
	}

	servicesSpec, err := service.List(ctx.K8S, projectModel.ID)
	if err != nil {
		ctx.Log.Error("Error: get serivce spec from cluster", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	var list = model.ServiceList{}
	var response []byte

	if serviceModel != nil {
		for _, val := range *serviceModel {
			val.Spec = servicesSpec["lb-"+val.ID]
			list = append(list, val)
		}

		response, err = list.ToJson()
		if err != nil {
			ctx.Log.Error("Error: convert struct to json", err.Error())
			e.HTTP.InternalServerError(w)
			return
		}

	} else {
		response, err = serviceModel.ToJson()
		if err != nil {
			ctx.Log.Error("Error: convert struct to json", err.Error())
			e.HTTP.InternalServerError(w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *model.Session
		projectModel *model.Project
		serviceModel *model.Service
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Debug("Get service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*model.Session)

	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceModel, err = ctx.Storage.Service().GetByNameOrID(session.Uid, projectModel.ID, serviceParam)
	if err == nil && serviceModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find service by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceSpec, err := service.Get(ctx.K8S, serviceModel.Project, "lb-"+serviceModel.ID)
	if err != nil {
		ctx.Log.Error("Error: get serivce spec from cluster", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceModel.Spec = serviceSpec

	response, err := serviceModel.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

type serviceReplaceS struct {
	*model.ServiceUpdateConfig
}

func (s *serviceReplaceS) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.New("service").IncorrectJSON(err)
	}

	return nil
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *model.Session
		projectModel *model.Project
		serviceModel *model.Service
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Debug("Update service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(serviceReplaceS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		e.HTTP.InternalServerError(w)
		return
	}

	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceModel, err = ctx.Storage.Service().GetByNameOrID(session.Uid, projectModel.ID, serviceParam)
	if err == nil && serviceModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find service by name or id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceModel.Description = rq.Description

	serviceModel, err = ctx.Storage.Service().Update(serviceModel)
	if err != nil {
		ctx.Log.Error("Error: insert service to db", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	cfg := rq.CreateServiceConfig()

	err = service.Update(ctx.K8S, serviceModel.Project, serviceModel.ID, cfg)
	if err != nil {
		ctx.Log.Error("Error: update service", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := serviceModel.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		er           error
		ctx          = c.Get()
		session      *model.Session
		projectModel *model.Project
		params       = mux.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Info("Remove service")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*model.Session)

	projectModel, err := ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	if !validator.IsUUID(serviceParam) {
		serviceModel, err := ctx.Storage.Service().GetByName(session.Uid, projectModel.ID, serviceParam)
		if err == nil && serviceModel == nil {
			e.New("service").NotFound().Http(w)
			return
		}
		if err != nil {
			ctx.Log.Error("Error: find service by id", err.Error())
			e.HTTP.InternalServerError(w)
			return
		}

		serviceParam = serviceModel.ID
	}

	// TODO: Clear entities from kubernetes

	err = ctx.Storage.Service().Remove(session.Uid, projectParam, serviceParam)
	if err != nil {
		ctx.Log.Error("Error: remove service from db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceLogsH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx            = c.Get()
		session        *model.Session
		projectModel   *model.Project
		serviceModel   *model.Service
		params         = mux.Vars(r)
		projectParam   = params["project"]
		serviceParam   = params["service"]
		query          = r.URL.Query()
		podQuery       = query.Get("pod")
		containerQuery = query.Get("container")
		ch             = make(chan bool, 1)
		notify         = w.(http.CloseNotifier).CloseNotify()
	)

	ctx.Log.Info("Show service log")

	go func() {
		<-notify
		ctx.Log.Debug("HTTP connection just closed.")
		ch <- true
	}()

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*model.Session)

	projectModel, err := ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	serviceModel, err = ctx.Storage.Service().GetByNameOrID(session.Uid, projectModel.ID, serviceParam)
	if err == nil && serviceModel == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find service by id", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	opts := service.ServiceLogsOption{
		Stream:     w,
		Pod:        podQuery,
		Container:  containerQuery,
		Follow:     true,
		Timestamps: true,
	}

	service.Logs(ctx.K8S, serviceModel.Project, &opts, ch)
}
