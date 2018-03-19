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
)

type Cluster interface {
	Get(ctx context.Context) (*vv1.ClusterList, error)
	Update(ctx context.Context, opts *rv1.ClusterUpdateOptions) (*vv1.Cluster, error)
}

type Deployment interface {
	List(ctx context.Context, namespace, service string) (*vv1.DeploymentList, error)
	Get(ctx context.Context, namespace, service, deployment string) (*vv1.Deployment, error)
	Update(ctx context.Context, namespace, service, deployment string, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error)
}

type Events interface {
}

type Namespace interface {
	Create(ctx context.Context, opts rv1.NamespaceCreateOptions) (*vv1.Namespace, error)
	List(ctx context.Context) (*vv1.NamespaceList, error)
	Get(ctx context.Context, name string) (*vv1.Namespace, error)
	Update(ctx context.Context, name string, opts rv1.NamespaceUpdateOptions) (*vv1.Namespace, error)
	Remove(ctx context.Context, name string, opts rv1.NamespaceRemoveOptions) error
}

type Node interface {
	List(ctx context.Context) (*vv1.NodeList, error)
	Get(ctx context.Context, name string) (*vv1.Node, error)
	GetSpec(ctx context.Context, name string) (*vv1.NodeSpec, error)
	Update(ctx context.Context, name string, opts rv1.NodeUpdateOptions) (*vv1.Node, error)
	SetInfo(ctx context.Context, name string, opts rv1.NodeInfoOptions) error
	SetState(ctx context.Context, name string, opts rv1.NodeStateOptions) error
	SetPodState(ctx context.Context, name string, opts rv1.NodeStateOptions) error
	SetVolumeState(ctx context.Context, name string, opts rv1.NodeStateOptions) error
	SetRouteState(ctx context.Context, name string, opts rv1.NodeStateOptions) error
	Remove(ctx context.Context, name string, opts rv1.NodeRemoveOptions) error
}

type Route interface {
	Create(ctx context.Context, namespace string, opts rv1.RouteCreateOptions) (*vv1.Route, error)
	List(ctx context.Context, namespace string) (*vv1.RouteList, error)
	Get(ctx context.Context, namespace, name string) (*vv1.Route, error)
	Update(ctx context.Context, namespace, name string, opts rv1.RouteUpdateOptions) (*vv1.Route, error)
	Remove(ctx context.Context, namespace, name string, opts rv1.RouteRemoveOptions) error
}

type Service interface {
	Create(ctx context.Context, namespace string, opts *rv1.ServiceCreateOptions) (*vv1.ServiceList, error)
	List(ctx context.Context, namespace string) (*vv1.ServiceList, error)
	Get(ctx context.Context, namespace, name string) (*vv1.Service, error)
	Update(ctx context.Context, namespace, name string, opts *rv1.ServiceUpdateOptions) (*vv1.NamespaceList, error)
	Remove(ctx context.Context, namespace, name string, opts rv1.ServiceRemoveOptions) error
}

type Trigger interface {
	Create(ctx context.Context, namespace, service string, opts rv1.TriggerCreateOptions) (*vv1.Trigger, error)
	List(ctx context.Context, namespace, service string) (*vv1.TriggerList, error)
	Get(ctx context.Context, namespace, service, name string) (*vv1.Trigger, error)
	Update(ctx context.Context, namespace, service, name string, opts rv1.TriggerUpdateOptions) (*vv1.Trigger, error)
	Remove(ctx context.Context, namespace, service, name string, opts rv1.TriggerRemoveOptions) error
}

type Volume interface {
	Create(ctx context.Context, namespace string, opts rv1.VolumeCreateOptions) (*vv1.Volume, error)
	List(ctx context.Context, namespace string) (*vv1.VolumeList, error)
	Get(ctx context.Context, namespace, name string) (*vv1.Volume, error)
	Update(ctx context.Context, namespace, name string, opts rv1.VolumeUpdateOptions) (*vv1.Volume, error)
	Remove(ctx context.Context, namespace, name string, opts rv1.VolumeRemoveOptions) error
}
