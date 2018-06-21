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

package mock

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
)

const (
	logLevel  = 6
	logPrefix = "storage:etcd:v3:mock"
)

type MockDB struct {
}

func New() (*MockDB, error) {
	return new(MockDB), nil
}

func (MockDB) Get(ctx context.Context, kind types.Kind, name string, obj interface{}) error {
	return nil
}

func (MockDB) List(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	return nil
}

func (MockDB) Map(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	return nil
}

func (MockDB) Create(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	return nil
}

func (MockDB) Update(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	return nil
}

func (MockDB) Upsert(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	return nil
}

func (MockDB) Remove(ctx context.Context, kind types.Kind, name string) error {
	return nil
}

func (MockDB) Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error {
	return nil
}
