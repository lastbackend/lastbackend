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
	cfg "github.com/lastbackend/lastbackend/pkg/daemon/config"

)

type Storage struct {
	*UserStorage
	*ProjectStorage
	*ServiceStorage
	*ImageStorage
	*BuildStorage
	*HookStorage
	*VolumeStorage
	*ActivityStorage
}

func (s *Storage) User() storage.IUser {
	if s == nil {
		return nil
	}
	return s.UserStorage
}

func (s *Storage) Project() storage.IProject {
	if s == nil {
		return nil
	}
	return s.ProjectStorage
}

func (s *Storage) Service() storage.IService {
	if s == nil {
		return nil
	}
	return s.ServiceStorage
}

func (s *Storage) Image() storage.IImage {
	if s == nil {
		return nil
	}
	return s.ImageStorage
}

func (s *Storage) Build() storage.IBuild {
	if s == nil {
		return nil
	}
	return s.BuildStorage
}

func (s *Storage) Hook() storage.IHook {
	if s == nil {
		return nil
	}
	return s.HookStorage
}

func (s *Storage) Volume() storage.IVolume {
	if s == nil {
		return nil
	}
	return s.VolumeStorage
}

func (s *Storage) Activity() storage.IActivity {
	if s == nil {
		return nil
	}
	return s.ActivityStorage
}

func Get() (*Storage, error) {

	var (
		store  = new(Storage)
		config = cfg.GetEtcdDB()
	)

	store.UserStorage = newUserStorage(config)
	store.ProjectStorage = newProjectStorage(config)
	store.ServiceStorage = newServiceStorage(config)
	store.ImageStorage = newImageStorage(config)
	store.BuildStorage = newBuildStorage(config)
	store.HookStorage = newHookStorage(config)
	store.VolumeStorage = newVolumeStorage(config)
	store.ActivityStorage = newActivityStorage(config)

	return store, nil
}
