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

func (d *Deployer) DeployFromTemplate(userID, serviceID string, tpl model.Template) *e.Err {

	var (
		er  error
		ctx = d.ctx
	)

	service, err := ctx.Storage.Service().GetByID(userID, serviceID)
	if err != nil {
		ctx.Log.Error(err.Err())
		return err
	}

	for _, val := range tpl.Secrets {
		_, er := ctx.K8S.Core().Secrets(service.Project).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.Services {
		_, er := ctx.K8S.Core().Services(service.Project).Create(&val)
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
		_, er = ctx.K8S.Core().PersistentVolumeClaims(service.Project).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.ServiceAccounts {
		_, er = ctx.K8S.Core().ServiceAccounts(service.Project).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.Deployments {
		_, er = ctx.K8S.Extensions().Deployments(service.Project).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.ReplicationControllers {
		_, er = ctx.K8S.Core().ReplicationControllers(service.Project).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	for _, val := range tpl.Pods {
		_, er = ctx.K8S.Core().Pods(service.Project).Create(&val)
		if er != nil {
			ctx.Log.Error(er.Error())
			return e.Service.Unknown(er)
		}
	}

	return nil
}
