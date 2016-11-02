package handler

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"net/http"
	"time"
)

func ProjectListH(w http.ResponseWriter, _ *http.Request) {
	var ctx = context.Get()
	ctx.Log.Info("get projects list")

	projects, err := ctx.Storage.Project().GetByUser("")
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(projects)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(buf)
}

func ProjectCreateH(w http.ResponseWriter, _ *http.Request) {
	var ctx = context.Get()
	ctx.Log.Info("create project")

	p := new(model.Project)
	p.Name = "test"
	p.Namespace = "test"
	p.Description = "test"
	p.User = ""

	p.Created = time.Now()
	p.Updated = time.Now()

	project, err := ctx.Storage.Project().Insert(p)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	buf, er := json.Marshal(project)
	if er != nil {
		ctx.Log.Error(er.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(buf)
}
