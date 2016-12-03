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
		err  error
		ctx  = context.Get()
		er   = new(e.Http)
		user = new(model.User)
	)

	_, _, err = ctx.HTTP.
		GET("/user").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		Request(user, er)
	if err != nil {
		return errors.New(e.Message(err.Error()))
	}

	if er.Code == 401 {
		return errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	ctx.Log.Info(fmt.Sprintf("User information:\r\n\r\n"+
		"Username: \t%s\n"+
		"E-mail: \t%s\n"+
		"Balance: \t%.0f\n"+
		"Organization: \t%v\n"+
		"Created: \t%s\n"+
		"Updated: \t%s", user.Username, user.Email,
		user.Balance, user.Organization, user.Created.String()[:10], user.Updated.String()[:10]))

	return nil
}
