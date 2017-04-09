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
	"github.com/satori/go.uuid"
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
func (s *ProjectStorage) GetByID(ctx context.Context, id uuid.UUID) (*types.Project, error) {

	project := new(types.Project)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, projectStorage, id.String(), "meta")
	meta := types.Meta{}
	if err := client.Get(ctx, key, meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	project.Meta = meta

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

	return s.GetByID(ctx, uuid.FromStringOrNil(id))
}

// List projects
func (s *ProjectStorage) List(ctx context.Context) (*types.ProjectList, error) {

	const filter = `\b(.+)projects\/[a-z0-9-]{36}\/meta\b`

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
		project.Meta = meta
		*projectList = append(*projectList, project)
	}

	return projectList, nil
}

// Insert new project into storage
func (s *ProjectStorage) Insert(ctx context.Context, name, description string) (*types.Project, error) {
	var (
		id      = uuid.NewV4()
		project = new(types.Project)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, name)
	if err := tx.Create(keyHelper, id.String(), 0); err != nil {
		return nil, err
	}

	meta := new(types.Meta)
	meta.ID = id
	meta.Name = name
	meta.Description = description
	meta.Labels = map[string]string{"tier": "ptoject"}
	meta.Updated = time.Now()
	meta.Created = time.Now()

	keyMeta := s.util.Key(ctx, projectStorage, id.String(), "meta")
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

	project.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return project, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, project.Meta.ID.String(), "meta")
	pmeta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, pmeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	meta := types.Meta{}
	meta = project.Meta
	meta.Updated = time.Now()

	tx := client.Begin(ctx)

	if pmeta.Name != project.Meta.Name {
		keyHelper1 := s.util.Key(ctx, "helper", projectStorage, pmeta.Name)
		tx.Delete(keyHelper1)

		keyHelper2 := s.util.Key(ctx, "helper", projectStorage, project.Meta.Name)
		if err := tx.Create(keyHelper2, project.Meta.ID.String(), 0); err != nil {
			return project, err
		}
	}

	keyMeta = s.util.Key(ctx, projectStorage, project.Meta.ID.String(), "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return project, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return project, nil
}

// Remove project model
func (s *ProjectStorage) Remove(ctx context.Context, id uuid.UUID) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, id.String(), "meta")
	meta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, meta); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, meta.Name)
	tx.Delete(keyHelper)

	key := s.util.Key(ctx, projectStorage, id.String())
	tx.DeleteDir(key)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func newProjectStorage(config store.Config, util IUtil) *ProjectStorage {
	s := new(ProjectStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
