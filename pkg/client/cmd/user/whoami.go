package user

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"fmt"
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

	fmt.Println(fmt.Sprintf("Username: %s\n" +
		"E-mail: %s\nBalance: %.0f\n" +
		"Organization: %t\nCreated: %s\n" +
		"Updated: %s", res.Username, res.Email,
		res.Balance, res.Organization, res.Created.String()[:10], res.Updated.String()[:10]))

	return nil
}
