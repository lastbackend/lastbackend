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

const ActivityTable string = "activities"

// Activity Service type for interface in interfaces folder
type ActivityStorage struct {
	storage.IActivity
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ActivityStorage) Insert(activity *model.Activity) (*model.Activity, error) {
	return nil, nil
}

func (s *ActivityStorage) ListProjectActivity(user, project string) (*model.ActivityList, error) {
	return nil, nil
}

func (s *ActivityStorage) ListServiceActivity(user, service string) (*model.ActivityList, error) {
	return nil, nil
}

func (s *ActivityStorage) RemoveByProject(user, project string) error {
	return nil
}

func (s *ActivityStorage) RemoveByService(user, service string) error {
	return nil
}

func newActivityStorage(config store.Config) *ActivityStorage {
	s := new(ActivityStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
