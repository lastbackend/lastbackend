package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func RemoveCmd(name string) {

	var ctx = context.Get()

	err := Remove(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Successful")
}

func Remove(name string) error {

	var (
		err error
		ctx = context.Get()
		er  = new(e.Http)
		res = new(struct{})
	)

	if len(name) == 0 {
		return e.BadParameter("name").Err()
	}

	_, _, err = ctx.HTTP.
		DELETE("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(res, er)
	if err != nil {
		return errors.New(err.Error())
	}

	if er.Code == 401 {
		return e.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	project, err := Current()
	if err != nil {
		return errors.New(err.Error())
	}

	if project != nil {
		if name == project.Name {
			err = ctx.Storage.Set("project", model.Project{})
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
