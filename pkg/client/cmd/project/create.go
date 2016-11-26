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
		ctx.Log.Error(err)
		return
	}
}

func Create(name, description string) error {

	var (
		err   error
		ctx   = context.Get()

	)
	token := struct {
		Token string `json:"token"`
	}{}
	err = ctx.Storage.Get("session", &token)
	if token.Token == "" {
		return errors.New(e.StatusAccessDenied)
	}

	er := new(e.Http)
	res := new(model.Project)

	_, _, err = ctx.HTTP.
		POST("/project").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer " + token.Token).
		BodyJSON(createS{name, description}).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	return nil
}
