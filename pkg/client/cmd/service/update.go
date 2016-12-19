package service

import (
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/libs/editor"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"strings"
)

type updateS struct {
	Name     string `json:"name"`
	Replicas int32  `json:"replicas"`
}

func UpdateCmd(name string) {

	var ctx = context.Get()

	//service, err := Inspect(name)
	//if err != nil {
	//	return err
	//}

	var config interface{}

	input := strings.NewReader("-replicas: 1")
	res, err := editor.Run(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = res.ToYAML(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = Update(name, config)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Successful")
}

func Update(name string, config interface{}) error {

	var (
		err error
		ctx = context.Get()
		er  = new(e.Http)
		res = new(model.Project)
	)

	_, _, err = ctx.HTTP.
		PUT("/service/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		BodyJSON(config).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code == 401 {
		return errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	return nil
}
