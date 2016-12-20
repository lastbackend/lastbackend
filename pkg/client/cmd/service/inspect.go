package service

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func InspectCmd(name string) {

	ctx := context.Get()

	service, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	service.DrawTable()
}

func Inspect(name string) (*model.Service, error) {

	var (
		err     error
		ctx     = context.Get()
		er      = new(e.Http)
		service = new(model.Service)
		project = new(model.Project)
	)

	if len(name) == 0 {
		return nil, e.BadParameter("name").Err()
	}

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if project.ID == "" {
		return nil, errors.New("Project didn't select")
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+project.ID+"/service/"+name).
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(service, er)

	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code == 404 {
		return nil, errors.New("Service "+name+" not found.")
	}

	if er.Code != 0 {
		return nil, errors.New(e.Message(er.Status))
	}

	return service, nil
}
