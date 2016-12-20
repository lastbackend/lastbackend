package service

import (
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/libs/editor"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"strings"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

func UpdateCmd(name string) {

	var ctx = context.Get()

	serviceModel, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	var config = model.ServiceConfig{
		Replicas:   serviceModel.Detail.Spec.Replicas,
		Command:    []string{"111", "222"},
		Args:       []string{},
		WorkingDir: "",
		Ports:      []model.Port{},
		Env:        []model.EnvVar{},
	}

	buf, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	buf1, err := yaml.Marshal(serviceModel.Detail.Spec)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info(string(buf1))

	res, err := editor.Run(strings.NewReader(string(buf)))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = res.ToYAML(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = Update(name, config)
	//if err != nil {
	//	ctx.Log.Error(err)
	//	return
	//}

	ctx.Log.Info("Successful")
}

func Update(name string, config model.ServiceConfig) error {

	var (
		err     error
		ctx     = context.Get()
		er      = new(e.Http)
		project = new(model.Project)
		res     = new(model.Project)
	)

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return errors.New(err.Error())
	}

	if project.ID == "" {
		return errors.New("Project didn't select")
	}

	_, _, err = ctx.HTTP.
		PUT("/project/"+project.ID+"/service/"+name).
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
