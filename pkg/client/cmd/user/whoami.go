package user

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func WhoamiCmd() {

	var (
		err error
		ctx = context.Get()
	)

	err = Whoami()
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Whoami() error {
	var (
		err error
		ctx = context.Get()
	)

	token := struct {
		Token string `json:"token"`
	}{}

	err = ctx.Storage.Get("session", &token)

	if err != nil {
		return errors.New(err.Error())
	}
	if token.Token == "" {
		return errors.New(e.Message(e.StatusAccessDenied))
	}

	er := new(e.Http)
	res := new(model.User)

	_, _, err = ctx.HTTP.
		GET("/user").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token.Token).
		Request(res, er)
	if err != nil {
		return errors.New(e.Message(err.Error()))
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	fmt.Println(fmt.Sprintf("Username: %s\n"+
		"E-mail: %s\nBalance: %.0f\n"+
		"Organization: %t\nCreated: %s\n"+
		"Updated: %s", res.Username, res.Email,
		res.Balance, res.Organization, res.Created.String()[:10], res.Updated.String()[:10]))

	return nil
}
