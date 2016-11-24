package user

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	m "github.com/lastbackend/lastbackend/libs/errors"
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
		err   error
		ctx   = context.Get()
		token *string
	)

	token, err = ctx.Session.Get()
	if token == nil {
		return errors.New(m.Message(e.StatusAccessDenied))
	}

	er := new(e.Http)
	res := new(model.User)

	_, _, err = ctx.HTTP.
		GET("/user").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+*token).
		Request(res, er)
	if err != nil {
		return err
	}

	if er != nil {
		if er.Code == 401 {
			return errors.New(m.Message(er.Status))
		}

		if er.Code == 500 {
			return errors.New(m.Message(er.Status))
		}
	}

	fmt.Println(fmt.Sprintf("Username: %s\n"+
		"E-mail: %s\nBalance: %.0f\n"+
		"Organization: %t\nCreated: %s\n"+
		"Updated: %s", res.Username, res.Email,
		res.Balance, res.Organization, res.Created.String()[:10], res.Updated.String()[:10]))

	return nil
}
