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

type Namespace interface {
	Get(ctx context.Context, name string) (*types.Namespace, error)
	List(ctx context.Context) (map[string]*types.Namespace, error)
	Insert(ctx context.Context, namespace *types.Namespace) error
	Update(ctx context.Context, namespace *types.Namespace) error
	Remove(ctx context.Context, namespace *types.Namespace) error
	Clear(ctx context.Context) error
}

type Service interface {
	Get(ctx context.Context, namespace, name string) (*types.Service, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Service, error)
	SetStatus(ctx context.Context, service *types.Service) error
	SetSpec(ctx context.Context, service *types.Service) error
	Insert(ctx context.Context, service *types.Service) error
	Update(ctx context.Context, service *types.Service) error
	Remove(ctx context.Context, service *types.Service) error
	Watch(ctx context.Context, event chan *types.Event) error
	WatchSpec(ctx context.Context, service chan *types.Service) error
	WatchStatus(ctx context.Context, service chan *types.Service) error
	Clear(ctx context.Context) error
}

type Deployment interface {
	Get(ctx context.Context, namespace, service, name string) (*types.Deployment, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Deployment, error)
	ListByService(ctx context.Context, namespace, service string) (map[string]*types.Deployment, error)
	SetSpec(ctx context.Context, d *types.Deployment) error
	SetStatus(ctx context.Context, d *types.Deployment) error
	Insert(ctx context.Context, d *types.Deployment) error
	Update(ctx context.Context, d *types.Deployment) error
	Remove(ctx context.Context, d *types.Deployment) error
	Watch(ctx context.Context, deployment chan *types.Deployment) error
	WatchSpec(ctx context.Context, deployment chan *types.Deployment) error
	WatchStatus(ctx context.Context, deployment chan *types.Deployment) error
	Clear(ctx context.Context) error
}

type Pod interface {
	Get(ctx context.Context, namespace, service, deployment, name string) (*types.Pod, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Pod, error)
	ListByService(ctx context.Context, namespace, service string) (map[string]*types.Pod, error)
	ListByDeployment(ctx context.Context, namespace, service, deployment string) (map[string]*types.Pod, error)
	SetMeta(ctx context.Context, pod *types.Pod) error
	SetSpec(ctx context.Context, pod *types.Pod) error
	SetStatus(ctx context.Context, pod *types.Pod) error
	Insert(ctx context.Context, pod *types.Pod) error
	Update(ctx context.Context, pod *types.Pod) error
	Remove(ctx context.Context, pod *types.Pod) error
	Watch(ctx context.Context, pod chan *types.Pod) error
	WatchSpec(ctx context.Context, pod chan *types.Pod) error
	WatchStatus(ctx context.Context, pod chan *types.Pod) error
	Clear(ctx context.Context) error
}

type Trigger interface {
	Get(ctx context.Context, namespace, service, name string) (*types.Trigger, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Trigger, error)
	ListByService(ctx context.Context, namespace, service string) (map[string]*types.Trigger, error)
	SetStatus(ctx context.Context, trigger *types.Trigger) error
	SetSpec(ctx context.Context, trigger *types.Trigger) error
	Insert(ctx context.Context, trigger *types.Trigger) error
	Update(ctx context.Context, trigger *types.Trigger) error
	Remove(ctx context.Context, trigger *types.Trigger) error
	Watch(ctx context.Context, trigger chan *types.Trigger) error
	WatchSpec(ctx context.Context, trigger chan *types.Trigger) error
	WatchStatus(ctx context.Context, trigger chan *types.Trigger) error
	Clear(ctx context.Context) error
}

type Route interface {
	Get(ctx context.Context, namespace, name string) (*types.Route, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Route, error)
	ListSpec(ctx context.Context) (map[string]*types.RouteSpec, error)
	SetStatus(ctx context.Context, route *types.Route) error
	SetSpec(ctx context.Context, route *types.Route) error
	Insert(ctx context.Context, route *types.Route) error
	Update(ctx context.Context, route *types.Route) error
	Remove(ctx context.Context, route *types.Route) error
	Watch(ctx context.Context, route chan *types.Route) error
	WatchSpec(ctx context.Context, route chan *types.Route) error
	WatchSpecEvents(ctx context.Context, event chan *types.RouteSpecEvent) error
	WatchStatus(ctx context.Context, route chan *types.Route) error
	Clear(ctx context.Context) error
}

