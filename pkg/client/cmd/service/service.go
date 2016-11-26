package service

import (
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
)

type serviceCreate struct {
	Name string `json:"name"`

}

func getToken(ctx *context.Context) (string, error) {
	var err error
	token := struct {
		Token string `json:"token"`
	}{}
	err = ctx.Storage.Get("session", &token)
	return token.Token, err
}


func Create(name string) error {
	var (
		err error
		ctx = context.Get()
		token string
		res model.Project
	)
	token, err = getToken(&ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		POST("/service").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer " + token).
		BodyJSON(serviceCreate{name}).
		Request(&res, req_err)
	if err != nil {}
	return err

}

func Inspect(name string) error {
	var (
		err error
		ctx = context.Get()
		token string
		res model.Project
	)
	token, err = getToken(&ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		GET("/service/" + name).
		AddHeader("Authorization", "Bearer " + token).
		Request(&res, req_err)
	if err != nil {}
	return err
}

func Remove(name string) error {
	var (
		err error
		ctx = context.Get()
		token string
		res model.Project
	)
	token, err = getToken(&ctx)
	req_err := new(e.Http)
	_, _, err = ctx.HTTP.
		DELETE("/service/" + name).
		AddHeader("Authorization", "Bearer " + token).
		Request(&res, req_err)
	if err != nil {}
	return err
}