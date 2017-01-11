package template

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func ListCmd() {

	var ctx = context.Get()

	templates, err := List()

	if err != nil {
		ctx.Log.Error(err)
		return
	}

	templates.DrawTable()
}

func List() (*model.TemplateList, error) {

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
		return nil, err
	}

	if er.Code == 401 {
		return nil, e.NotLoggedMessage
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return templates, err
}
