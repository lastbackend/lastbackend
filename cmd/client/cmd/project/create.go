package project

import (
	"github.com/lastbackend/lastbackend/cmd/client/context"
	"github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type CreateTemplate struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

type ProjectCreateS struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func Create(name, description string) {

	var (
		ctx = context.Get()
	)

	er := errors.Http{}
	res := model.Project{}

	ctx.HTTP.
		POST("/project").
		AddHeader("Content-Type", "application/json").
		BodyJSON(ProjectCreateS{name, description}).
		Request(&res, &er)

	ctx.Log.Info(res.ID)

}
