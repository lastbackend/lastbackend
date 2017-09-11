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
	"strings"
)

const logLevel = 5

type Storage struct {
	*VendorStorage
	*AppStorage
	*RepoStorage
	*ServiceStorage
	*ImageStorage
	*BuildStorage
	*HookStorage
	*VolumeStorage
	*ActivityStorage
	*NodeStorage
	*PodStorage
	*SystemStorage
	*EndpointStorage
}

func (s *Storage) Vendor() IVendor {
	if s == nil {
		return nil
	}
	return s.VendorStorage
}

func (s *Storage) App() IApp {
	if s == nil {
		return nil
	}
	return s.AppStorage
}

func (s *Storage) Repo() IRepo {
	if s == nil {
		return nil
	}
	return s.RepoStorage
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

func (s *Storage) Node() INode {
	if s == nil {
		return nil
	}
	return s.NodeStorage
}

func (s *Storage) Pod() IPod {
	if s == nil {
		return nil
	}
	return s.PodStorage
}

func (s *Storage) System() ISystem {
	if s == nil {
		return nil
	}
	return s.SystemStorage
}

func (s *Storage) Endpoint() IEndpoint {
	if s == nil {
		return nil
	}
	return s.EndpointStorage
}

func keyCreate(args ...string) string {
	return strings.Join([]string(args), "/")
}

func Get(config store.Config) (*Storage, error) {
	s := new(Storage)
	s.VendorStorage = newVendorStorage(config)
	s.AppStorage = newAppStorage(config)
	s.RepoStorage = newRepoStorage(config)
	s.ServiceStorage = newServiceStorage(config)
	s.ImageStorage = newImageStorage(config)
	s.BuildStorage = newBuildStorage(config)
	s.HookStorage = newHookStorage(config)
	s.VolumeStorage = newVolumeStorage(config)
	s.ActivityStorage = newActivityStorage(config)
	s.NodeStorage = newNodeStorage(config)
	s.PodStorage = newPodStorage(config)
	s.SystemStorage = newSystemStorage(config)
	s.EndpointStorage = newEndpointStorage(config)
	return s, nil
}

func New(c store.Config) (store.IStore, store.DestroyFunc, error) {
	return createEtcd3Storage(c)
}
