package project

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

type createS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func CreateCmd(name, description string) {

	var ctx = context.Get()

	err := Create(name, description)
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
}

func Create(name, description string) error {

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
		POST("/project").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+*token).
		BodyJSON(createS{name, description}).
		Request(&res, &er) // TODO: Need handle er
	if err != nil {
		return err
	}

	// TODO: Need handle response status code

	return nil
}
