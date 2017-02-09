package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
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

	if project == nil {
		ctx.Log.Info("Project didn't select")
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
		return nil, e.NotLoggedMessage
	}

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if project.ID == "" {
		return nil, nil
	}

	return project, nil
}
