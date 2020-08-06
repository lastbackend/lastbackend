//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"fmt"
	"github.com/lastbackend/lastbackend/internal/pkg/storage/bbolt"
	"github.com/lastbackend/lastbackend/internal/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/internal/pkg/storage/types"
)

const (
	MockDriver  = "mock"
	BboltDriver = "bbolt"
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

type IStorage interface {
	List(collection string, listOutPtr interface{}) error
	Get(collection, key string, outPtr interface{}) error
	Set(collection, key string, obj interface{}) error
	Put(collection, key string, obj interface{}) error
	Del(collection, key string) error
	Close() error
}

type BboltConfig bbolt.Options

func Get(driver string, opts interface{}) (IStorage, error) {

	if driver != MockDriver && opts == nil {
		return nil, fmt.Errorf("options can not be is nil")
	}

	switch driver {
	case MockDriver:
		return mock.New()
	case BboltDriver:
		fallthrough
	default:
		o := opts.(BboltConfig)
		return bbolt.New(bbolt.Options(o))
	}
}

func GetOpts() *types.Opts {
	return new(types.Opts)
}

func NewWatcher() chan *types.WatcherEvent {
	return make(chan *types.WatcherEvent)
}
