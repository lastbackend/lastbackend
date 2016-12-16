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
	)

	_, _, err = ctx.HTTP.
		GET("/service/"+name).
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(service, er)

	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return nil, errors.New(e.Message(er.Status))
	}

	return service, nil
}
