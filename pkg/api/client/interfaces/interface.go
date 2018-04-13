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

package interfaces

import (
	"context"

	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"io"
)

type Cluster interface {
	Node(hostname string) Node

	Get(ctx context.Context) (*vv1.ClusterList, error)
	Update(ctx context.Context, opts *rv1.ClusterUpdateOptions) (*vv1.Cluster, error)
}

type Node interface {
	List(ctx context.Context) (*vv1.NodeList, error)
	Connect(ctx context.Context, opts *rv1.NodeConnectOptions) error
	Get(ctx context.Context) (*vv1.Node, error)
	GetSpec(ctx context.Context) (*vv1.NodeSpec, error)
	SetMeta(ctx context.Context, opts *rv1.NodeMetaOptions) (*vv1.Node, error)
	SetStatus(ctx context.Context, opts *rv1.NodeStatusOptions) error
	SetPodStatus(ctx context.Context, pod string, opts *rv1.NodePodStatusOptions) error
	SetVolumeStatus(ctx context.Context, volume string, opts *rv1.NodeVolumeStatusOptions) error
	SetRouteStatus(ctx context.Context, route string, opts *rv1.NodeRouteStatusOptions) error
	Remove(ctx context.Context, opts *rv1.NodeRemoveOptions) error
}

type Ingress interface {
	List(ctx context.Context) (*vv1.IngressList, error)
	Connect(ctx context.Context, opts *rv1.IngressConnectOptions) error
	Get(ctx context.Context) (*vv1.Ingress, error)
	GetSpec(ctx context.Context) (*vv1.IngressSpec, error)
	SetMeta(ctx context.Context, opts *rv1.IngressMetaOptions) (*vv1.Ingress, error)
	SetStatus(ctx context.Context, opts *rv1.IngressStatusOptions) error
	SetRouteStatus(ctx context.Context, route string, opts *rv1.IngressRouteStatusOptions) error
	Remove(ctx context.Context, opts *rv1.IngressRemoveOptions) error
	Logs(ctx context.Context, pod, container string, opts *rv1.IngressLogsOptions) (io.ReadCloser, error)
}

type Namespace interface {
	Service(name ...string) Service
	Secret(name ...string) Secret
	Volume(name ...string) Volume
	Route(name ...string) Route

	Create(ctx context.Context, opts *rv1.NamespaceCreateOptions) (*vv1.Namespace, error)
	List(ctx context.Context) (*vv1.NamespaceList, error)
	Get(ctx context.Context) (*vv1.Namespace, error)
	Update(ctx context.Context, opts *rv1.NamespaceUpdateOptions) (*vv1.Namespace, error)
	Remove(ctx context.Context, opts *rv1.NamespaceRemoveOptions) error
}

type Service interface {
	Deployment(name ...string) Deployment
	Trigger(name ...string) Trigger

	Create(ctx context.Context, opts *rv1.ServiceCreateOptions) (*vv1.Service, error)
	List(ctx context.Context) (*vv1.ServiceList, error)
	Get(ctx context.Context) (*vv1.Service, error)
	Update(ctx context.Context, opts *rv1.ServiceUpdateOptions) (*vv1.NamespaceList, error)
	Remove(ctx context.Context, opts *rv1.ServiceRemoveOptions) error
	Logs(ctx context.Context, opts *rv1.ServiceLogsOptions) (io.ReadCloser, error)
}

type Deployment interface {
	Pod(name ...string) Pod
	List(ctx context.Context) (*vv1.DeploymentList, error)
	Get(ctx context.Context) (*vv1.Deployment, error)
	Update(ctx context.Context, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error)
}

type Pod interface {
	List(ctx context.Context) (*vv1.PodList, error)
	Get(ctx context.Context) (*vv1.Pod, error)
	Logs(ctx context.Context, opts *rv1.PodLogsOptions) (io.ReadCloser, error)
}

type Events interface {
}

type Secret interface {
	Create(ctx context.Context, opts *rv1.SecretCreateOptions) (*vv1.Secret, error)
	List(ctx context.Context) (*vv1.SecretList, error)
	Update(ctx context.Context, opts *rv1.SecretUpdateOptions) (*vv1.Secret, error)
	Remove(ctx context.Context, opts *rv1.SecretRemoveOptions) error
}

type Route interface {
	Create(ctx context.Context, opts *rv1.RouteCreateOptions) (*vv1.Route, error)
	List(ctx context.Context) (*vv1.RouteList, error)
	Get(ctx context.Context) (*vv1.Route, error)
	Update(ctx context.Context, opts *rv1.RouteUpdateOptions) (*vv1.Route, error)
	Remove(ctx context.Context, opts *rv1.RouteRemoveOptions) error
}

type Trigger interface {
	Create(ctx context.Context, opts *rv1.TriggerCreateOptions) (*vv1.Trigger, error)
	List(ctx context.Context) (*vv1.TriggerList, error)
	Get(ctx context.Context) (*vv1.Trigger, error)
	Update(ctx context.Context, opts *rv1.TriggerUpdateOptions) (*vv1.Trigger, error)
	Remove(ctx context.Context, opts *rv1.TriggerRemoveOptions) error
}

type Volume interface {
	Create(ctx context.Context, opts *rv1.VolumeCreateOptions) (*vv1.Volume, error)
	List(ctx context.Context) (*vv1.VolumeList, error)
	Get(ctx context.Context) (*vv1.Volume, error)
	Update(ctx context.Context, opts *rv1.VolumeUpdateOptions) (*vv1.Volume, error)
	Remove(ctx context.Context, opts *rv1.VolumeRemoveOptions) error
}
