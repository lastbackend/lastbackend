package service

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	p "github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func InspectCmd(name string) {

	ctx := context.Get()

	service, projectName, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	service.DrawTable(projectName)
}

func Inspect(name string) (*model.Service, string, error) {

	var (
		err     error
		ctx     = context.Get()
		er      = new(e.Http)
		service = new(model.Service)
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
		GET("/project/"+project.Name+"/service/"+name).
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(service, er)

	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, "", errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return nil, "", errors.New(e.Message(er.Status))
	}

	return service, project.Name, nil
}
