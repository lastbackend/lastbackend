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
	"github.com/lastbackend/lastbackend/pkg/storage/etcd3"
	"golang.org/x/net/context"
)

const (
	ErrKeyExists        = "key exists"
	ErrOperationFailure = "operation failure"
	ErrKeyNotFound      = "key not found"
	ErrUnreachable      = "server unreachable"
)

// DestroyFunc is to destroy any resources used by the storage returned in Create() together.
type DestroyFunc func()

// FilterFunc takes an API object and returns true if the object satisfies some requirements.
type FilterFunc func(obj interface{}) bool

// Interface offers a common interface for object marshaling/unmarshaling operations and
// hides all the storage-related operations behind it.
type Interface interface {
	// Create adds a new object at a key unless it already exists.
	Create(ctx context.Context, key string, obj, out interface{}, ttl uint64) error
	// Get unmarshals json found at key into objPtr. On a not found error, will either
	// return a zero object.
	Get(ctx context.Context, key string, objPtr interface{}) error
	// List unmarshalls jsons found at directory defined by key and opaque them
	// into list object.
	List(ctx context.Context, key string, listObj interface{}) error
	// Delete removes the specified key and returns the value that existed at that spot.
	// If key didn't exist, it will return NotFound storage error.
	Delete(ctx context.Context, key string, out interface{}) error
	// Delete removes the specified key and returns the value that existed at that spot.
	// If key didn't exist, it will return NotFound storage error.
	Begin(ctx context.Context) *etcd3.Tx
}

