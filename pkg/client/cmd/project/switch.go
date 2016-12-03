package project

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func SwitchCmd(name string) error {
	var ctx = context.Get()
	token := struct {
		Token string `json:"token"`
	}{}
	ctx.Storage.Get("session", &token)

	//var project = new(model.Project)

	er := new(e.Http)
	res := new(model.Project)

	_, _, err := ctx.HTTP.
		GET("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+token.Token).
		Request(&res, er)

	if err != nil {
		return errors.New(err.Error())
	}

	err = ctx.Storage.Set("project", res)
	if err != nil {
		return errors.New(err.Error())
	}
	//err = ctx.Storage.Get(name, &project)

	fmt.Println("The project \"", name, "\" was selected as the current")

	return nil
}
