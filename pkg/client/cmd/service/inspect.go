package service

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func InspectCmd(name string) {

	ctx := context.Get()

	err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}


func Inspect(name string) error {
	var (
		err   error
		ctx   = context.Get()
		token string
		res   model.Service
	)
	token, err = getToken(ctx)

	if err != nil {
		return err
	}

	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		GET("/service/"+name).
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, req_err)

	if req_err.Code != 0 {
		return errors.New(e.Message(req_err.Status))
	}

	printData(res)
	return err
}