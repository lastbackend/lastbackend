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
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

const logLevel = 5

type Storage struct {
	*VendorStorage
	*NamespaceStorage
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

func Get(config store.Config, log logger.ILogger) (*Storage, error) {
	s := new(Storage)
	s.VendorStorage = newVendorStorage(config, log)
	s.NamespaceStorage = newNamespaceStorage(config, log)
	s.ServiceStorage = newServiceStorage(config, log)
	s.ImageStorage = newImageStorage(config, log)
	s.BuildStorage = newBuildStorage(config, log)
	s.HookStorage = newHookStorage(config, log)
	s.VolumeStorage = newVolumeStorage(config, log)
	s.ActivityStorage = newActivityStorage(config, log)
	s.NodeStorage = newNodeStorage(config, log)
	s.PodStorage = newPodStorage(config, log)
	s.SystemStorage = newSystemStorage(config, log)
	s.EndpointStorage = newEndpointStorage(config, log)
	return s, nil
}

func New(c store.Config, log logger.ILogger) (store.IStore, store.DestroyFunc, error) {
	return createEtcd3Storage(c, log)
}
