package service

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	p "github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListServiceCmd() {

	var ctx = context.Get()

	services, projectName, err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	if services != nil {
		services.DrawTable(projectName)
	}
}

func List() (*model.ServiceList, string, error) {

	var (
		err      error
		ctx      = context.Get()
		er       = new(e.Http)
		services = new(model.ServiceList)
	)

	project, err := p.Current()
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if project == nil {
		ctx.Log.Info("Project didn't select")
		return nil, "", nil
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+project.Name+"/service").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(services, er)
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, "", e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, "", errors.New(er.Message)
	}

	if len(*services) == 0 {
		ctx.Log.Info("You don't have any services")
		return nil, "", nil
	}

	return services, project.Name, nil
}
