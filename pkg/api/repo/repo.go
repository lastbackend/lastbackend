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

package repo

import (
	"context"
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	ins "github.com/lastbackend/lastbackend/pkg/api/repo/interfaces"
	"github.com/lastbackend/lastbackend/pkg/api/repo/routes/request"
	u "github.com/lastbackend/lastbackend/pkg/api/repo/utils"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

var utils ins.IUtil = new(u.Util)

type repo struct {
	Context context.Context
}

func New(ctx context.Context) *repo {
	return &repo{ctx}
}

func SetUtils(u ins.IUtil) {
	utils = u
}

func (r *repo) List() (types.RepoList, error) {
	var (
		storage = ctx.Get().GetStorage()
		list    = types.RepoList{}
	)

	log.V(logLevel).Debug("Repo: list repo")

	apps, err := storage.Repo().List(r.Context)
	if err != nil {
		log.V(logLevel).Error("Repo: list repo err: %s", err.Error())
		return list, err
	}

	log.V(logLevel).Debugf("Repo: list repo result: %d", len(apps))

	for _, item := range apps {
		var repo = item
		list = append(list, repo)
	}

	return list, nil
}

func (r *repo) Get(name string) (*types.Repo, error) {
	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Repo: get repo %s", name)

	repo, err := storage.Repo().GetByName(r.Context, name)
	if err != nil {
		if err.Error() == store.ErrKeyNotFound {
			log.V(logLevel).Warnf("Repo: repo by name `%s` not found", name)
			return nil, nil
		}
		log.V(logLevel).Errorf("Repo: get repo by name `%s` err: %s", name, err.Error())
		return nil, err
	}

	return repo, nil
}

func (r *repo) Create(rq *request.RequestRepoCreateS) (*types.Repo, error) {

	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Repo: create repo %#v", rq)

	var repo = types.Repo{}
	repo.Meta.SetDefault()
	repo.Meta.Name = utils.NameCreate(r.Context, rq.Name)

	if err = storage.Repo().Insert(r.Context, &repo); err != nil {
		log.V(logLevel).Errorf("Repo: insert repo err: %s", err.Error())
		return nil, err
	}

	return &repo, nil
}

func (r *repo) Update(repo *types.Repo, rq *request.RequestRepoUpdateS) (*types.Repo, error) {
	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Repo: update repo %#v", repo)

	repo.Meta.Technology = *rq.Technology
	repo.Meta.Description = *rq.Description

	if err = storage.Repo().Update(r.Context, repo); err != nil {
		log.V(logLevel).Errorf("Repo: update repo err: %s", err.Error())
		return repo, err
	}

	return repo, nil
}

func (r *repo) Remove(name string) error {
	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Repo: remove repo %s", name)

	err = storage.Repo().Remove(r.Context, name)
	if err != nil {
		log.V(logLevel).Errorf("Repo: remove repo err: %s", err.Error())
		return err
	}

	return nil
}