package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"net/http"
)

func ProjectListH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("List project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	projectListModel, err := ctx.Storage.Project().ListByUser(session.Uid)
	if err != nil {
		ctx.Log.Error("Error: find projects by user", err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := projectListModel.ToJson()
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

func ProjectInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		er           error
		err          *e.Err
		session      *model.Session
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("Get project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)
	var projectModel *model.Project

	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("project").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Err())
		err.Http(w)
		return
	}

	response, err := projectModel.ToJson()
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

type projectCreateS struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *projectCreateS) decodeAndValidate(reader io.Reader) *e.Err {

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
		return e.New("project").IncorrectJSON(err)
	}

	if s.Name == nil {
		return e.New("project").BadParameter("name")
	}

	if s.Name != nil && !validator.IsProjectName(*s.Name) {
		return e.New("project").BadParameter("name")
	}

	if s.Description == nil {
		s.Description = new(string)
	}

	return nil
}

func ProjectCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		session *model.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("Create project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(projectCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	projectModel := new(model.Project)
	projectModel.User = session.Uid
	projectModel.Name = *rq.Name
	projectModel.Description = *rq.Description

	exists, er := ctx.Storage.Project().ExistByName(projectModel.User, projectModel.Name)
	if er != nil {
		ctx.Log.Error("Error: check exists by name", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}
	if exists {
		e.New("project").NotUnique("name").Http(w)
		return
	}

	project, err := ctx.Storage.Project().Insert(projectModel)
	if err != nil {
		ctx.Log.Error("Error: insert project to db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	namespace := &v1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:      project.ID,
			Namespace: project.ID,
			Labels: map[string]string{
				"user": session.Username,
			},
		},
	}

	_, er = ctx.K8S.Core().Namespaces().Create(namespace)
	if er != nil {
		ctx.Log.Error("Error: create namespace", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := project.ToJson()
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

type projectReplaceS struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *projectReplaceS) decodeAndValidate(reader io.Reader) *e.Err {

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
		return e.New("project").IncorrectJSON(err)
	}

	if s.Name != nil && !validator.IsProjectName(*s.Name) {
		return e.New("project").BadParameter("name")
	}

	if s.Description == nil {
		s.Description = new(string)
	}

	return nil
}

func ProjectUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		er           error
		err          *e.Err
		session      *model.Session
		projectModel *model.Project
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("Update project handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(projectReplaceS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		err.Http(w)
		return
	}

	projectModel, err = ctx.Storage.Project().GetByNameOrID(session.Uid, projectParam)
	if err == nil && projectModel == nil {
		e.New("project").NotFound().Http(w)
		return
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Err())
		err.Http(w)
		return
	}

	if rq.Name == nil || *rq.Name == "" {
		rq.Name = &projectModel.Name
	}

	if !validator.IsUUID(projectParam) && projectModel.Name != *rq.Name {
		exists, er := ctx.Storage.Project().ExistByName(projectModel.User, *rq.Name)
		if er != nil {
			e.HTTP.InternalServerError(w)
		}
		if exists {
			e.New("project").NotUnique("name").Http(w)
			return
		}
	}

	projectModel.Name = *rq.Name
	projectModel.Description = *rq.Description

	projectModel, err = ctx.Storage.Project().Update(projectModel)

	if err != nil {
		ctx.Log.Error("Error: insert project to db", err.Err())
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := projectModel.ToJson()
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

func ProjectRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		er           error
		ctx          = c.Get()
		session      *model.Session
		params       = mux.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Info("Remove project")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	if !validator.IsUUID(projectParam) {
		projectModel, err := ctx.Storage.Project().GetByName(session.Uid, projectParam)
		if err == nil && projectModel == nil {
			e.New("project").NotFound().Http(w)
			return
		}
		if err != nil {
			ctx.Log.Error("Error: find project by id", err.Err())
			err.Http(w)
			return
		}

		projectParam = projectModel.ID
	}

	var opts = new(api.DeleteOptions)

	er = ctx.K8S.Core().Namespaces().Delete(projectParam, opts)
	if er != nil {
		ctx.Log.Error("Error: remove namespace", er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	volumes, err := ctx.Storage.Volume().ListByProject(projectParam)
	if err != nil {
		ctx.Log.Error("Error: get volumes from db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	if volumes != nil {
		for _, val := range *volumes {
			er = ctx.K8S.Core().PersistentVolumes().Delete(val.Name, &api.DeleteOptions{})
			if er != nil {
				ctx.Log.Error("Error: remove persistent volume", er.Error())
				e.HTTP.InternalServerError(w)
				return
			}

			err := ctx.Storage.Volume().Remove(val.ID)
			if err != nil {
				ctx.Log.Error("Error: remove volume from db", err)
				e.HTTP.InternalServerError(w)
				return
			}
		}
	}

	err = ctx.Storage.Service().RemoveByProject(session.Uid, projectParam)
	if err != nil {
		ctx.Log.Error("Error: remove services from db", err)
		e.HTTP.InternalServerError(w)
		return
	}

	err = ctx.Storage.Project().Remove(session.Uid, projectParam)
	if err != nil {
		ctx.Log.Error("Error: remove project from db", err)
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
