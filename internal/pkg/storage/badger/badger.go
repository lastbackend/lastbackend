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

package badger

import (
	"context"
	badger "github.com/dgraph-io/badger"
	"github.com/lastbackend/lastbackend/internal/pkg/storage/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

type Storage struct {
}

func (s *Storage) Info(ctx context.Context, collection string, name string) (*types.System, error) {
	return new(types.System), nil
}

func (s *Storage) Get(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {

	return nil
}

func (s *Storage) List(ctx context.Context, collection string, q string, obj interface{}, opts *types.Opts) error {

	return nil
}

func (s *Storage) Map(ctx context.Context, collection string, q string, obj interface{}, opts *types.Opts) error {

	return nil
}

func (s *Storage) Put(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {

	return nil
}

func (s *Storage) Set(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {

	return nil
}

func (s *Storage) Del(ctx context.Context, collection string, name string) error {

	return nil
}

func (s *Storage) Watch(ctx context.Context, collection string, event chan *types.WatcherEvent, opts *types.Opts) error {

	<-ctx.Done()
	return nil
}

func (s Storage) Filter() types.Filter {
	return new(Filter)
}

func (s Storage) Key() types.Key {
	return new(Key)
}

func (s Storage) Collection() types.Collection {
	return new(Collection)
}

func New() (*Storage, error) {
	db := new(Storage)
	return db, nil
}

func test() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
