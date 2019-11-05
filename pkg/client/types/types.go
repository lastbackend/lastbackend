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
	"github.com/lastbackend/lastbackend/internal/api/types/v1/request"
	"github.com/lastbackend/lastbackend/internal/api/types/v1/views"
	"io"
	"net/http"
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
	API(args ...string) APIClientV1
	Controller(args ...string) ControllerClientV1
	Get(ctx context.Context) (*views.Cluster, error)
}

type NodeClientV1 interface {
	List(ctx context.Context) (*views.NodeList, error)
	Connect(ctx context.Context, opts *request.NodeConnectOptions) error
	Get(ctx context.Context) (*views.Node, error)
	SetStatus(ctx context.Context, opts *request.NodeStatusOptions) (*views.NodeManifest, error)
	Remove(ctx context.Context, opts *request.NodeRemoveOptions) error
}

type DiscoveryClientV1 interface {
	List(ctx context.Context) (*views.DiscoveryList, error)
	Get(ctx context.Context) (*views.Discovery, error)
	Connect(ctx context.Context, opts *request.DiscoveryConnectOptions) error
	SetStatus(ctx context.Context, opts *request.DiscoveryStatusOptions) (*views.DiscoveryManifest, error)
}

type IngressClientV1 interface {
	List(ctx context.Context) (*views.IngressList, error)
	Get(ctx context.Context) (*views.Ingress, error)
	Connect(ctx context.Context, opts *request.IngressConnectOptions) error
	SetStatus(ctx context.Context, opts *request.IngressStatusOptions) (*views.IngressManifest, error)
}

type ExporterClientV1 interface {
	List(ctx context.Context) (*views.ExporterList, error)
	Get(ctx context.Context) (*views.Exporter, error)
	Connect(ctx context.Context, opts *request.ExporterConnectOptions) error
	SetStatus(ctx context.Context, opts *request.ExporterStatusOptions) (*views.ExporterManifest, error)
}

type APIClientV1 interface {
	List(ctx context.Context) (*views.APIList, error)
	Get(ctx context.Context) (*views.API, error)
}

type ControllerClientV1 interface {
	List(ctx context.Context) (*views.ControllerList, error)
	Get(ctx context.Context) (*views.Controller, error)
}

type NamespaceClientV1 interface {
	Secret(args ...string) SecretClientV1
	Config(args ...string) ConfigClientV1
	Service(args ...string) ServiceClientV1
	Job(args ...string) JobClientV1
	Route(args ...string) RouteClientV1
	Volume(args ...string) VolumeClientV1
	Create(ctx context.Context, opts *request.NamespaceManifest) (*views.Namespace, error)
	Apply(ctx context.Context, opts *request.NamespaceApplyManifest) (*views.NamespaceApplyStatus, error)
	List(ctx context.Context) (*views.NamespaceList, error)
	Get(ctx context.Context) (*views.Namespace, error)
	Update(ctx context.Context, opts *request.NamespaceManifest) (*views.Namespace, error)
	Remove(ctx context.Context, opts *request.NamespaceRemoveOptions) error
}

type ServiceClientV1 interface {
	Deployment(args ...string) DeploymentClientV1
	Create(ctx context.Context, opts *request.ServiceManifest) (*views.Service, error)
	List(ctx context.Context) (*views.ServiceList, error)
	Get(ctx context.Context) (*views.Service, error)
	Update(ctx context.Context, opts *request.ServiceManifest) (*views.Service, error)
	Remove(ctx context.Context, opts *request.ServiceRemoveOptions) error
	Logs(ctx context.Context, opts *request.ServiceLogsOptions) (io.ReadCloser, *http.Response, error)
}

type JobClientV1 interface {
	Task(args ...string) TaskClientV1

	Create(ctx context.Context, opts *request.JobManifest) (*views.Job, error)
	Run(ctx context.Context, opts *request.TaskManifest) (*views.Task, error)
	List(ctx context.Context) (*views.JobList, error)
	Get(ctx context.Context) (*views.Job, error)
	Update(ctx context.Context, opts *request.JobManifest) (*views.Job, error)
	Remove(ctx context.Context, opts *request.JobRemoveOptions) error
	Logs(ctx context.Context, opts *request.JobLogsOptions) (io.ReadCloser, *http.Response, error)
}

type TaskClientV1 interface {
	Pod(args ...string) PodClientV1

	Create(ctx context.Context, opts *request.TaskManifest) (*views.Task, error)
	List(ctx context.Context) (*views.TaskList, error)
	Get(ctx context.Context) (*views.Task, error)
	Cancel(ctx context.Context, opts *request.TaskCancelOptions) (*views.Task, error)
	Remove(ctx context.Context, opts *request.TaskRemoveOptions) error
}

type DeploymentClientV1 interface {
	Pod(args ...string) PodClientV1
	List(ctx context.Context) (*views.DeploymentList, error)
	Get(ctx context.Context) (*views.Deployment, error)
	Create(ctx context.Context, opts *request.DeploymentManifest) (*views.Deployment, error)
	Update(ctx context.Context, opts *request.DeploymentManifest) (*views.Deployment, error)
	Remove(ctx context.Context, opts *request.DeploymentRemoveOptions) error
}

type PodClientV1 interface {
	List(ctx context.Context) (*views.PodList, error)
	Get(ctx context.Context) (*views.Pod, error)
	Logs(ctx context.Context, opts *request.PodLogsOptions) (io.ReadCloser, *http.Response, error)
}

type EventsClientV1 interface {
}

type SecretClientV1 interface {
	Get(ctx context.Context) (*views.Secret, error)
	Create(ctx context.Context, opts *request.SecretManifest) (*views.Secret, error)
	List(ctx context.Context) (*views.SecretList, error)
	Update(ctx context.Context, opts *request.SecretManifest) (*views.Secret, error)
	Remove(ctx context.Context, opts *request.SecretRemoveOptions) error
}

type ConfigClientV1 interface {
	Get(ctx context.Context) (*views.Config, error)
	Create(ctx context.Context, opts *request.ConfigManifest) (*views.Config, error)
	List(ctx context.Context) (*views.ConfigList, error)
	Update(ctx context.Context, opts *request.ConfigManifest) (*views.Config, error)
	Remove(ctx context.Context, opts *request.ConfigRemoveOptions) error
}

type RouteClientV1 interface {
	Create(ctx context.Context, opts *request.RouteManifest) (*views.Route, error)
	List(ctx context.Context) (*views.RouteList, error)
	Get(ctx context.Context) (*views.Route, error)
	Update(ctx context.Context, opts *request.RouteManifest) (*views.Route, error)
	Remove(ctx context.Context, opts *request.RouteRemoveOptions) error
}

type VolumeClientV1 interface {
	Create(ctx context.Context, opts *request.VolumeManifest) (*views.Volume, error)
	List(ctx context.Context) (*views.VolumeList, error)
	Get(ctx context.Context) (*views.Volume, error)
	Update(ctx context.Context, opts *request.VolumeManifest) (*views.Volume, error)
	Remove(ctx context.Context, opts *request.VolumeRemoveOptions) error
}
