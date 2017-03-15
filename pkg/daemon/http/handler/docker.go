package handler

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/docker"
	"net/http"
)

func DockerRepositorySearchH(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		ctx    = c.Get()
		params = r.URL.Query()
		name   = params.Get("name")
	)

	ctx.Log.Debug("Search docker repository handler")

	repoListModel, err := docker.GetRepository(name)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := repoListModel.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}

func DockerRepositoryTagListH(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		ctx    = c.Get()
		params = r.URL.Query()
		owner  = params.Get("owner")
		name   = params.Get("name")
	)

	ctx.Log.Debug("List docker repository tags handler")

	tagListModel, err := docker.ListTag(owner, name)
	if err != nil {
		ctx.Log.Error(err)
		e.HTTP.InternalServerError(w)
		return
	}

	response, err := tagListModel.ToJson()
	if err != nil {
		ctx.Log.Error("Error: convert struct to json", err.Error())
		e.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}
