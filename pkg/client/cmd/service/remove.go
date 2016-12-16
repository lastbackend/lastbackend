package service

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func RemoveCmd(name string) {

	ctx := context.Get()

	err := Remove(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Successful")
}

func Remove(name string) error {

	var (
		err     error
		ctx     = context.Get()
		service = new(model.Project)
		er      = new(e.Http)
	)

	_, _, err = ctx.HTTP.
		DELETE("/service/"+name).
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(service, er)

	if err != nil {
		return err
	}

	if er.Code == 401 {
		return errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	return nil
}
