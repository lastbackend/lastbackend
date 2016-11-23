package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func GetCmd(name string) {

	var ctx = context.Get()

	err := Get(name)
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
}

func Get(name string) error {

	var (
		err   error
		ctx   = context.Get()
		token *string
	)

	token, err = ctx.Session.Get()
	if token == nil {
		return errors.New(e.StatusAccessDenied)
	}

	er := e.Http{}
	res := model.Project{}

	_, _, err = ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+*token).
		Request(&res, &er) // TODO: Need handle er
	if err != nil {
		return err
	}

	// TODO: Need handle response status code

	var header []string = []string{"ID", "Name", "Created", "Updated"}
	var data [][]string

	d := []string{
		res.ID,
		res.Name,
		res.Created.String()[:10],
		res.Updated.String()[:10],
	}

	data = append(data, d)

	table.PrintTable(header, data, []string{})

	return nil
}
