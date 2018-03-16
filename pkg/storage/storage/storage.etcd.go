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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type Hook interface {
	Get(ctx context.Context, id string) (*types.Hook, error)
	Insert(ctx context.Context, hook *types.Hook) error
	Remove(ctx context.Context, id string) error
}

type Namespace interface {
	Get(ctx context.Context, name string) (*types.Namespace, error)
	List(ctx context.Context) ([]*types.Namespace, error)
	Insert(ctx context.Context, namespace *types.Namespace) error
	Update(ctx context.Context, project *types.Namespace) error
	Remove(ctx context.Context, id string) error
}

type Service interface {
	CountByNamespace(ctx context.Context, namespace string) (int, error)
	Get(ctx context.Context, namespace, name string) (*types.Service, error)
	ListByNamespace(ctx context.Context, namespace string) ([]*types.Service, error)
	Insert(ctx context.Context, service *types.Service) error
	Update(ctx context.Context, service *types.Service) error
	UpdateSpec(ctx context.Context, service *types.Service) error
	Remove(ctx context.Context, service *types.Service) error
	RemoveByNamespace(ctx context.Context, namespace string) error

	Watch(ctx context.Context, service chan *types.Service) error
	SpecWatch(ctx context.Context, service chan *types.Service) error
}

type Deployment interface {
	Get(ctx context.Context, namespace, name string) (*types.Deployment, error)
	ListByNamespace(ctx context.Context, namespace string) ([]*types.Deployment, error)
	ListByService(ctx context.Context, namespace, service string) ([]*types.Deployment, error)
	SetState(ctx context.Context, d *types.Deployment) error
	Insert(ctx context.Context, d *types.Deployment) error
	Update(ctx context.Context, d *types.Deployment) error
	Remove(ctx context.Context, d *types.Deployment) error
	Watch(ctx context.Context, deployment chan *types.Deployment) error
	SpecWatch(ctx context.Context, service chan *types.Deployment) error
}

type Pod interface {
	Get(ctx context.Context, namespace, name string) (*types.Pod, error)
	ListByNamespace(ctx context.Context, namespace string) ([]*types.Pod, error)
	ListByService(ctx context.Context, namespace, service string) ([]*types.Pod, error)
	ListByDeployment(ctx context.Context, namespace, service, deployment string) ([]*types.Pod, error)
	SetState(ctx context.Context, pod *types.Pod) error
	Insert(ctx context.Context, pod *types.Pod) error
	Upsert(ctx context.Context, pod *types.Pod) error
	Update(ctx context.Context, pod *types.Pod) error
	Destroy(ctx context.Context, pod *types.Pod) error
	Remove(ctx context.Context, pod *types.Pod) error
	Watch(ctx context.Context, pod chan *types.Pod) error
}

type Volume interface {
	GetByToken(ctx context.Context, token string) (*types.Volume, error)
	ListByNamespace(ctx context.Context, namespace string) ([]*types.Volume, error)
	Insert(ctx context.Context, volume *types.Volume) error
	Remove(ctx context.Context, id string) error
}

type Cluster interface {
	Info(ctx context.Context) (*types.Cluster, error)
	Update(ctx context.Context, cluster *types.Cluster) error
}

type Node interface {
	List(ctx context.Context) ([]*types.Node, error)

	Get(ctx context.Context, name string) (*types.Node, error)
	Insert(ctx context.Context, node *types.Node) error

	Update(ctx context.Context, node *types.Node) error

	SetState(ctx context.Context, node *types.Node) error
	SetInfo(ctx context.Context, node *types.Node) error
	SetNetwork(ctx context.Context, node *types.Node) error

	SetAvailable(ctx context.Context, node *types.Node) error
	SetUnavailable(ctx context.Context, node *types.Node) error

	InsertPod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error
	UpdatePod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error
	RemovePod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error

	Remove(ctx context.Context, name string) error
	Watch(ctx context.Context, node chan *types.Node) error
}

type System interface {
	ProcessSet(ctx context.Context, process *types.Process) error

	Elect(ctx context.Context, process *types.Process) (bool, error)
	ElectUpdate(ctx context.Context, process *types.Process) error
	ElectWait(ctx context.Context, process *types.Process, lead chan bool) error
}

type Route interface {
	Get(ctx context.Context, namespace, name string) (*types.Route, error)
	ListByNamespace(ctx context.Context, namespace string) ([]*types.Route, error)
	ListByService(ctx context.Context, namespace, service string) ([]*types.Route, error)
	Insert(ctx context.Context, route *types.Route) error
	Update(ctx context.Context, route *types.Route) error
	Remove(ctx context.Context, route *types.Route) error
}

type Endpoint interface {
	Get(ctx context.Context, name string) ([]string, error)
	Upsert(ctx context.Context, name string, ips []string) error
	Remove(ctx context.Context, name string) error
	Watch(ctx context.Context, endpoint chan string) error
}
