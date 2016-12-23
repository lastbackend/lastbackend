package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/template"
	"io"
	"io/ioutil"
	"net/http"
)

type deployTemplateS struct {
	Project *string `json:"project,omitempty"`
	Target  *string `json:"target,omitempty"`
}

func (d *deployTemplateS) decodeAndValidate(reader io.Reader) *e.Err {

	var (
		err error
		ctx = c.Get()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Error(err)
		return e.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return e.New("service").IncorrectJSON(err)
	}

	if d.Project == nil {
		return e.New("service").BadParameter("project")
	}

	if d.Target == nil {
		return e.New("service").BadParameter("target")
	}

	return nil
}

func TemplateListH(w http.ResponseWriter, _ *http.Request) {

	var (
		er             error
		ctx            = c.Get()
		response_empty = func() {
			w.WriteHeader(200)
			_, er = w.Write([]byte("[]"))
			if er != nil {
				ctx.Log.Error("Error: write response", er.Error())
				return
			}
			return
		}
	)

	templates, err := template.List()
	if err != nil {
		ctx.Log.Error(err.Err())
		response_empty()
		return
	}

	if templates == nil {
		response_empty()
		return
	}

	response, err := templates.ToJson()
	if er != nil {
		ctx.Log.Error(err.Err())
		response_empty()
		return
	}

	w.WriteHeader(200)
	_, er = w.Write(response)
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func TemplateDeployH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		ctx     = c.Get()
		session *model.Session
	)

	ctx.Log.Debug("Deploy template")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	tpl := template.CreateDefaultDeploymentConfig("redis", "redis")

	srv, err := service.Create(&tpl.Deployments[0])
	if err != nil {
		ctx.Log.Info(err.Err())
		return
	}

	detail, err := srv.Deploy(ctx.K8S, "785f7418-d731-48d6-868b-074ebd749ae3")
	if err != nil {
		ctx.Log.Info(err.Err())
		return
	}

	var serviceModel = new(model.Service)
	serviceModel.User = session.Uid
	serviceModel.Project = "785f7418-d731-48d6-868b-074ebd749ae3"
	serviceModel.Name = detail.ObjectMeta.Name

	serviceModel, err = ctx.Storage.Service().Insert(serviceModel)
	if err != nil {
		return
	}

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
