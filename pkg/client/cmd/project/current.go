package project

import (
	"errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func CurrentCmd() {

	var ctx = context.Get()

	project, err := Current()

	if err != nil {
		ctx.Log.Error(err)
		return
	}

	project.DrawTable()
}

func Current() (*model.Project, error) {

	var (
		err     error
		ctx     = context.Get()
		project = new(model.Project)
	)

	if ctx.Token == "" {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if project.ID == "" {
		ctx.Log.Info("Project didn't select")
		return nil, nil
	}

	return project, nil
}
