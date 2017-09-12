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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/pkg/errors"
	"time"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const repoStorage = "repos"

// Repo Service type for interface in interfaces folder
type RepoStorage struct {
	IRepo
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get repo by name
func (s *RepoStorage) GetByName(ctx context.Context, name string) (*types.Repo, error) {

	log.V(logLevel).Debugf("Storage: Repo: get by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Repo: get repo err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Repo: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	repo := new(types.Repo)
	keyMeta := keyCreate(repoStorage, name, "meta")
	if err := client.Get(ctx, keyMeta, &repo.Meta); err != nil {
		log.V(logLevel).Errorf("Storage: Repo: get repo `%s` meta err: %s", name, err.Error())
		return nil, err
	}

	return repo, nil
}

// List projects
func (s *RepoStorage) List(ctx context.Context) ([]*types.Repo, error) {

	log.V(logLevel).Debug("Storage: Repo: get repo list")

	const filter = `\b(.+)` + repoStorage + `\/.+\/(meta)\b`

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Repo: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	repos := []*types.Repo{}
	keyRepos := keyCreate(repoStorage)
	if err := client.List(ctx, keyRepos, filter, &repos); err != nil {
		log.V(logLevel).Errorf("Storage: Repo: get repos list err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("Storage: Repo: get repo list result: %d", len(repos))

	return repos, nil
}

// Insert new repo into storage
func (s *RepoStorage) Insert(ctx context.Context, repo *types.Repo) error {

	log.V(logLevel).Debug("Storage: Repo: insert repo: %#v", repo)

	if repo == nil {
		err := errors.New("repo can not be nil")
		log.V(logLevel).Errorf("Storage: Repo: insert repo err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Repo: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(repoStorage, repo.Meta.Name, "meta")
	if err := client.Create(ctx, keyMeta, repo.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Repo: insert repo err: %s", err.Error())
		return err
	}

	return nil
}

// Update repo model
func (s *RepoStorage) Update(ctx context.Context, repo *types.Repo) error {

	log.V(logLevel).Debugf("Storage: Repo: update repo: %#v", repo)

	if repo == nil {
		err := errors.New("repo can not be nil")
		log.V(logLevel).Errorf("Storage: Repo: update repo err: %s", err.Error())
		return err
	}

	repo.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Repo: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	meta := types.RepoMeta{}
	meta = repo.Meta
	meta.Updated = time.Now()

	keyMeta := keyCreate(repoStorage, repo.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Repo: update repo meta err: %s", err.Error())
		return err
	}

	return nil
}

// Remove repo model
func (s *RepoStorage) Remove(ctx context.Context, name string) error {

	log.V(logLevel).Debugf("Storage: Repo: remove repo: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Repo: remove repo err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Repo: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyRepo := keyCreate(repoStorage, name)
	if err := client.DeleteDir(ctx, keyRepo); err != nil {
		log.V(logLevel).Errorf("Storage: Repo: remove repo `%s` err: %s", name, err.Error())
		return err
	}

	return nil
}

func newRepoStorage(config store.Config) *RepoStorage {
	s := new(RepoStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
