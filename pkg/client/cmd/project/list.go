package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListCmd() {

	var ctx = context.Get()

	projects, err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	projects.DrawTable()
}

func List() (*model.ProjectList, error) {

	var (
		err      error
		ctx      = context.Get()
		er       = new(e.Http)
		projects = new(model.ProjectList)
	)

	_, _, err = ctx.HTTP.
		GET("/project").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(projects, er)
	if err != nil {
		return nil,err
	}

	if er.Code == 401 {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return nil, errors.New(e.Message(er.Status))
	}

	if len(*projects) == 0 {
		ctx.Log.Info("You don't have any projects")
		return nil, nil
	}

	return projects, nil
}
