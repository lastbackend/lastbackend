//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

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
