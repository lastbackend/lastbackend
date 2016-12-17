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

	services, err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	if services != nil {
		services.DrawTable()
	}
}

func List() (*model.ServiceList, error) {

	var (
		err      error
		ctx      = context.Get()
		er       = new(e.Http)
		services = new(model.ServiceList)
	)

	project, err := p.Current()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+project.Name+"/service").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(services, er)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return nil, errors.New(e.Message(er.Status))
	}

	if len(*services) == 0 {
		ctx.Log.Info("You don't have any services")
		return nil, nil
	}

	return services, nil
}
