package service

import (
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
}

func Remove(name string) error {
	var (
		err   error
		ctx   = context.Get()
		token string
		res   model.Project
	)
	token, err = getToken(ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		DELETE("/service/"+name).
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, req_err)

	return err
}
