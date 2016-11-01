package rethinkdb

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Interface interface {
	GetByID(string) (*model.Build, *error)
}

// Service Service type for interface in interfaces folder
type BuildService struct{}

func (BuildService) GetByID(uuid string) (*model.Build, *e.Err) {

	var err error
	var build = new(model.Build)
	ctx := context.Get()

	res, err := r.Table("builds").Get(uuid).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.NotFound(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}
