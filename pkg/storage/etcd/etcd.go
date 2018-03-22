//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package etcd

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

const logLevel = 3

var c struct {
	store store.Store
	dfunc store.DestroyFunc
}

type Storage struct {
	context.Context
	context.CancelFunc

	*ClusterStorage
	*DeploymentStorage
	*TriggerStorage
	*NodeStorage
	*NamespaceStorage
	*PodStorage
	*ServiceStorage
	*RouteStorage
	*VolumeStorage
	*SystemStorage
}

func (s *Storage) Cluster() storage.Cluster {
	if s == nil {
		return nil
	}
	return s.ClusterStorage
}

func (s *Storage) Deployment() storage.Deployment {
	if s == nil {
		return nil
	}
	return s.DeploymentStorage
}

func (s *Storage) Trigger() storage.Trigger {
	if s == nil {
		return nil
	}
	return s.TriggerStorage
}

func (s *Storage) Node() storage.Node {
	if s == nil {
		return nil
	}
	return s.NodeStorage
}

func (s *Storage) Namespace() storage.Namespace {
	if s == nil {
		return nil
	}
	return s.NamespaceStorage
}

func (s *Storage) Route() storage.Route {
	if s == nil {
		return nil
	}
	return s.RouteStorage
}

func (s *Storage) Pod() storage.Pod {
	if s == nil {
		return nil
	}
	return s.PodStorage
}

func (s *Storage) Service() storage.Service {
	if s == nil {
		return nil
	}
	return s.ServiceStorage
}

func (s *Storage) Volume() storage.Volume {
	if s == nil {
		return nil
	}
	return s.VolumeStorage
}

func (s *Storage) System() storage.System {
	if s == nil {
		return nil
	}
	return s.SystemStorage
}

func keyCreate(args ...string) string {
	return strings.Join([]string(args), "/")
}

func keyDirCreate(args ...string) string {
	key := strings.Join([]string(args), "/")
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	return key
}

func New() (*Storage, error) {

	log.Debug("Etcd: define storage")

	var (
		err error
	)

	if c.store, c.dfunc, err = v3.GetClient(context.Background()); err != nil {
		log.Errorf("etcd: store initialize err: %s", err)
		return nil, err
	}

	s := new(Storage)

	s.ClusterStorage = newClusterStorage()
	s.NodeStorage = newNodeStorage()

	s.NamespaceStorage = newNamespaceStorage()
	s.ServiceStorage = newServiceStorage()
	s.DeploymentStorage = newDeploymentStorage()
	s.PodStorage = newPodStorage()

	s.TriggerStorage = newTriggerStorage()

	s.RouteStorage = newRouteStorage()
	s.SystemStorage = newSystemStorage()
	s.VolumeStorage = newVolumeStorage()

	return s, nil
}

func getClient(ctx context.Context) (store.Store, store.DestroyFunc, error) {
	return c.store, c.dfunc, nil
}
