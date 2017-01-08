package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/template"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
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

	if d.Template != nil {
		if d.Name == nil {
			d.Name = d.Template
		}
	}

	if d.Image != nil && d.Url == nil {
		if !validator.IsServiceName(*d.Image) {
			return e.New("service").BadParameter("docker")
		}

		source, err := converter.DockerNamespaceParse(*d.Image)
		if err != nil {
			return e.New("service").BadParameter("image")
		}

		if d.Name == nil {
			d.Name = &source.Repo
		}
	}

	if d.Url != nil {
		if !validator.IsGitUrl(*d.Url) {
			e.New("service").BadParameter("url")
		}

		source, err := converter.GitUrlParse(*d.Url)
		if err != nil {
			return e.New("service").BadParameter("url")
		}

		if d.Name == nil {
			d.Name = &source.Repo
		}
	}

	return nil
}

func DeployH(w http.ResponseWriter, r *http.Request) {

	var (
		er      error
		err     *e.Err
		ctx     = c.Get()
		session *model.Session
		tpl     *template.Template
		cfg     = new(template.PatchConfig)
	)

	ctx.Log.Debug("Deploy handler")

	s, ok := context.GetOk(r, `session`)
	if !ok {
		ctx.Log.Error("Error: get session context")
		e.New("user").Unauthorized().Http(w)
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
		tpl, err = template.Get(*rq.Template)
		if err == nil && tpl == nil {
			ctx.Log.Error("Error: tempalte " + *rq.Template + " not found")
			e.New("template").NotFound().Http(w)
			return
		}
		if err != nil {
			ctx.Log.Error("Error: deploy from template", err.Err())
			err.Http(w)
			return
		}
	}

	// If you are not using a template, then create a standard configuration template
	if tpl == nil {
		tpl = template.CreateDefaultDeploymentConfig(*rq.Name)
	}

	cfg.Name = *rq.Name

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
	}

	// Deploy from service template config
	err = tpl.Provision(*rq.Project, session.Uid, *rq.Project, cfg)
	if err != nil {
		ctx.Log.Error("Error: template provision failed", err.Err())
		err.Http(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
