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

package types

import (
	"io"
	"context"

	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type ClientV1 interface {
	Cluster() ClusterClientV1
	Namespace(args ...string) NamespaceClientV1
	Secret(args ...string) SecretClientV1
}

type ClusterClientV1 interface {
	Node(args ...string) NodeClientV1
	Ingress(args ...string) IngressClientV1

	Get(ctx context.Context) (*vv1.Cluster, error)
}

type NodeClientV1 interface {
	List(ctx context.Context) (*vv1.NodeList, error)
	Connect(ctx context.Context, opts *rv1.NodeConnectOptions) error
	Get(ctx context.Context) (*vv1.Node, error)
	GetSpec(ctx context.Context) (*vv1.NodeManifest, error)
	SetMeta(ctx context.Context, opts *rv1.NodeMetaOptions) (*vv1.Node, error)
	SetStatus(ctx context.Context, opts *rv1.NodeStatusOptions) error
	SetPodStatus(ctx context.Context, pod string, opts *rv1.NodePodStatusOptions) error
	SetVolumeStatus(ctx context.Context, volume string, opts *rv1.NodeVolumeStatusOptions) error
	SetRouteStatus(ctx context.Context, route string, opts *rv1.NodeRouteStatusOptions) error
	Remove(ctx context.Context, opts *rv1.NodeRemoveOptions) error
}

type IngressClientV1 interface {
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

type NamespaceClientV1 interface {
	Service(args ...string) ServiceClientV1
	Route(args ...string) RouteClientV1
	Volume(args ...string) VolumeClientV1
	Create(ctx context.Context, opts *rv1.NamespaceCreateOptions) (*vv1.Namespace, error)
	List(ctx context.Context) (*vv1.NamespaceList, error)
	Get(ctx context.Context) (*vv1.Namespace, error)
	Update(ctx context.Context, opts *rv1.NamespaceUpdateOptions) (*vv1.Namespace, error)
	Remove(ctx context.Context, opts *rv1.NamespaceRemoveOptions) error
}

type ServiceClientV1 interface {
	Deployment(args ...string) DeploymentClientV1
	Trigger(args ...string) TriggerClientV1

	Create(ctx context.Context, opts *rv1.ServiceManifest) (*vv1.Service, error)
	List(ctx context.Context) (*vv1.ServiceList, error)
	Get(ctx context.Context) (*vv1.Service, error)
	Update(ctx context.Context, opts *rv1.ServiceManifest) (*vv1.Service, error)
	Remove(ctx context.Context, opts *rv1.ServiceRemoveOptions) error
	Logs(ctx context.Context, opts *rv1.ServiceLogsOptions) (io.ReadCloser, error)
}

type DeploymentClientV1 interface {
	Pod(args ...string) PodClientV1

	List(ctx context.Context) (*vv1.DeploymentList, error)
	Get(ctx context.Context) (*vv1.Deployment, error)
	Update(ctx context.Context, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error)
}

type PodClientV1 interface {
	List(ctx context.Context) (*vv1.PodList, error)
	Get(ctx context.Context) (*vv1.Pod, error)
	Logs(ctx context.Context, opts *rv1.PodLogsOptions) (io.ReadCloser, error)
}

type EventsClientV1 interface {
}

type SecretClientV1 interface {
	Get(ctx context.Context) (*vv1.Secret, error)
	Create(ctx context.Context, opts *rv1.SecretCreateOptions) (*vv1.Secret, error)
	List(ctx context.Context) (*vv1.SecretList, error)
	Update(ctx context.Context, opts *rv1.SecretUpdateOptions) (*vv1.Secret, error)
	Remove(ctx context.Context, opts *rv1.SecretRemoveOptions) error
}

type RouteClientV1 interface {
	Create(ctx context.Context, opts *rv1.RouteCreateOptions) (*vv1.Route, error)
	List(ctx context.Context) (*vv1.RouteList, error)
	Get(ctx context.Context) (*vv1.Route, error)
	Update(ctx context.Context, opts *rv1.RouteUpdateOptions) (*vv1.Route, error)
	Remove(ctx context.Context, opts *rv1.RouteRemoveOptions) error
}

type TriggerClientV1 interface {
	Create(ctx context.Context, opts *rv1.TriggerCreateOptions) (*vv1.Trigger, error)
	List(ctx context.Context) (*vv1.TriggerList, error)
	Get(ctx context.Context) (*vv1.Trigger, error)
	Update(ctx context.Context, opts *rv1.TriggerUpdateOptions) (*vv1.Trigger, error)
	Remove(ctx context.Context, opts *rv1.TriggerRemoveOptions) error
}

type VolumeClientV1 interface {
	Create(ctx context.Context, opts *rv1.VolumeCreateOptions) (*vv1.Volume, error)
	List(ctx context.Context) (*vv1.VolumeList, error)
	Get(ctx context.Context) (*vv1.Volume, error)
	Update(ctx context.Context, opts *rv1.VolumeUpdateOptions) (*vv1.Volume, error)
	Remove(ctx context.Context, opts *rv1.VolumeRemoveOptions) error
}
