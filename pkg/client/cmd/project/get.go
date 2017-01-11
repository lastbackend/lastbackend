package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func GetCmd(name string) {

	var ctx = context.Get()

	project, err := Get(name)

	if err != nil {
		ctx.Log.Error(err)
		return
	}

	project.DrawTable()
}

func Get(name string) (*model.Project, error) {

	var (
		err     error
		ctx     = context.Get()
		er      = new(e.Http)
		project = new(model.Project)
	)

	if len(name) == 0 {
		return nil, e.BadParameter("name").Err()
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(&project, er)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if er.Code == 401 {
		return nil, e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return project, nil
}
