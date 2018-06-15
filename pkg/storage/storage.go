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
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
	"context"
)

const (
	NamespaceKind  types.Kind = "namespaces"
	ServiceKind    types.Kind = "services"
	DeploymentKind types.Kind = "deployments"
	ClusterKind    types.Kind = "cluster"
	PodKind        types.Kind = "pods"
	IngressKind    types.Kind = "ingresses"
	SystemKind     types.Kind = "systems"
	NodeKind       types.Kind = "nodes"
	EndpointKind   types.Kind = "endpoints"
	UtilsKind      types.Kind = "utils"
)

type Storage interface {
	Get(ctx context.Context, kind types.Kind, selfLink string, obj interface{}) error
	List(ctx context.Context, kind types.Kind, q string, obj interface{}) error
	Map(ctx context.Context, kind types.Kind, q string, obj interface{}) error
	Update(ctx context.Context, kind types.Kind, selfLink string, obj interface{}, opts *types.Opts) error
	Create(ctx context.Context, kind types.Kind, selfLink string, obj interface{}, opts *types.Opts) error
	Remove(ctx context.Context, kind types.Kind, selfLink string) error
	Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error
}

func Get(driver string) (Storage, error) {
	switch driver {
	default:
		return etcd.NewV3()
	}
}
