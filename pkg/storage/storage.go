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

const debugLevel = 5

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
	*NodeStorage
	*PodStorage
	*SystemStorage
	*EndpointStorage
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

	var store = new(Storage)

	if _util == nil {
		// Set default util helpers
		_util = new(util)
	}

	store.VendorStorage = newVendorStorage(config, log, _util)
	store.NamespaceStorage = newNamespaceStorage(config, log, _util)
	store.ServiceStorage = newServiceStorage(config, log, _util)
	store.ImageStorage = newImageStorage(config, log, _util)
	store.BuildStorage = newBuildStorage(config, log, _util)
	store.HookStorage = newHookStorage(config, log, _util)
	store.VolumeStorage = newVolumeStorage(config, log, _util)
	store.ActivityStorage = newActivityStorage(config, log, _util)
	store.NodeStorage = newNodeStorage(config, log, _util)
	store.PodStorage = newPodStorage(config, log, _util)
	store.SystemStorage = newSystemStorage(config, log, _util)
	store.EndpointStorage = newEndpointStorage(config, log, _util)
	return store, nil
}

func SetLogger() {

}

func New(c store.Config, log logger.ILogger) (store.IStore, store.DestroyFunc, error) {
	return createEtcd3Storage(c, log)
}
