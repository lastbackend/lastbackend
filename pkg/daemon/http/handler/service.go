package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/deployer"
	"github.com/lastbackend/lastbackend/utils"
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
		e.User.AccessDenied().Http(w)
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
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)
	var service *model.Service

	if !utils.IsUUID(id) {
		service, err = ctx.Storage.Service().GetByName(session.Uid, id)
	} else {
		service, err = ctx.Storage.Service().GetByID(session.Uid, id)
	}

	if err == nil && service == nil {
		e.Service.NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find service by id", err.Err())
		err.Http(w)
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

type serviceCreate struct {
	Name     *string `json:"name,omitempty"`
	Project  *string `json:"project,omitempty"`
	Template *struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"template,omitempty"`
}

func (s *serviceCreate) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.User.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.Service.IncorrectJSON(err)
	}

	if s.Name == nil {
		return e.Service.BadParameter("name")
	}

	if s.Name != nil && !utils.IsServiceName(*s.Name) {
		return e.Service.BadParameter("name")
	}

	if s.Project == nil {
		return e.Service.BadParameter("project")
	}

	if s.Template != nil {
		if s.Template.Name == "" {
			return e.Service.BadParameter("template_name")
		}

		if s.Template.Version == "" {
			s.Template.Version = "latest"
		}
	} else {
		// TODO: this condition will be relevant as long as the establishment of a service only from the template
		return e.Service.BadParameter("template")
	}

	return nil
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		ctx     = c.Get()
		session *model.Session
	)

	ctx.Log.Debug("Create service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)
	ctx.Log.Debug(session.Uid)

	// request body struct
	rq := new(serviceCreate)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err.Err())
		err.Http(w)
		return
	}

	var httperr = new(e.Http)
	var tpl = new(model.Template)

	_, _, er = ctx.TemplateRegistry.
		GET(fmt.Sprintf("/template/%s/%s", rq.Template.Name, rq.Template.Version)).
		Request(tpl, httperr)
	if er != nil {
		ctx.Log.Error(httperr.Message)
		e.HTTP.InternalServerError(w)
		return
	}

	service := new(model.Service)
	service.User = session.Uid
	service.Project = *rq.Project
	service.Name = *rq.Name

	exists, er := ctx.Storage.Service().CheckExistsByName(service.User, service.Project, service.Name)
	if er != nil {
		ctx.Log.Error("Error: check exists by name", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}
	if exists {
		e.Service.NameExists().Http(w)
		return
	}

	service, err := ctx.Storage.Service().Insert(service)
	if err != nil {
		ctx.Log.Error("Error: insert service to db", err.Err())
		e.HTTP.InternalServerError(w)
		return
	}

	d := deployer.Get()

	err = d.DeployFromTemplate(service.User, service.ID, *tpl)
	if err != nil {
		ctx.Log.Error("Error: deploy service from tempalte", err.Err())
		err.Http(w)
		return
	}

	response, err := service.ToJson()
	if er != nil {
		ctx.Log.Error("Error: convert struct to json", err.Err())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

type serviceReplace struct {
	Name *string `json:"name,omitempty"`
}

func (s *serviceReplace) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.User.Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return e.Service.IncorrectJSON(err)
	}

	if s.Name != nil && !utils.IsServiceName(*s.Name) {
		return e.Service.BadParameter("name")
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
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(serviceReplace)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	if !utils.IsUUID(name) {
		service, err = ctx.Storage.Service().GetByName(session.Uid, name)
	} else {
		service, err = ctx.Storage.Service().GetByID(session.Uid, id)
	}

	if err == nil && service == nil {
		e.Service.NotFound().Http(w)
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

	if !utils.IsUUID(id) && service.Name != *rq.Name {
		exists, er := ctx.Storage.Service().CheckExistsByName(service.User, service.Project, service.Name)
		if er != nil {
			ctx.Log.Error("Error: check exists by name", er.Error())
			e.HTTP.InternalServerError(w)
			return
		}
		if exists {
			e.Service.NameExists().Http(w)
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
		e.User.AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	if !utils.IsUUID(id) {
		service, err := ctx.Storage.Service().GetByName(session.Uid, id)
		if err == nil && service == nil {
			e.Service.NotFound().Http(w)
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
