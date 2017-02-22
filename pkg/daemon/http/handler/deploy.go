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
  Project  string `json:"project,omitempty"`
  Name     string `json:"name,omitempty"`
  Template string `json:"template,omitempty"`
  Image    string `json:"image,omitempty"`
  Url      string `json:"url,omitempty"`
  Config *struct {
    Scale   int32      `json:"scale,omitempty"`
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

  if d.Project == "" {
    return e.New("service").BadParameter("project")
  }

  if d.Template != "" {
    if d.Name == "" {
      d.Name = d.Template
    }
  }

  if d.Image != "" && d.Url == ""  {
    if !validator.IsServiceName(d.Image) {
      return e.New("service").BadParameter("docker")
    }

    source, err := converter.DockerNamespaceParse(d.Image)
    if err != nil {
      return e.New("service").BadParameter("image")
    }

    if d.Name == "" {
      d.Name = source.Repo
    }
  }

  if d.Url != "" {
    if !validator.IsGitUrl(d.Url) {
      e.New("service").BadParameter("url")
    }

    source, err := converter.GitUrlParse(d.Url)
    if err != nil {
      return e.New("service").BadParameter("url")
    }

    if d.Name == "" {
      d.Name = source.Repo
    }
  }

  return nil
}

func DeployH(w http.ResponseWriter, r *http.Request) {

  var (
    err     error
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
    e.HTTP.InternalServerError(w)
    return
  }

  // Load template from registry
  if rq.Template != "" {
    tpl, err = template.Get(rq.Template)
    if err == nil && tpl == nil {
      ctx.Log.Error("Error: tempalte " + rq.Template + " not found")
      e.New("template").NotFound().Http(w)
      return
    }
    if err != nil {
      ctx.Log.Error("Error: deploy from template", err.Error())
      e.HTTP.InternalServerError(w)
      return
    }
  }

  // If you are not using a template, then create a standard configuration template
  if tpl == nil {
    tpl = template.CreateDefaultDeploymentConfig(rq.Name)
  }

  cfg.Name = rq.Name

  // Set image as default for docker image
  if rq.Image != "" {
    cfg.Image = rq.Image
  }

  // If have custom config, then need patch this config
  if rq.Config != nil {

    if rq.Config.Scale != 0 {
      cfg.Scale = rq.Config.Scale
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

  projectModel, err := ctx.Storage.Project().GetByNameOrID(session.Uid, rq.Project)
  if err == nil && projectModel == nil {
    e.New("service").NotFound().Http(w)
    return
  }
  if err != nil {
    ctx.Log.Error("Error: find project by id", err.Error())
    e.HTTP.InternalServerError(w)
    return
  }

  // Deploy from service template config
  err = tpl.Provision(projectModel.ID, session.Uid, projectModel.ID, cfg)
  if err != nil {
    ctx.Log.Error("Error: template provision failed", err.Error())
    e.HTTP.InternalServerError(w)
    return
  }

  w.WriteHeader(http.StatusOK)
  _, err = w.Write([]byte{})
  if err != nil {
    ctx.Log.Error("Error: write response", err.Error())
    return
  }
}
