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

	"github.com/lastbackend/lastbackend/pkg/storage/types"
)

const (
	logLevel  = 6
	logPrefix = "storage:mock"
)

type MockDB struct {
	store map[types.Kind]map[string]interface{}
}

func (db *MockDB) Get(ctx context.Context, kind types.Kind, name string, obj interface{}) error {
	db.check(kind)

	return nil
}

func (db *MockDB) List(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	db.check(kind)
	return nil
}

func (db *MockDB) Map(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	db.check(kind)
	return nil
}

func (db *MockDB) Create(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	db.check(kind)
	return nil
}

func (db *MockDB) Update(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	db.check(kind)
	return nil
}

func (db *MockDB) Upsert(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	db.check(kind)
	return nil
}

func (db *MockDB) Remove(ctx context.Context, kind types.Kind, name string) error {
	db.check(kind)
	return nil
}

func (db *MockDB) Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error {
	db.check(kind)
	return nil
}

func (db *MockDB) check(kind types.Kind) {
	if _, ok := db.store[kind]; !ok {
		db.store[kind] = make(map[string]interface{})
	}
}

func New() (*MockDB, error) {
	db := new(MockDB)
	db.store = make(map[types.Kind]map[string]interface{})
	return new(MockDB), nil
}
