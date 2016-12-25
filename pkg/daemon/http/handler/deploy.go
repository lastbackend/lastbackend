package handler

import (
	"encoding/json"
	"github.com/gorilla/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/template"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"net/http"
)

type deployS struct {
	Project  *string `json:"project,omitempty"`
	Name     *string `json:"name,omitempty"`
	Template *string `json:"template,omitempty"`
	Image    *string `json:"image,omitempty"`
	Url      *string `json:"url,omitempty"`
	Config   *struct {
		Scale   *int32      `json:"scale,omitempty"`
		Ports   *portList   `json:"ports,omitempty"`
		Env     *envList    `json:"env,omitempty"`
		Volumes *volumeList `json:"volumes,omitempty"`
	} `json:"config,omitempty"`
}

type portList []template.Port
type envList []template.EnvVar
type volumeList []template.Volume

func (d *deployS) decodeAndValidate(reader io.Reader) *e.Err {

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

	if d.Image != nil && !validator.IsServiceName(*d.Image) {
		return e.New("service").BadParameter("docker")
	}

	if d.Url != nil && !validator.IsGitUrl(*d.Url) {
		ctx.Log.Error("Error: not implement")
		e.New("service").NotImplemented()
	}

	return nil
}

func DeployH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		ctx     = c.Get()
		session *model.Session
		tpl     = new(template.Template)
	)

	ctx.Log.Debug("Deploy handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").AccessDenied().Http(w)
		return
	}

	session = s.(*model.Session)

	// request body struct
	rq := new(deployS)
	if err := rq.decodeAndValidate(r.Body); err != nil {
		ctx.Log.Error("Error: validation incomming data", err.Err())
		err.Http(w)
		return
	}

	// Load template from registry
	if rq.Template != nil {
		tpl, err := template.Get(*rq.Template)
		if err == nil && tpl == nil {
			err = e.New("template").NotFound()
		}
		if err != nil {
			ctx.Log.Error("Error: deploy from template", err.Err())
			err.Http(w)
			return
		}
	}

	cfg := new(template.PatchConfig)

	// If you are not using a template, then create a standard configuration template
	if tpl == nil {
		tpl = template.CreateDefaultDeploymentConfig(*rq.Name)
	}

	// Set image as default for docker image
	if rq.Image != nil {
		cfg.Image = *rq.Image
	}

	// If have custom config, then need patch this config
	if rq.Config != nil {

		if rq.Config.Scale != nil {
			cfg.Scale = *rq.Config.Scale
		}

		if rq.Config.Ports != nil {
			for _, item := range *rq.Config.Ports {
				ports := template.Port{
					Name:          item.Name,
					ContainerPort: item.ContainerPort,
					Protocol:      item.Protocol,
				}

				cfg.Ports = append(cfg.Ports, ports)
			}
		}

		if rq.Config.Env != nil {
			for _, item := range *rq.Config.Env {
				env := template.EnvVar{
					Name:  item.Name,
					Value: item.Value,
				}

				cfg.Env = append(cfg.Env, env)
			}
		}

		if rq.Config.Volumes != nil {
			for _, item := range *rq.Config.Volumes {
				volume := template.Volume{
					Name:      item.Name,
					ReadOnly:  item.ReadOnly,
					MountPath: item.MountPath,
				}

				cfg.Volumes = append(cfg.Volumes, volume)
			}
		}

		tpl.Patch(cfg)
	}

	// Deploy from service template config
	err := tpl.Provision(*rq.Project, session.Uid, *rq.Project)
	if err != nil {
		ctx.Log.Error("Error: template provision failed", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
