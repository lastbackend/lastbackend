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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const BuildTable string = "builds"

// Service Build type for interface in interfaces folder
type BuildStorage struct {
	IBuild
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get build model by id
func (s *BuildStorage) GetByID(user, id string) (*types.Build, error) {
	return nil, nil
}

// Get builds by image
func (s *BuildStorage) ListByImage(user, id string) (*types.BuildList, error) {
	return nil, nil
}

// Insert new build into storage
func (s *BuildStorage) Insert(build *types.Build) (*types.Build, error) {
	return nil, nil
}

func newBuildStorage(config store.Config) *BuildStorage {
	s := new(BuildStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
