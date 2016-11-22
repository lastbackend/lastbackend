package user

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/table"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"strconv"
	"github.com/lastbackend/lastbackend/libs/model"
)

func WhoamiCmd() {

	var (
		err error
		ctx = context.Get()
	)

	err = Whoami()
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
}

func Whoami() error {
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
	res := model.User{}

	_, _, err = ctx.HTTP.
		GET("/user").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+*token).
		Request(&res, &er) // TODO: Need handle er
	if err != nil {
		return err
	}

	// TODO: Need handle response status code

	var header []string = []string{"Username", "Email", "Balance", "Organization", "Created", "Updated"}
	var data [][]string

	organization := strconv.FormatBool(res.Organization)
	balance := strconv.FormatFloat(float64(res.Balance), 'f', 2, 32)
	d := []string{
		res.Username,
		res.Email,
		balance,
		organization,
		res.Created.String()[:10],
		res.Updated.String()[:10],
	}

	data = append(data, d)
	d = d[:0]

	table.PrintTable(header, data, []string{})

	return nil
}
