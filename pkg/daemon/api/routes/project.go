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
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/api/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func ProjectListH(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		session *types.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("List project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	projectList, err := ctx.Storage.Project().ListByUser(session.Username)
	if err != nil {
		ctx.Log.Error("Error: find projects by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProjectList(projectList).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ProjectInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *types.Session
		project      *types.Project
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Info("Get project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	if validator.IsUUID(projectParam) {
		project, err = ctx.Storage.Project().GetByID(session.Username, projectParam)
	} else {
		project, err = ctx.Storage.Project().GetByName(session.Username, projectParam)
	}
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}

}

type projectCreateS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *projectCreateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var (
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("project").IncorrectJSON(err)
	}

	if s.Name == "" {
		return errors.New("project").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsProjectName(s.Name) {
		return errors.New("project").BadParameter("name")
	}

	return nil
}

func ProjectCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		session *types.Session
		ctx     = c.Get()
	)

	ctx.Log.Debug("Create project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// request body struct
	rq := new(projectCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	project, err := ctx.Storage.Project().GetByName(session.Username, rq.Name)
	if err != nil {
		ctx.Log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if project != nil {
		errors.New("project").NotUnique("name").Http(w)
		return
	}

	project, err = ctx.Storage.Project().Insert(session.Username, rq.Name, rq.Description)
	if err != nil {
		ctx.Log.Error("Error: insert project to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

type projectUpdateS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *projectUpdateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var (
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("project").IncorrectJSON(err)
	}

	if s.Name == "" {
		return errors.New("project").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsProjectName(s.Name) {
		return errors.New("project").BadParameter("name")
	}

	return nil
}

func ProjectUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *types.Session
		project      = new(types.Project)
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("Update project handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// request body struct
	rq := new(projectUpdateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	if validator.IsUUID(projectParam) {
		project, err = ctx.Storage.Project().GetByID(session.Username, projectParam)
	} else {
		project, err = ctx.Storage.Project().GetByName(session.Username, projectParam)
	}
	if err != nil {
		ctx.Log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	project.Name = rq.Name
	project.Description = rq.Description

	project, err = ctx.Storage.Project().Update(session.Username, project)
	if err != nil {
		ctx.Log.Error("Error: update project to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ProjectRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		ctx          = c.Get()
		session      *types.Session
		project      = new(types.Project)
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Info("Remove project")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	if validator.IsUUID(projectParam) {
		project, err = ctx.Storage.Project().GetByID(session.Username, projectParam)
	} else {
		project, err = ctx.Storage.Project().GetByName(session.Username, projectParam)
	}
	if err != nil {
		ctx.Log.Error("Error: find project by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	// Todo: remove all services by project id
	// Todo: remove all activity by project id

	//err = ctx.Storage.Service().RemoveByProject(session.Username, projectParam)
	//if err != nil {
	//	ctx.Log.Error("Error: remove services from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	//err = ctx.Storage.Activity().RemoveByProject(session.Username, projectParam)
	//if err != nil {
	//	ctx.Log.Error("Error: remove activity from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	err = ctx.Storage.Project().Remove(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: remove project from db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte{})
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}

func ProjectActivityListH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ctx = c.Get()
	)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`[]`))
	if err != nil {
		ctx.Log.Error("Error: write response", err.Error())
		return
	}
}
