package service

import (
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
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		GET("/service/"+name).
		AddHeader("Authorization", "Bearer "+token).
		Request(&res, req_err)

	printData(res)
	return err
}