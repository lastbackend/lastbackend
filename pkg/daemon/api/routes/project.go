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
	"github.com/lastbackend/lastbackend/pkg/apis/views/v1"
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
		log     = c.Get().GetLogger()
		storage = c.Get().GetStorage()
	)

	log.Debug("List project handler")

	projectList, err := storage.Project().List(r.Context())
	if err != nil {
		log.Error("Error: find projects by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProjectList(projectList).ToJson()
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

func ProjectInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		err       error
		log       = c.Get().GetLogger()
		storage   = c.Get().GetStorage()
		project   *types.Project
		params    = utils.Vars(r)
		projectID = params["project"]
	)

	log.Info("Get project handler")

	if validator.IsUUID(projectID) {
		project, err = storage.Project().GetByID(r.Context(), projectID)
	} else {
		project, err = storage.Project().GetByName(r.Context(), projectID)
	}
	if err != nil {
		log.Error("Error: find project by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
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

type projectCreateS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *projectCreateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var (
		log = c.Get().GetLogger()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
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
		log     = c.Get().GetLogger()
		storage = c.Get().GetStorage()
	)

	log.Debug("Create project handler")

	// request body struct
	rq := new(projectCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	project, err := storage.Project().GetByName(r.Context(), rq.Name)
	if err != nil {
		log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if project != nil {
		errors.New("project").NotUnique("name").Http(w)
		return
	}

	project, err = storage.Project().Insert(r.Context(), rq.Name, rq.Description)
	if err != nil {
		log.Error("Error: insert project to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
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

type projectUpdateS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *projectUpdateS) decodeAndValidate(reader io.Reader) *errors.Err {

	var (
		log = c.Get().GetLogger()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
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
		log          = c.Get().GetLogger()
		storage      = c.Get().GetStorage()
		project      = new(types.Project)
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	log.Debug("Update project handler")

	// request body struct
	rq := new(projectUpdateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	if validator.IsUUID(projectParam) {
		project, err = storage.Project().GetByID(r.Context(), projectParam)
	} else {
		project, err = storage.Project().GetByName(r.Context(), projectParam)
	}
	if err != nil {
		log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	project.Meta.Name = rq.Name
	project.Meta.Description = rq.Description

	project, err = storage.Project().Update(r.Context(), project)
	if err != nil {
		log.Error("Error: update project to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewProject(project).ToJson()
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

func ProjectRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		log          = c.Get().GetLogger()
		storage      = c.Get().GetStorage()
		project      = new(types.Project)
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	log.Info("Remove project")

	if validator.IsUUID(projectParam) {
		project, err = storage.Project().GetByID(r.Context(), projectParam)
	} else {
		project, err = storage.Project().GetByName(r.Context(), projectParam)
	}
	if err != nil {
		log.Error("Error: find project by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	// Todo: remove all services by project id
	// Todo: remove all activity by project id

	//err = storage.Service().RemoveByProject(session.Username, projectParam)
	//if err != nil {
	//	log.Error("Error: remove services from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	//err = storage.Activity().RemoveByProject(session.Username, projectParam)
	//if err != nil {
	//	log.Error("Error: remove activity from db", err)
	//	e.HTTP.InternalServerError(w)
	//	return
	//}

	err = storage.Project().Remove(r.Context(), project.Meta.ID)
	if err != nil {
		log.Error("Error: remove project from db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ProjectActivityListH(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		c.Get().GetLogger().Error("Error: write response", err.Error())
		return
	}
}
