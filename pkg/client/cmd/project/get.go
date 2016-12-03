package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func GetCmd(name string) {

	var ctx = context.Get()

	err := Get(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Get(name string) error {

	var (
		err     error
		ctx     = context.Get()
		er      = new(e.Http)
		project = new(model.Project)
	)

	if len(name) == 0 {
		return e.BadParameter("name").Err()
	}

	_, _, err = ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(&project, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code == 401 {
		return errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	project.DrawTable()

	return nil
}
