//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"context"
	"io"
	"net/http"

	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type ClientV1 interface {
	Cluster() ClusterClientV1
	Namespace(args ...string) NamespaceClientV1
}

type ClusterClientV1 interface {
	Node(args ...string) NodeClientV1
	Ingress(args ...string) IngressClientV1
	Exporter(args ...string) ExporterClientV1
	Discovery(args ...string) DiscoveryClientV1
	Get(ctx context.Context) (*vv1.Cluster, error)
}

type NodeClientV1 interface {
	List(ctx context.Context) (*vv1.NodeList, error)
	Connect(ctx context.Context, opts *rv1.NodeConnectOptions) error
	Get(ctx context.Context) (*vv1.Node, error)
	SetStatus(ctx context.Context, opts *rv1.NodeStatusOptions) (*vv1.NodeManifest, error)
	Remove(ctx context.Context, opts *rv1.NodeRemoveOptions) error
}

type DiscoveryClientV1 interface {
	List(ctx context.Context) (*vv1.DiscoveryList, error)
	Get(ctx context.Context) (*vv1.Discovery, error)
	Connect(ctx context.Context, opts *rv1.DiscoveryConnectOptions) error
	SetStatus(ctx context.Context, opts *rv1.DiscoveryStatusOptions) (*vv1.DiscoveryManifest, error)
}

type IngressClientV1 interface {
	List(ctx context.Context) (*vv1.IngressList, error)
	Get(ctx context.Context) (*vv1.Ingress, error)
	Connect(ctx context.Context, opts *rv1.IngressConnectOptions) error
	SetStatus(ctx context.Context, opts *rv1.IngressStatusOptions) (*vv1.IngressManifest, error)
}

type ExporterClientV1 interface {
	List(ctx context.Context) (*vv1.ExporterList, error)
	Get(ctx context.Context) (*vv1.Exporter, error)
	Connect(ctx context.Context, opts *rv1.ExporterConnectOptions) error
	SetStatus(ctx context.Context, opts *rv1.ExporterStatusOptions) (*vv1.ExporterManifest, error)
}

type NamespaceClientV1 interface {
	Secret(args ...string) SecretClientV1
	Config(args ...string) ConfigClientV1
	Service(args ...string) ServiceClientV1
	Job(args ...string) JobClientV1
	Route(args ...string) RouteClientV1
	Volume(args ...string) VolumeClientV1
	Create(ctx context.Context, opts *rv1.NamespaceManifest) (*vv1.Namespace, error)
	Apply(ctx context.Context, opts *rv1.NamespaceApplyManifest) (*vv1.NamespaceApplyStatus, error)
	List(ctx context.Context) (*vv1.NamespaceList, error)
	Get(ctx context.Context) (*vv1.Namespace, error)
	Update(ctx context.Context, opts *rv1.NamespaceManifest) (*vv1.Namespace, error)
	Remove(ctx context.Context, opts *rv1.NamespaceRemoveOptions) error
}

type ServiceClientV1 interface {
	Deployment(args ...string) DeploymentClientV1
	Create(ctx context.Context, opts *rv1.ServiceManifest) (*vv1.Service, error)
	List(ctx context.Context) (*vv1.ServiceList, error)
	Get(ctx context.Context) (*vv1.Service, error)
	Update(ctx context.Context, opts *rv1.ServiceManifest) (*vv1.Service, error)
	Remove(ctx context.Context, opts *rv1.ServiceRemoveOptions) error
	Logs(ctx context.Context, opts *rv1.ServiceLogsOptions) (io.ReadCloser, *http.Response, error)
}

type JobClientV1 interface {
	Tasks(args ...string) TaskClientV1
	Create(ctx context.Context, opts *rv1.JobManifest) (*vv1.Job, error)
	Run(ctx context.Context, opts *rv1.TaskManifest) (*vv1.Task, error)
	List(ctx context.Context) (*vv1.JobList, error)
	Get(ctx context.Context) (*vv1.Job, error)
	Update(ctx context.Context, opts *rv1.JobManifest) (*vv1.Job, error)
	Remove(ctx context.Context, opts *rv1.JobRemoveOptions) error
	Logs(ctx context.Context, opts *rv1.JobLogsOptions) (io.ReadCloser, *http.Response, error)
}

type TaskClientV1 interface {
	Pod(args ...string) PodClientV1

	List(ctx context.Context) (*vv1.TaskList, error)
	Get(ctx context.Context) (*vv1.Task, error)
	Cancel(ctx context.Context, opts *rv1.TaskCancelOptions) (*vv1.Task, error)
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
	Logs(ctx context.Context, opts *rv1.PodLogsOptions) (io.ReadCloser, *http.Response, error)
}

type EventsClientV1 interface {
}

type SecretClientV1 interface {
	Get(ctx context.Context) (*vv1.Secret, error)
	Create(ctx context.Context, opts *rv1.SecretManifest) (*vv1.Secret, error)
	List(ctx context.Context) (*vv1.SecretList, error)
	Update(ctx context.Context, opts *rv1.SecretManifest) (*vv1.Secret, error)
	Remove(ctx context.Context, opts *rv1.SecretRemoveOptions) error
}

type ConfigClientV1 interface {
	Get(ctx context.Context) (*vv1.Config, error)
	Create(ctx context.Context, opts *rv1.ConfigManifest) (*vv1.Config, error)
	List(ctx context.Context) (*vv1.ConfigList, error)
	Update(ctx context.Context, opts *rv1.ConfigManifest) (*vv1.Config, error)
	Remove(ctx context.Context, opts *rv1.ConfigRemoveOptions) error
}

type RouteClientV1 interface {
	Create(ctx context.Context, opts *rv1.RouteManifest) (*vv1.Route, error)
	List(ctx context.Context) (*vv1.RouteList, error)
	Get(ctx context.Context) (*vv1.Route, error)
	Update(ctx context.Context, opts *rv1.RouteManifest) (*vv1.Route, error)
	Remove(ctx context.Context, opts *rv1.RouteRemoveOptions) error
}

type VolumeClientV1 interface {
	Create(ctx context.Context, opts *rv1.VolumeManifest) (*vv1.Volume, error)
	List(ctx context.Context) (*vv1.VolumeList, error)
	Get(ctx context.Context) (*vv1.Volume, error)
	Update(ctx context.Context, opts *rv1.VolumeManifest) (*vv1.Volume, error)
	Remove(ctx context.Context, opts *rv1.VolumeRemoveOptions) error
}
