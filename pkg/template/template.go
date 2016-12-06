package template

import (
	"encoding/json"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service"
	"github.com/lastbackend/lastbackend/pkg/volume"
	"io/ioutil"
)

const packageName = "template"

type Template model.Template
type TemplateList model.TemplateList

func Get(name, version string) (*Template, *e.Err) {

	var (
		er      error
		ctx     = context.Get()
		httperr = new(e.Http)
		tpl     = new(Template)
	)

	_, _, er = ctx.TemplateRegistry.
		GET(fmt.Sprintf("/template/%s/%s", name, version)).
		Request(tpl, httperr)
	if er != nil {
		return nil, e.New(packageName).Unknown(er)
	}

	if httperr.Code != 0 {
		switch httperr.Status {
		case e.StatusNotFound:
			return nil, nil
		default:
			return nil, e.New(packageName).Unknown(er)
		}
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

func (t *Template) Provision(namespace, user, project string) *e.Err {

	var (
		er  error
		ctx = context.Get()
	)

	for _, val := range t.PersistentVolumes {
		pv, err := volume.Create(user, project, &val)
		if err != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
		}

		err = pv.Deploy()
		if err != nil {
			ctx.Log.Info(err.Err())
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
		s, err := service.Create(user, project, &val)
		if err != nil {
			ctx.Log.Info(err.Err())
			return err
		}

		err = s.Deploy(namespace)
		if err != nil {
			ctx.Log.Info(err.Err())
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
		s, err := service.Create(user, project, &val)
		if err != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
		}

		err = s.Deploy(namespace)
		if er != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
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
		s, err := service.Create(user, project, &val)
		if err != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
		}

		err = s.Deploy(namespace)
		if err != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
		}
	}

	for _, val := range t.Pods {
		s, err := service.Create(user, project, &val)
		if err != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
		}

		err = s.Deploy(namespace)
		if er != nil {
			ctx.Log.Info(err.Err())
			return e.New("template").Unknown(err.Err())
		}
	}

	return nil
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
