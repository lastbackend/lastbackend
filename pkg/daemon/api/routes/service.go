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
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type serviceCreateS struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Template    string     `json:"template"`
	Image       string     `json:"image"`
	Url         string     `json:"url"`
	Resources   *Resources `json:"resources,omitempty"`
	source      *types.ServiceSource
}

type Resources struct {
	Region string `json:"region"`
	Memory int    `json:"memory"`
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
		return errors.New("service").IncorrectJSON(err)
	}

	if s.Template == "" && s.Image == "" && s.Url == "" {
		return errors.New("service").BadParameter("template,image,url")
	}

	if s.Template != "" {
		if s.Name == "" {
			s.Name = s.Template
		}
	}

	if s.Image != "" && s.Url == "" {
		source, err := converter.DockerNamespaceParse(s.Image)
		if err != nil {
			return errors.New("service").BadParameter("image")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}

		s.source = &types.ServiceSource{
			Type:   types.SourceDockerType,
			Hub:    source.Hub,
			Owner:  source.Owner,
			Repo:   source.Repo,
			Branch: source.Branch,
		}
	}

	if s.Url != "" {
		if !validator.IsGitUrl(s.Url) {
			return errors.New("service").BadParameter("url")
		}

		source, err := converter.GitUrlParse(s.Url)
		if err != nil {
			return errors.New("service").BadParameter("url")
		}

		if s.Name == "" {
			s.Name = source.Repo
		}

		s.source = &types.ServiceSource{
			Type:   types.SourceGitType,
			Hub:    source.Hub,
			Owner:  source.Owner,
			Repo:   source.Repo,
			Branch: source.Branch,
		}
	}

	s.Name = strings.ToLower(s.Name)

	if s.Name == "" {
		return errors.New("service").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsServiceName(s.Name) {
		return errors.New("service").BadParameter("name")
	}

	if s.Resources != nil {
		// Set default resource memory if not setting
		if s.Resources.Memory == 0 {
			s.Resources.Memory = 128
		}

		// Set default resource region if not setting
		if s.Resources.Region == "" {
			s.Resources.Region = types.WestEuropeRegion
		}
	}

	return nil
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		err          error
		session      *types.Session
		project      *types.Project
		service      *types.Service
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("Create service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	// request body struct
	rq := new(serviceCreateS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	project, err = ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if project == nil {
		errors.New("project").NotFound().Http(w)
		return
	}

	service, err = ctx.Storage.Service().GetByName(session.Username, project.Name, rq.Name)
	if err != nil {
		ctx.Log.Error("Error: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if service != nil {
		errors.New("service").NotUnique("name").Http(w)
		return
	}

	// Load template from registry
	if rq.Template != "" {
		// TODO: send request for get template config from registry
		// TODO: Set service source with types.SourceTemplateType type field
	}

	// If you are not using a template, then create a standard configuration template
	//if tpl == nil {
	// TODO: Generate default template for service
	//}

	// Patch config if exists custom configurations
	if rq.Resources != nil {
		// TODO: If have custom config, then need patch this config
	}

	config := &types.ServiceConfig{
		Replicas: 1,
		Memory:   rq.Resources.Memory,
		Region:   rq.Resources.Region,
	}

	// TODO: Create service from template

	service, err = ctx.Storage.Service().Insert(session.Username, projectParam, rq.Name, rq.Description, rq.source, config)
	if err != nil {
		ctx.Log.Error("Error: insert project to db", err)
		errors.HTTP.InternalServerError(w)
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

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		session      *types.Session
		project      *types.Project
		ctx          = c.Get()
		params       = utils.Vars(r)
		projectParam = params["project"]
	)

	ctx.Log.Debug("List service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	project, err := ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by name", err.Error())
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
		params       = utils.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Debug("Get service handler")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

	project, err := ctx.Storage.Project().GetByName(session.Username, projectParam)
	if err != nil {
		ctx.Log.Error("Error: find project by name", err.Error())
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

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		ctx          = c.Get()
		session      *types.Session
		params       = utils.Vars(r)
		projectParam = params["project"]
		serviceParam = params["service"]
	)

	ctx.Log.Info("Remove service")

	session = utils.Session(r)
	if session == nil {
		ctx.Log.Error(http.StatusText(http.StatusUnauthorized))
		errors.HTTP.Unauthorized(w)
		return
	}

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
