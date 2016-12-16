package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func SwitchCmd(name string) {

	var ctx = context.Get()

	project, err := Switch(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Infof("The project `%s` was selected as the current", project.Name)
}

func Switch(name string) (*model.Project, error) {

	var (
		ctx     = context.Get()
		er      = new(e.Http)
		project = new(model.Project)
	)

	_, _, err := ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(&project, er)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return nil, errors.New(e.Message(er.Status))
	}

	err = ctx.Storage.Set("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return project, nil
}
