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

package routes

import (
	"github.com/lastbackend/lastbackend/pkg/api/repo"
	"github.com/lastbackend/lastbackend/pkg/api/repo/routes/request"
	"github.com/lastbackend/lastbackend/pkg/api/repo/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const logLevel = 2

func RepoListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
	)

	log.V(logLevel).Debug("Handler: Repo: list repo")

	ns := repo.New(r.Context())
	items, err := ns.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: find repo list err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewRepoList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: convert struct to json err: ", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: write response err: %s", err.Error())
		return
	}
}

func RepoInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		name = utils.Vars(r)["repo"]
	)

	log.V(logLevel).Debugf("Handler: Repo: get repo `%s`", name)

	ns := repo.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: find repo by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Repo: repo `%s` not found", name)
		errors.New("repo").NotFound().Http(w)
		return
	}

	response, err := v1.NewRepo(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: write response err: %s", err.Error())
		return
	}
}

func RepoCreateH(w http.ResponseWriter, r *http.Request) {

	log.V(logLevel).Debug("Handler: Repo: create repo")

	// request body struct
	rq := new(request.RequestRepoCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := repo.New(r.Context())
	item, err := ns.Get(rq.Name)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Handler: Repo: check exists by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item != nil {
		log.Warnf("Handler: Repo: name `%s` not unique", rq.Name)
		errors.New("repo").NotUnique("name").Http(w)
		return
	}

	n, err := ns.Create(rq)
	response, err := v1.NewRepo(n).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: write response err: %s", err.Error())
		return
	}
}

func RepoUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err    error
		params = utils.Vars(r)
		name   = params["repo"]
	)

	log.V(logLevel).Debugf("Handler: Repo: update repo `%s`", name)

	// request body struct
	rq := new(request.RequestRepoUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: validation incoming data", err)
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := repo.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: check exists by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Repo: repo name `%s` not found", name)
		errors.New("repo").NotFound().Http(w)
		return
	}

	item, err = ns.Update(item, rq)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: update repo `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
	}

	response, err := v1.NewRepo(item).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: write response err: %s", err.Error())
		return
	}
}

func RepoRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		err  error
		name = utils.Vars(r)["repo"]
	)

	log.V(logLevel).Debugf("Handler: Repo: remove repo %s", name)

	ns := repo.New(r.Context())
	item, err := ns.Get(name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: find repo by name `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.Warnf("Handler: Repo: repo name `%s` not found", name)
		errors.New("repo").NotFound().Http(w)
		return
	}

	// TODO: remove all activity by repo name

	ns.Remove(item.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Repo: remove repo `%s` err: %s", name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: write response err: %s", err.Error())
		return
	}
}

func RepoActivityListH(w http.ResponseWriter, r *http.Request) {

	var (
		name = utils.Vars(r)["repo"]
	)

	log.V(logLevel).Debugf("Handler: Repo: list repo `%s` activity", name)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		log.V(logLevel).Errorf("Handler: Repo: write response err: %s", err.Error())
		return
	}
}
