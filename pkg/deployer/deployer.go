package deployer

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
)

var deployer *Deployer

type Deployer struct {
	ctx *context.Context
}

func Get() *Deployer {

	if deployer == nil {
		deployer = &Deployer{
			ctx: context.Get(),
		}
	}

	return deployer
}

func (d *Deployer) DeployProjectFromTemplate(namespace string, tpl model.Template) *e.Err {

	var (
		er  error
		ctx = d.ctx
	)

	for _, val := range tpl.Secrets {
		_, er := ctx.K8S.Core().Secrets(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.PersistentVolumes {
		_, er = ctx.K8S.Core().PersistentVolumes().Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.PersistentVolumeClaims {
		_, er = ctx.K8S.Core().PersistentVolumeClaims(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.Services {
		_, er := ctx.K8S.Core().Services(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.ServiceAccounts {
		_, er = ctx.K8S.Core().ServiceAccounts(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.Deployments {
		_, er = ctx.K8S.Extensions().Deployments(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.ReplicationControllers {
		_, er = ctx.K8S.Core().ReplicationControllers(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.Pods {
		_, er = ctx.K8S.Core().Pods(namespace).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	return nil
}
