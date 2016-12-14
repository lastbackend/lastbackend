package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("List service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	services, err := ctx.Storage.Service().GetByUser(session.Uid)
	if err != nil {
		ctx.Log.Error("Error: find services by user", err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := services.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
		params  = mux.Vars(r)
		id      = params["id"]
	)

	ctx.Log.Debug("Get service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)
	var sm *model.Service

	if !validator.IsUUID(id) {
		sm, err = ctx.Storage.Service().GetByName(session.Uid, id)
	} else {
		sm, err = ctx.Storage.Service().GetByID(session.Uid, id)
	}

	if err == nil && sm == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find service by id", err.Err())
		err.Http(w)
		return
	}

	serviceSpec, err := service.Get(sm.Project, sm.Name)
	if err != nil {
		ctx.Log.Error("Error: get serivce spec from cluster", err.Err())
		err.Http(w)
		return
	}

	buf, er := json.Marshal(serviceSpec)
	if er != nil {
		ctx.Log.Error("Error: parse spec", er.Error())
		e.Unknown(er).Http(w)
		return
	}

	ctx.Log.Info(">>>>>>>> ", string(buf))

	response, err := sm.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

type serviceReplaceS struct {
	Name *string `json:"name,omitempty"`
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

	if s.Name != nil && !validator.IsServiceName(*s.Name) {
		return e.New("service").BadParameter("name")
	}

	return nil
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		service *model.Service
		ctx     = c.Get()
		params  = mux.Vars(r)
		id      = params["id"]
		name    = params["id"]
	)

	ctx.Log.Debug("Update service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(serviceReplaceS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	if !validator.IsUUID(name) {
		service, err = ctx.Storage.Service().GetByName(session.Uid, name)
	} else {
		service, err = ctx.Storage.Service().GetByID(session.Uid, id)
	}

	if err == nil && service == nil {
		e.New("service").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find service by name or id", err.Err())
		err.Http(w)
		return
	}

	if rq.Name == nil || *rq.Name == "" {
		rq.Name = &service.Name
	}

	service.Name = *rq.Name

	if !validator.IsUUID(id) && service.Name != *rq.Name {
		exists, er := ctx.Storage.Service().CheckExistsByName(service.User, service.Project, service.Name)
		if er != nil {
			ctx.Log.Error("Error: check exists by name", er.Error())
			e.HTTP.InternalServerError(w)
			return
		}
		if exists {
			e.New("service").NotUnique("name").Http(w)
			return
		}
	}

	service, err = ctx.Storage.Service().Update(service)
	if err != nil {
		ctx.Log.Error("Error: insert service to db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := service.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		ctx     = c.Get()
		session *model.Session
		params  = mux.Vars(r)
		id      = params["id"]
	)

	ctx.Log.Info("Remove service")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	if !validator.IsUUID(id) {
		service, err := ctx.Storage.Service().GetByName(session.Uid, id)
		if err == nil && service == nil {
			e.New("service").NotFound().Http(w)
			return
		}
		if err != nil {
			ctx.Log.Error("Error: find service by id", err.Err())
			err.Http(w)
			return
		}

		id = service.ID
	}

	// TODO: Clear entities from kubernetes

	err := ctx.Storage.Service().Remove(id)
	if err != nil {
		ctx.Log.Error("Error: remove service from db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
