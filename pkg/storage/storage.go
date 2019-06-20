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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	v3 "github.com/lastbackend/lastbackend/pkg/storage/etcd/v3"
	"github.com/lastbackend/lastbackend/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/spf13/viper"
)

const (
	NamespaceKind  types.Kind = "namespace"
	ServiceKind    types.Kind = "service"
	DeploymentKind types.Kind = "deployment"
	ClusterKind    types.Kind = "cluster"
	PodKind        types.Kind = "pod"
	IngressKind    types.Kind = "ingress"
	ExporterKind   types.Kind = "exporter"
	SystemKind     types.Kind = "system"
	NodeKind       types.Kind = "node"
	RouteKind      types.Kind = "route"
	VolumeKind     types.Kind = "volume"
	TriggerKind    types.Kind = "trigger"
	SecretKind     types.Kind = "secret"
	EndpointKind   types.Kind = "endpoint"
	UtilsKind      types.Kind = "utils"
	ManifestKind   types.Kind = "manifest"
	NetworkKind    types.Kind = "network"
	SubnetKind     types.Kind = "subnet"
	TaskKind       types.Kind = "task"
	JobKind        types.Kind = "job"
	TestKind       types.Kind = "test"
)

type Storage interface {
	Info(ctx context.Context, collection, name string) (*types.System, error)
	Get(ctx context.Context, collection, name string, obj interface{}, opts *types.Opts) error
	List(ctx context.Context, collection, q string, obj interface{}, opts *types.Opts) error
	Map(ctx context.Context, collection, q string, obj interface{}, opts *types.Opts) error
	Put(ctx context.Context, collection, name string, obj interface{}, opts *types.Opts) error
	Set(ctx context.Context, collection, name string, obj interface{}, opts *types.Opts) error
	Del(ctx context.Context, collection, name string) error
	Watch(ctx context.Context, collection string, event chan *types.WatcherEvent, opts *types.Opts) error
	Collection() types.Collection
	Filter() types.Filter
}

func Get(v *viper.Viper) (Storage, error) {

	if !v.IsSet("storage.driver") {
		return nil, errors.New("storage driver not set")
	}

	switch v.GetString("storage.driver") {
	case "mock":
		return mock.New()
	default:

		config := new(v3.Config)

		config.Prefix = v.GetString("storage.etcd.prefix")
		config.Endpoints = v.GetStringSlice("storage.etcd.endpoints")

		config.TLS.CA = v.GetString("storage.etcd.tls.ca")
		config.TLS.Cert = v.GetString("storage.etcd.tls.cert")
		config.TLS.Key = v.GetString("storage.etcd.tls.key")

		return etcd.New(config)
	}
}

func GetOpts() *types.Opts {
	return new(types.Opts)
}

func NewWatcher() chan *types.WatcherEvent {
	return make(chan *types.WatcherEvent)
}
