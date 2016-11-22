package user

import (
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func SignUpCmd() {

	var (
		err      error
		ctx      = context.Get()
		username string
		email    string
		password string
	)

	fmt.Print("Username: ")
	fmt.Scan(&username)

	fmt.Print("Email: ")
	fmt.Scan(&email)

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
	password = string(pass)

	err = SignUp(username, email, password)
	if err != nil {
		ctx.Log.Error(err) // TODO: Need handle error and print to console
		return
	}
}

type userCreateS struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(username, email, password string) error {

	var (
		err error
		ctx = context.Get()
	)

	er := e.Http{}
	res := struct {
		Token string `json:"token"`
	}{}

	_, _, err = ctx.HTTP.
		POST("/user").
		AddHeader("Content-Type", "application/json").
		BodyJSON(userCreateS{username, email, password}).
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
