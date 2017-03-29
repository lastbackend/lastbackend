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
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/types"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		session      *types.Session
		project      *types.Project
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("List service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		errors.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*types.Session)

	project, err := ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	serviceList, err := ctx.Storage.Service().ListByProject(session.Username, project.Name)
	if err != nil {
		ctx.Log.Error("Error: find service list by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewServiceList(serviceList).ToJson()
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

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		session      *types.Session
		project      *types.Project
		service      *types.Service
		ctx          = c.Get()
		params       = mux.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Debug("Get service handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		errors.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*types.Session)

	project, err := ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by id", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	service, err = ctx.Storage.Service().GetByName(session.Username, project.Name, serviceParam)
	if err != nil {
		ctx.Log.Error("Error: find service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if service == nil {
		errors.New("service").NotFound().Http(w)
		return
	}

	response, err := v1.NewService(service).ToJson()
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

type serviceCreateS struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *serviceCreateS) decodeAndValidate(reader io.Reader) *errors.Err {

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

	if s.Name == nil {
		return errors.New("project").BadParameter("name")
	}

	*s.Name = strings.ToLower(*s.Name)

	if len(*s.Name) < 4 && len(*s.Name) > 64 && !validator.IsServiceName(*s.Name) {
		return errors.New("service").BadParameter("name")
	}

	if s.Description == nil {
		s.Description = new(string)
	}

	return nil
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = c.Get()
		session      *types.Session
		params       = mux.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Info("Remove service")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		errors.New("user").Unauthorized().Http(w)
		return
	}

	session = s.(*types.Session)

	projectModel, err := ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if projectModel == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	// Todo: remove all activity by service name

	err = ctx.Storage.Service().Remove(session.Username, projectModel.Name, serviceParam)
	if err != nil {
		ctx.Log.Error("Error: remove service from db", err)
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
