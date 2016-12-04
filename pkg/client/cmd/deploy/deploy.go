package deploy

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

type createS struct {
	Project string `json:"project"`
	Target  string `json:"target"`
}

func DeployTargetCmd(target string) {

	var ctx = context.Get()

	err := DeployTarget(target)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func DeployTarget(target string) error {

	var (
		err     error
		ctx     = context.Get()
		project = new(model.Project)
		er      = new(e.Http)
		res     = new(struct{})
	)

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return errors.New(err.Error())
	}

	if project.ID == "" {
		return errors.New("Project didn't select")
	}

	_, _, err = ctx.HTTP.
		POST("/deploy").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		BodyJSON(createS{project.ID, target}).
		Request(res, er)
	if err != nil {
		return err
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	return nil
}
