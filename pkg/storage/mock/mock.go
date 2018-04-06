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

package mock

import (
	"context"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const logLevel = 5

type Storage struct {
	context.Context
	context.CancelFunc

	*ClusterStorage
	*DeploymentStorage
	*TriggerStorage
	*NodeStorage
	*IngressStorage
	*NamespaceStorage
	*PodStorage
	*ServiceStorage
	*RouteStorage
	*VolumeStorage
	*SecretStorage
	*SystemStorage
}

func (s *Storage) Cluster() storage.Cluster {
	return s.ClusterStorage
}

func (s *Storage) Deployment() storage.Deployment {
	return s.DeploymentStorage
}

func (s *Storage) Trigger() storage.Trigger {
	return s.TriggerStorage
}

func (s *Storage) Node() storage.Node {
	return s.NodeStorage
}

func (s *Storage) Ingress() storage.Ingress {
	return s.IngressStorage
}

func (s *Storage) Namespace() storage.Namespace {
	return s.NamespaceStorage
}

func (s *Storage) Route() storage.Route {
	return s.RouteStorage
}

func (s *Storage) Pod() storage.Pod {
	return s.PodStorage
}

func (s *Storage) Service() storage.Service {
	return s.ServiceStorage
}

func (s *Storage) Volume() storage.Volume {
	return s.VolumeStorage
}

func (s *Storage) Secret() storage.Secret {
	return s.SecretStorage
}

func (s *Storage) System() storage.System {
	return s.SystemStorage
}

func keyCreate(args ...string) string {
	return strings.Join([]string(args), "/")
}

func New() (*Storage, error) {

	log.Debug("Etcd: define mock storage")

	s := new(Storage)

	s.ClusterStorage = newClusterStorage()
	s.NodeStorage = newNodeStorage()
	s.IngressStorage = newIngressStorage()

	s.NamespaceStorage = newNamespaceStorage()
	s.ServiceStorage = newServiceStorage()
	s.DeploymentStorage = newDeploymentStorage()
	s.PodStorage = newPodStorage()

	s.TriggerStorage = newTriggerStorage()

	s.RouteStorage = newRouteStorage()
	s.SystemStorage = newSystemStorage()
	s.VolumeStorage = newVolumeStorage()
	s.SecretStorage = newSecretStorage()

	return s, nil
}

func getClient(_ context.Context) (store.Store, store.DestroyFunc, error) {

	log.V(logLevel).Debug("Etcd3: initialization storage")
	return nil, nil, nil
}
