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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

type Storage struct {
	*UserStorage
	*VendorStorage
	*ProjectStorage
	*ServiceStorage
	*ImageStorage
	*BuildStorage
	*HookStorage
	*VolumeStorage
	*ActivityStorage
}

func (s *Storage) User() IUser {
	if s == nil {
		return nil
	}
	return s.UserStorage
}

func (s *Storage) Vendor() IVendor {
	if s == nil {
		return nil
	}
	return s.VendorStorage
}

func (s *Storage) Project() IProject {
	if s == nil {
		return nil
	}
	return s.ProjectStorage
}

func (s *Storage) Service() IService {
	if s == nil {
		return nil
	}
	return s.ServiceStorage
}

func (s *Storage) Image() IImage {
	if s == nil {
		return nil
	}
	return s.ImageStorage
}

func (s *Storage) Build() IBuild {
	if s == nil {
		return nil
	}
	return s.BuildStorage
}

func (s *Storage) Hook() IHook {
	if s == nil {
		return nil
	}
	return s.HookStorage
}

func (s *Storage) Volume() IVolume {
	if s == nil {
		return nil
	}
	return s.VolumeStorage
}

func (s *Storage) Activity() IActivity {
	if s == nil {
		return nil
	}
	return s.ActivityStorage
}

func Get(config store.Config) (*Storage, error) {

	var (
		store  = new(Storage)
		helper = new(Helper)
	)

	store.UserStorage = NewUserStorage(config, helper)
	store.VendorStorage = NewVendorStorage(config, helper)
	store.ProjectStorage = NewProjectStorage(config, helper)
	store.ServiceStorage = NewServiceStorage(config, helper)
	store.ImageStorage = NewImageStorage(config, helper)
	store.BuildStorage = NewBuildStorage(config, helper)
	store.HookStorage = NewHookStorage(config, helper)
	store.VolumeStorage = NewVolumeStorage(config, helper)
	store.ActivityStorage = NewActivityStorage(config, helper)

	return store, nil
}

func New(c store.Config) (store.IStore, store.DestroyFunc, error) {
	return createEtcd3Storage(c)
}
