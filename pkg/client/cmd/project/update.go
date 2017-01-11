package project

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"strings"
	"time"
)

type updateS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func UpdateCmd(name, newProjectName, description string) {

	var ctx = context.Get()
	var choice string

	if description == "" {
		ctx.Log.Info("Description is empty, field will be cleared\n" +
			"Want to continue? [Y\\n]")

		for {
			fmt.Scan(&choice)

			switch strings.ToLower(choice) {
			case "y":
				break
			case "n":
				return
			default:
				ctx.Log.Error("Incorrect input. [Y\n]")
				continue
			}

			break
		}
	}

	err := Update(name, newProjectName, description)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	ctx.Log.Info("Successful")
}

func Update(name, newProjectName, description string) error {

	var (
		err error
		ctx = context.Get()
		er  = new(e.Http)
		res = new(model.Project)
	)

	_, _, err = ctx.HTTP.
		PUT("/project/"+name).
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		BodyJSON(updateS{newProjectName, description}).
		Request(&res, er)
	if err != nil {
		return err
	}

	if er.Code == 401 {
		return e.NotLoggedMessage
	}

	if er.Code != 0 {
		return errors.New(er.Message)
	}

	project, err := Current()
	if err != nil {
		return errors.New(err.Error())
	}

	if project != nil {
		if name == project.Name {
			project.Name = newProjectName
			project.Description = description
			project.Updated = time.Now()

			err = ctx.Storage.Set("project", project)
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
