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
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"context"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	"github.com/lastbackend/lastbackend/pkg/storage/mock"
)

const (
	NamespaceKind  types.Kind = "namespace"
	ServiceKind    types.Kind = "service"
	DeploymentKind types.Kind = "deployment"
	ClusterKind    types.Kind = "cluster"
	PodKind        types.Kind = "pod"
	IngressKind    types.Kind = "ingress"
	SystemKind     types.Kind = "system"
	NodeKind       types.Kind = "node"
	RouteKind      types.Kind = "route"
	VolumeKind     types.Kind = "volume"
	TriggerKind    types.Kind = "trigger"
	SecretKind     types.Kind = "secret"
	EndpointKind   types.Kind = "endpoint"
	UtilsKind      types.Kind = "utils"
)

type Storage interface {
	Get(ctx context.Context, kind types.Kind, name string, obj interface{}) error
	List(ctx context.Context, kind types.Kind, q string, obj interface{}) error
	Map(ctx context.Context, kind types.Kind, q string, obj interface{}) error
	Create(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error
	Update(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error
	Upsert(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error
	Remove(ctx context.Context, kind types.Kind, name string) error
	Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error
	Filter() types.Filter
}

func Get(driver string) (Storage, error) {
	switch driver {
	case "mock":
		return mock.New()
	default:
		return etcd.New()
	}
}

func NewWatcher() chan *types.WatcherEvent {
	return make(chan *types.WatcherEvent)
}