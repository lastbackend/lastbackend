package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/common"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/volume"
	"io/ioutil"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"strings"
)

const packageName = "template"

type Template model.Template
type TemplateList model.TemplateList

type PatchConfig struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Scale   int32    `json:"scale"`
	Ports   []Port   `json:"ports"`
	Env     []EnvVar `json:"env"`
	Volumes []Volume `json:"volumes"`
}

// Port represents a network port in a single container
type Port struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int32  `json:"container"`
	Protocol      string `json:"protocol,omitempty"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// VolumeMount describes a mounting of a Volume within a container.
type Volume struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readonly,omitempty"`
	MountPath string `json:"mountpath"`
}

func Get(name string) (*Template, *e.Err) {

	var (
		er      error
		ctx     = context.Get()
		httperr = new(e.Http)
		tpl     = new(Template)
	)

	parts := strings.Split(name, ":")

	name = parts[0]
	version := "latest"

	if len(parts) == 2 {
		version = parts[1]
	}

	_, _, er = ctx.TemplateRegistry.
		GET(fmt.Sprintf("/template/%s/%s", name, version)).
		Request(tpl, httperr)
	if er != nil {
		return nil, e.New(packageName).Unknown(er)
	}

	if httperr.Code == 404 {
		return nil, nil
	}

	if httperr.Code != 0 {
		return nil, e.New(packageName).Unknown(er)
	}

	return tpl, nil
}

func List() (*TemplateList, *e.Err) {

	var (
		er        error
		ctx       = context.Get()
		templates = new(TemplateList)
	)

	_, resp, er := ctx.TemplateRegistry.GET("/template").Do()
	if er != nil {
		return nil, e.New(packageName).Unknown(er)
	}

	buf, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		return nil, e.New(packageName).Unknown(er)
	}

	er = json.Unmarshal(buf, templates)
	if er != nil {
		return nil, e.New(packageName).Unknown(er)
	}

	return templates, nil
}

func CreateDefaultDeploymentConfig(name string) *Template {

	dp := new(v1beta1.Deployment)

	common.Set_defaults_v1beta1_deployment(dp)

	dp.Name = name
	dp.Labels = map[string]string{
		"app": name,
	}

	dp.Spec.Template.Labels = map[string]string{
		"app":  name,
		"role": "placeholder",
	}

	dp.Spec.Template.Name = name
	dp.Spec.Template.Spec.Containers = make([]v1.Container, 1)
	dp.Spec.Template.Spec.Containers[0].Name = name
	dp.Spec.Template.Spec.Containers[0].Image = "alpine"
	dp.Spec.Template.Spec.Containers[0].ImagePullPolicy = v1.PullAlways

	var tpl = new(Template)
	tpl.Deployments = make([]v1beta1.Deployment, 1)
	tpl.Deployments[0] = *dp

	return tpl
}

