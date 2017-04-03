//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"golang.org/x/net/context"
)

const (
	ErrKeyExists        = "key exists"
	ErrOperationFailure = "operation failure"
	ErrKeyNotFound      = "key not found"
)

type DestroyFunc func()

type FilterFunc func(obj interface{}) bool

type IStore interface {
	Create(ctx context.Context, key string, obj, out interface{}, ttl uint64) error
	Get(ctx context.Context, key string, objPtr interface{}) error
	List(ctx context.Context, key, keyRegexpFilter string, listObjPtr interface{}) error
	Update(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error
	Delete(ctx context.Context, key string, out interface{}) error
	Begin(ctx context.Context) ITx
}

type ITx interface {
	Create(key string, obj interface{}, ttl uint64) error
	Update(key string, obj interface{}, ttl uint64) error
	Delete(key string)
	Commit() error
}
