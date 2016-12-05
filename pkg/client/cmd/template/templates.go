package template

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListCmd() {

	var ctx = context.Get()

	err := List()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func List() error {

	var (
		ctx       = context.Get()
		er        = new(e.Http)
		templates = new(model.TemplateList)
	)

	_, _, err := ctx.HTTP.
		GET("/template").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(&templates, er)

	if err != nil {
		return err
	}

	if er.Code == 401 {
		return errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	templates.DrawTable()

	return err
}
