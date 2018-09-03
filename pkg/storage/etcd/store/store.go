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

package store

import (
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"golang.org/x/net/context"
)

type DestroyFunc func()

type Store interface {
	Info(ctx context.Context, key string) (*types.Runtime, error)
	Count(ctx context.Context, key, keyRegexFilter string) (int, error)
	Put(ctx context.Context, key string, obj, out interface{}, ttl uint64) error
	Get(ctx context.Context, key string, objPtr interface{}, rev *int64) error
	List(ctx context.Context, key, filter string, listObjPtr interface{}, rev *int64) error
	Map(ctx context.Context, key, filter string, mapObj interface{}, rev *int64) error
	Set(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64, force bool, rev *int64) error
	Del(ctx context.Context, key string) error
	Watch(ctx context.Context, key, filter string, rev *int64) (types.Watcher, error)
	Begin(ctx context.Context) TX
	Decode(ctx context.Context, value []byte, out interface{}) error
}

type TX interface {
	Put(key string, obj interface{}, ttl uint64) error
	Set(key string, obj interface{}, ttl uint64, force bool) error
	Del(key string)
	Commit() error
}
