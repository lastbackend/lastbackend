package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	em "github.com/lastbackend/lastbackend/libs/errors"
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
		err   error
		ctx   = context.Get()
		token *string
	)

	token, err = ctx.Session.Get()
	if token == nil {
		return errors.New(e.StatusAccessDenied)
	}

	er := new(e.Http)
	res := []model.Project{}

	_, _, err = ctx.HTTP.
		GET("/project").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+*token).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code != 0 {
		return errors.New(em.Message(er.Status))
	}

	var header []string = []string{"ID", "Name", "Created", "Updated"}
	var data [][]string

	for i := 0; i < len(res); i++ {
		d := []string{
			res[i].ID,
			res[i].Name,
			res[i].Created.String()[:10],
			res[i].Updated.String()[:10],
		}

		data = append(data, d)
	}

	table.PrintTable(header, data, []string{})

	return nil
}
