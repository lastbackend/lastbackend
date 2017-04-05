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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"time"
)

const projectStorage = "projects"

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	IProject
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get project by name
func (s *ProjectStorage) GetByID(ctx context.Context, id string) (*types.Project, error) {

	project := new(types.Project)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, projectStorage, id, "meta")
	meta := new(types.ProjectMeta)
	if err := client.Get(ctx, key, meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	project.ID = meta.ID
	project.Name = meta.Name
	project.Description = meta.Description
	project.Labels = meta.Labels
	project.Created = meta.Created
	project.Updated = meta.Updated

	return project, nil
}

// Get project by name
func (s *ProjectStorage) GetByName(ctx context.Context, name string) (*types.Project, error) {

	var id string

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, "helper", projectStorage, name)
	if err := client.Get(ctx, key, &id); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// List projects
func (s *ProjectStorage) List(ctx context.Context) (*types.ProjectList, error) {

	const filter = `\b(.+)\/meta\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, projectStorage)
	metaList := []types.Meta{}
	if err := client.List(ctx, key, filter, &metaList); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if metaList == nil {
		return nil, nil
	}

	projectList := new(types.ProjectList)
	for _, meta := range metaList {
		project := types.Project{}
		project.ID = meta.ID
		project.Name = meta.Name
		project.Description = meta.Description
		project.Labels = meta.Labels
		project.Created = meta.Created
		project.Updated = meta.Updated

		*projectList = append(*projectList, project)
	}

	return projectList, nil
}

// Insert new project into storage
func (s *ProjectStorage) Insert(ctx context.Context, name, description string) (*types.Project, error) {
	var (
		id      = generator.GetUUIDV4()
		project = new(types.Project)
	)

	project.ID = id
	project.Name = name
	project.Description = description
	project.Labels = map[string]string{"tier": "ptoject"}
	project.Updated = time.Now()
	project.Created = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, name)
	if err := tx.Create(keyHelper, &project.ID, 0); err != nil {
		return nil, err
	}

	meta := new(types.ProjectMeta)
	meta.ID = id
	meta.Name = name
	meta.Description = description
	meta.Labels = map[string]string{"tier": "ptoject"}
	meta.Updated = time.Now()
	meta.Created = time.Now()

	keyMeta := s.util.Key(ctx, projectStorage, id, "meta")
	if err := tx.Create(keyMeta, meta, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return project, nil
}

// Update project model
func (s *ProjectStorage) Update(ctx context.Context, project *types.Project) (*types.Project, error) {

	project.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return project, err
	}
	defer destroy()

	meta := new(types.ProjectMeta)
	meta.ID = project.ID
	meta.Name = project.Name
	meta.Description = project.Description
	meta.Labels = project.Labels
	meta.Updated = time.Now()

	key := s.util.Key(ctx, projectStorage, project.ID, "meta")
	if err := client.Update(ctx, key, meta, nil, 0); err != nil {
		return project, err
	}

	return project, nil
}

// Remove project model
func (s *ProjectStorage) Remove(ctx context.Context, id string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, id, "meta")
	meta := new(types.ProjectMeta)
	if err := client.Get(ctx, keyMeta, meta); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, meta.Name)
	tx.Delete(keyHelper)

	key := s.util.Key(ctx, projectStorage, id)
	tx.Delete(key)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewProjectStorage(config store.Config, util IUtil) *ProjectStorage {
	s := new(ProjectStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