func (t *Template) Provision(namespace, user, project string, patchConfig *PatchConfig) *e.Err {

	t.Patch(patchConfig)

	var (
		er         error
		ctx        = context.Get()
		deployFunc = func(val interface{}) *e.Err {

			var config *v1beta1.Deployment

			switch val.(type) {
			case *v1beta1.Deployment:
				config = val.(*v1beta1.Deployment)
			case *v1.ReplicationController:
				config = converter.Convert_ReplicationController_to_Deployment(val.(*v1.ReplicationController))
			case *v1.Pod:
				config = converter.Convert_Pod_to_Deployment(val.(*v1.Pod))
			default:
				return e.New("service").Unknown(errors.New("unknown type config"))
			}

			var serviceModel = new(model.Service)
			serviceModel.User = user
			serviceModel.Project = project
			serviceModel.Name = config.Name

			serviceModel, err := ctx.Storage.Service().Insert(serviceModel)
			if err != nil {
				return err
			}

			config.Name = serviceModel.ID

			_, err = service.Deploy(ctx.K8S, namespace, config)
			if err != nil {
				return err
			}

			return nil
		}
	)

	for _, val := range t.PersistentVolumes {
		pv, err := volume.Create(user, project, &val)
		if err != nil {
			return e.New("template").Unknown(err.Err())
		}

		err = pv.Deploy()
		if err != nil {
			return e.New("template").Unknown(err.Err())
		}
	}

	for _, val := range t.Services {
		_, er = ctx.K8S.Core().Services(namespace).Create(&val)
		if er != nil {
			ctx.Log.Info(er.Error())
			return e.New("template").Unknown(er)
		}
	}

	for _, val := range t.Secrets {
		_, er = ctx.K8S.Core().Secrets(namespace).Create(&val)
		if er != nil {
			ctx.Log.Info(er.Error())
			return e.New("template").Unknown(er)
		}
	}

	for _, val := range t.Deployments {
		err := deployFunc(&val)
		if err != nil {
			return err
		}
	}

	for _, val := range t.PersistentVolumeClaims {
		_, er = ctx.K8S.Core().PersistentVolumeClaims(namespace).Create(&val)
		if er != nil {
			ctx.Log.Info(er.Error())
			return e.New("template").Unknown(er)
		}
	}

	for _, val := range t.ServiceAccounts {
		_, er = ctx.K8S.Core().ServiceAccounts(namespace).Create(&val)
		if er != nil {
			ctx.Log.Info(er.Error())
			return e.New("template").Unknown(er)
		}
	}

	for _, val := range t.DaemonSets {
		_, er = ctx.K8S.Extensions().DaemonSets(namespace).Create(&val)
		if er != nil {
			ctx.Log.Info(er.Error())
			return e.New("template").Unknown(er)
		}
	}

	for _, val := range t.Jobs {
		err := deployFunc(&val)
		if err != nil {
			return err
		}
	}

	for _, val := range t.Ingresses {
		_, er = ctx.K8S.Extensions().Ingresses(namespace).Create(&val)
		if er != nil {
			ctx.Log.Info(er.Error())
			return e.New("template").Unknown(er)
		}
	}

	for _, val := range t.ReplicationControllers {
		err := deployFunc(&val)
		if err != nil {
			return err
		}
	}

	for _, val := range t.Pods {
		err := deployFunc(&val)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) Patch(config *PatchConfig) {
	if config == nil {
		return
	}

	for dp_index := range t.Deployments {
		var dp = &t.Deployments[dp_index]

		dp.Name =  config.Name+"-"+dp.Name

		if _, ok := dp.Spec.Template.Labels["role"]; ok && dp.Spec.Template.Labels["role"] == "placeholder" {

			dp.Name =  config.Name

			delete(dp.Spec.Template.Labels, "role")

			if config.Scale > 1 {
				dp.Spec.Replicas = &config.Scale
			}

			for c_index := range dp.Spec.Template.Spec.Containers {

				var c = &dp.Spec.Template.Spec.Containers[c_index]

				if config.Image != "" {
					c.Image = config.Image
				}

				for _, p := range config.Ports {
					c.Ports = append(c.Ports, v1.ContainerPort{
						Protocol:      v1.Protocol(p.Protocol),
						ContainerPort: p.ContainerPort,
						Name:          p.Name,
					})
				}

				for _, env := range config.Env {
					c.Env = append(c.Env, v1.EnvVar{
						Name:  env.Name,
						Value: env.Value,
					})
				}

				for _, volume := range config.Volumes {
					c.VolumeMounts = append(c.VolumeMounts, v1.VolumeMount{
						Name:      volume.Name,
						ReadOnly:  volume.ReadOnly,
						MountPath: volume.MountPath,
					})
				}
			}
		}
	}

	return
}

func (t *Template) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, e.New("template").Unknown(err)
	}

	return buf, nil
}

func (t *TemplateList) ToJson() ([]byte, *e.Err) {

	if t == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(t)
	if err != nil {
		return nil, e.New("template").Unknown(err)
	}

	return buf, nil
}
