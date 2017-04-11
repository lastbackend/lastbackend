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
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
)

// Util helpers
var _util IUtil

type Storage struct {
	*VendorStorage
	*NamespaceStorage
	*ServiceStorage
	*ImageStorage
	*BuildStorage
	*HookStorage
	*VolumeStorage
	*ActivityStorage
}

func SetUtil(u IUtil) {
	_util = u
}

func (s *Storage) Vendor() IVendor {
	if s == nil {
		return nil
	}
	return s.VendorStorage
}

func (s *Storage) Namespace() INamespace {
	if s == nil {
		return nil
	}
	return s.NamespaceStorage
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

	var store = new(Storage)

	if _util == nil {
		// Set default util helpers
		_util = new(util)
	}

	store.VendorStorage = newVendorStorage(config, _util)
	store.NamespaceStorage = newNamespaceStorage(config, _util)
	store.ServiceStorage = newServiceStorage(config, _util)
	store.ImageStorage = newImageStorage(config, _util)
	store.BuildStorage = newBuildStorage(config, _util)
	store.HookStorage = newHookStorage(config, _util)
	store.VolumeStorage = newVolumeStorage(config, _util)
	store.ActivityStorage = newActivityStorage(config, _util)

	return store, nil
}

func New(c store.Config) (store.IStore, store.DestroyFunc, error) {
	return createEtcd3Storage(c)
}
