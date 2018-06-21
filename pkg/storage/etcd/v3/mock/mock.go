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

	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/types"

	s "github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
)

const (
	logLevel  = 6
	logPrefix = "storage:etcd:v3:mock"
)

type mockstore struct {
	store.Store
}

func GetMockClient() (s.Store, s.DestroyFunc, error) {
	var (
		df store.DestroyFunc = func() {}
		st                   = new(mockstore)
	)
	return st, df, nil
}

func (s *mockstore) Count(ctx context.Context, key, keyRegexFilter string) (int, error) {
	return int(0), nil
}

func (s *mockstore) Create(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	return nil
}

func (s *mockstore) Get(ctx context.Context, key string, outPtr interface{}) error {
	return nil
}

func (s *mockstore) List(ctx context.Context, key, keyRegexFilter string, listOutPtr interface{}) error {
	return nil
}

func (s *mockstore) Map(ctx context.Context, key, keyRegexFilter string, mapOutPtr interface{}) error {
	return nil
}

func (s *mockstore) MapList(ctx context.Context, key string, keyRegexFilter string, mapOutPtr interface{}) error {
	return nil
}

func (s *mockstore) Update(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	return nil
}

func (s *mockstore) Upsert(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	return nil
}

func (s *mockstore) Delete(ctx context.Context, key string) error {
	return nil
}

func (s *mockstore) DeleteDir(ctx context.Context, key string) error {
	return nil
}

func (s *mockstore) Watch(ctx context.Context, key, keyRegexFilter string) (types.Watcher, error) {
	return newWatcher(), nil
}

func (s *mockstore) Begin(ctx context.Context) store.TX {
	return nil
}

func (s *mockstore) Decode(ctx context.Context, value []byte, out interface{}) error {
	return nil
}
