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

type Storage struct {
	store map[types.Kind]map[string]interface{}
}

func (s *Storage) Get(ctx context.Context, kind types.Kind, name string, obj interface{}) error {
	s.check(kind)

	return nil
}

func (s *Storage) List(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	s.check(kind)
	return nil
}

func (s *Storage) Map(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	s.check(kind)
	return nil
}

func (s *Storage) Create(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	s.check(kind)
	return nil
}

func (s *Storage) Update(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	s.check(kind)
	return nil
}

func (s *Storage) Upsert(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	s.check(kind)
	return nil
}

func (s *Storage) Remove(ctx context.Context, kind types.Kind, name string) error {
	s.check(kind)
	return nil
}

func (s *Storage) Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error {
	s.check(kind)
	return nil
}

func (s Storage) Filter() types.Filter {
	return new(Filter)
}

func (s Storage) Key() types.Key {
	return new(Key)
}

func (s *Storage) check(kind types.Kind) {
	if _, ok := s.store[kind]; !ok {
		s.store[kind] = make(map[string]interface{})
	}
}

func New() (*Storage, error) {
	db := new(Storage)
	db.store = make(map[types.Kind]map[string]interface{})
	return db, nil
}
