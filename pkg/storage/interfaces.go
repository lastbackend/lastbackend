package storage

import (
	"github.com/lastbackend/lastbackend/pkg/runtime"
	"golang.org/x/net/context"
)

// MatchValue defines a pair (<index name>, <value for that index>).
type MatchValue struct {
	IndexName string
	Value     string
}

type Filter interface {
	Filter(obj runtime.Object) bool
	Trigger() []MatchValue
}

// Versioner abstracts setting and retrieving metadata fields from database response
// onto the object ot list.
type Versioner interface {
	UpdateObject(obj runtime.Object, resourceVersion uint64) error
	UpdateList(obj runtime.Object, resourceVersion uint64) error
	ObjectResourceVersion(obj runtime.Object) (uint64, error)
}

// Interface offers a common interface for object marshaling/unmarshaling operations and
// hides all the storage-related operations behind it.
type Interface interface {
	Versioner() Versioner
	Get(ctx context.Context, key string, obj runtime.Object, ignoreNotFound bool) error
	//GetToList(ctx context.Context, key string, filter Filter, listObj runtime.Object) error
	//List(ctx context.Context, key string, filter Filter, listObj runtime.Object) error
	Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error
	//Delete(ctx context.Context, key string, out runtime.Object, preconditions *Preconditions) error
	//Watch(ctx context.Context, key string, filter Filter) (watch.Interface, error)
	//WatchList(ctx context.Context, key string, filter Filter) (watch.Interface, error)
}
