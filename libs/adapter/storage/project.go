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
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	db "github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	storage.IProject
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ProjectStorage) GetByNameOrID(user, nameOrID string) (*model.Project, error) {
	return nil, nil
}

func (s *ProjectStorage) GetByName(user, name string) (*model.Project, error) {
	return nil, nil
}

func (s *ProjectStorage) ExistByName(userID, name string) (bool, error) {
	return false, nil
}

func (s *ProjectStorage) GetByID(user, id string) (*model.Project, error) {
	return nil, nil
}

func (s *ProjectStorage) ListByUser(id string) (*model.ProjectList, error) {
	return nil, nil
}

// Insert new project into storage
func (s *ProjectStorage) Insert(project *model.Project) (*model.Project, error) {
	return nil, nil
}

// Update build model
func (s *ProjectStorage) Update(project *model.Project) (*model.Project, error) {
	return nil, nil
}

// Remove build model
func (s *ProjectStorage) Remove(user, id string) error {
	return nil
}

func newProjectStorage(config store.Config) *ProjectStorage {
	s := new(ProjectStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
