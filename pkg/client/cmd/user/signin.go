package user

import (
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func SignInCmd() {

	var (
		login    string
		password string
		ctx      = context.Get()
	)

	fmt.Print("Login: ")
	fmt.Scan(&login)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
	password = string(pass)

	err = SignIn(login, password)
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
}

type userLoginS struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SignIn(login, password string) error {

	var (
		err error
		ctx = context.Get()
	)

	er := e.Http{}
	res := struct {
		Token string `json:"token"`
	}{}

	_, _, err = ctx.HTTP.
		POST("/session").
		AddHeader("Content-Type", "application/json").
		BodyJSON(userLoginS{login, password}).
		Request(&res, &er) // TODO: Need handle er
	if err != nil {
		return err
	}

	// TODO: Need handle response status code

	if res.Token == "" {
		return errors.New(e.StatusAccessDenied)
	}

	err = ctx.Session.Set(res.Token)
	if err != nil {
		return err
	}

	return nil
}
