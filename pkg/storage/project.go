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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"time"
)

const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	IProject
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get project by name for user
func (s *ProjectStorage) GetByName(username, name string) (*types.Project, error) {
	var (
		project = new(types.Project)
		key     = fmt.Sprintf("%s/%s/%s/info", ProjectTable, username, name)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, key, project); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}

// List project by username
func (s *ProjectStorage) ListByUser(username string) (*types.ProjectList, error) {
	var (
		projectList = new(types.ProjectList)
		key         = fmt.Sprintf("%s/%s", ProjectTable, username)
		filter      = `\b(.+)\/info\b`
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.List(ctx, key, filter, projectList); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return projectList, nil
}

// Insert new project into storage  for user
func (s *ProjectStorage) Insert(username, name, description string) (*types.Project, error) {
	var (
		project = new(types.Project)
		keyInfo = fmt.Sprintf("%s/%s/%s/info", ProjectTable, username, name)
	)

	project.Name = name
	project.User = username
	project.Description = description
	project.Updated = time.Now()
	project.Created = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return project, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Create(ctx, keyInfo, project, nil, 0); err != nil {
		return project, err
	}

	return project, nil
}

// Update project model
func (s *ProjectStorage) Update(project *types.Project) (*types.Project, error) {
	return nil, nil
}

// Remove project model
func (s *ProjectStorage) Remove(username, id string) error {
	var (
		key = fmt.Sprintf("%s/%s/%s/info", ProjectTable, username, id)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Delete(ctx, key, nil); err != nil {
		return err
	}

	return nil
}

func newProjectStorage(config store.Config) *ProjectStorage {
	s := new(ProjectStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
