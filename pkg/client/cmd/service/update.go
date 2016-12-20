package service

import (
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/libs/editor"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"k8s.io/client-go/1.5/pkg/util/json"
	"strings"
)

type updateS struct {
	Name     string `json:"name"`
	Replicas int32  `json:"replicas"`
}

func UpdateCmd(name string) {

	var ctx = context.Get()

	serviceModel, _, err := Inspect(name)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	var config = model.Config{
		Replicas:   serviceModel.Spec.Spec.Replicas,
		Command:    []string{},
		Args:       []string{},
		WorkingDir: "",
		Ports:      []model.PortConfig{},
		Env:        []model.EnvVarConfig{},
		Volumes:    []model.VolumeConfig{},
	}

	buf, err := json.Marshal(config)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	res, err := editor.Run(strings.NewReader(string(buf)))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Lines())

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

func Update(name string, config interface{}) error {

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
