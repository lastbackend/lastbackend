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

	vv1 "github.com/lastbackend/lastbackend/pkg/api/views/v1"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/request/v1"
)

type Cluster interface {
	Get(ctx context.Context) (*vv1.NamespaceList, error)
	Update(ctx context.Context, opts *rv1.ClusterUpdateOpts) (*vv1.Cluster, error)
}

type Deployment interface {
	List(ctx context.Context, namespace, service string) (*vv1.DeploymentList, error)
	Get(ctx context.Context, namespace, service, deployment string) (*vv1.Deployment, error)
	Update(ctx context.Context, namespace, service, deployment string, opts *rv1.DeploymentUpdateOpts) (*vv1.Deployment, error)
}

type Events interface {
}

type Namespace interface {
	Create(ctx context.Context, opts rv1.NamespaceCreateOpts) (*vv1.Namespace, error)
	List(ctx context.Context) (*vv1.NamespaceList, error)
	Get(ctx context.Context, name string) (*vv1.Namespace, error)
	Update(ctx context.Context, name string, opts rv1.NamespaceUpdateOpts) (*vv1.Namespace, error)
	Remove(ctx context.Context, name string, opts rv1.NamespaceRemoveOpts) error
}

type Node interface {
	List(ctx context.Context) (*vv1.NodeList, error)
	Get(ctx context.Context, name string) (*vv1.Node, error)
	Update(ctx context.Context, name string, opts rv1.NodeUpdateOpts) (*vv1.Node, error)
	Remove(ctx context.Context, name string, opts rv1.NodeRemoveOpts) error
}

type Route interface {
	Create(ctx context.Context, namespace string, opts rv1.RouteCreateOpts) (*vv1.Route, error)
	List(ctx context.Context, namespace string) (*vv1.RouteList, error)
	Get(ctx context.Context, namespace, name string) (*vv1.Route, error)
	Update(ctx context.Context, namespace, name string, opts rv1.RouteUpdateOpts) (*vv1.Route, error)
	Remove(ctx context.Context, namespace, name string, opts rv1.RouteRemoveOpts) error
}

type Service interface {
	Create(ctx context.Context, namespace string, opts *rv1.ServiceCreateOpts) (*vv1.ServiceList, error)
	List(ctx context.Context, namespace string) (*vv1.ServiceList, error)
	Get(ctx context.Context, namespace, name string) (*vv1.Service, error)
	Update(ctx context.Context, namespace, name string, opts *rv1.ServiceUpdateOpts) (*vv1.NamespaceList, error)
	Remove(ctx context.Context, namespace, name string, opts rv1.ServiceRemoveOpts) error
}

type Trigger interface {
	Create(ctx context.Context, namespace, service string, opts rv1.TriggerCreateOpts) (*vv1.Trigger, error)
	List(ctx context.Context, namespace, service string) (*vv1.TriggerList, error)
	Get(ctx context.Context, namespace, service, name string) (*vv1.Trigger, error)
	Update(ctx context.Context, namespace, service, name string, opts rv1.TriggerUpdateOpts) (*vv1.Trigger, error)
	Remove(ctx context.Context, namespace, service, name string, opts rv1.TriggerRemoveOpts) error
}

type Volume interface {
	Create(ctx context.Context, namespace string, opts rv1.VolumeCreateOpts) (*vv1.Volume, error)
	List(ctx context.Context, namespace string) (*vv1.VolumeList, error)
	Get(ctx context.Context, namespace, name string) (*vv1.Volume, error)
	Update(ctx context.Context, namespace, name string, opts rv1.VolumeUpdateOpts) (*vv1.Volume, error)
	Remove(ctx context.Context, namespace, name string, opts rv1.VolumeRemoveOpts) error
}