type Secret interface {
	Get(ctx context.Context, namespace, name string) (*types.Secret, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Secret, error)
	Insert(ctx context.Context, route *types.Secret) error
	Update(ctx context.Context, route *types.Secret) error
	Remove(ctx context.Context, route *types.Secret) error
	Clear(ctx context.Context) error
}

type Volume interface {
	Get(ctx context.Context, namespace, name string) (*types.Volume, error)
	ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Volume, error)
	SetStatus(ctx context.Context, volume *types.Volume) error
	SetSpec(ctx context.Context, volume *types.Volume) error
	Insert(ctx context.Context, volume *types.Volume) error
	Update(ctx context.Context, volume *types.Volume) error
	Remove(ctx context.Context, volume *types.Volume) error
	Watch(ctx context.Context, volume chan *types.Volume) error
	WatchSpec(ctx context.Context, volume chan *types.Volume) error
	WatchStatus(ctx context.Context, volume chan *types.Volume) error
	Clear(ctx context.Context) error
}

type Cluster interface {
	Insert(ctx context.Context, cluster *types.Cluster) error
	Get(ctx context.Context) (*types.Cluster, error)
	Update(ctx context.Context, cluster *types.Cluster) error
	Clear(ctx context.Context) error
}

type Node interface {
	List(ctx context.Context) (map[string]*types.Node, error)
	Get(ctx context.Context, name string) (*types.Node, error)
	GetSpec(ctx context.Context, node *types.Node) (*types.NodeSpec, error)
	Insert(ctx context.Context, node *types.Node) error
	Update(ctx context.Context, node *types.Node) error
	SetStatus(ctx context.Context, node *types.Node) error
	SetInfo(ctx context.Context, node *types.Node) error
	SetNetwork(ctx context.Context, node *types.Node) error
	SetOnline(ctx context.Context, node *types.Node) error
	SetOffline(ctx context.Context, node *types.Node) error
	InsertPod(ctx context.Context, node *types.Node, pod *types.Pod) error
	UpdatePod(ctx context.Context, node *types.Node, pod *types.Pod) error
	RemovePod(ctx context.Context, node *types.Node, pod *types.Pod) error
	InsertVolume(ctx context.Context, node *types.Node, volume *types.Volume) error
	RemoveVolume(ctx context.Context, node *types.Node, volume *types.Volume) error
	Remove(ctx context.Context, node *types.Node) error
	Watch(ctx context.Context, node chan *types.Node) error
	WatchStatus(ctx context.Context, event chan *types.NodeStatusEvent) error
	WatchPodSpec(ctx context.Context, event chan *types.PodSpecEvent) error
	WatchVolumeSpec(ctx context.Context, event chan *types.VolumeSpecEvent) error
	Clear(ctx context.Context) error
}

type Ingress interface {
	List(ctx context.Context) (map[string]*types.Ingress, error)
	Get(ctx context.Context, name string) (*types.Ingress, error)
	GetSpec(ctx context.Context, ingress *types.Ingress) (*types.IngressSpec, error)
	Insert(ctx context.Context, ingress *types.Ingress) error
	Update(ctx context.Context, ingress *types.Ingress) error
	SetStatus(ctx context.Context, ingress *types.Ingress) error
	Remove(ctx context.Context, ingress *types.Ingress) error
	Watch(ctx context.Context, ingress chan *types.Ingress) error
	WatchStatus(ctx context.Context, event chan *types.IngressStatusEvent) error
	Clear(ctx context.Context) error
}

type System interface {
	ProcessSet(ctx context.Context, process *types.Process) error
	Elect(ctx context.Context, process *types.Process) (bool, error)
	ElectUpdate(ctx context.Context, process *types.Process) error
	ElectWait(ctx context.Context, process *types.Process, lead chan bool) error
	Clear(ctx context.Context) error
}
