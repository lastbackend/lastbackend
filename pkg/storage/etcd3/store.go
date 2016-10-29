package etcd3

import (
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/runtime"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"golang.org/x/net/context"
	"path"
	"reflect"
	"strings"
)

type store struct {
	client     *clientv3.Client
	serializer runtime.Serializer
	versioner  storage.Versioner
	pathPrefix string
	//watcher    *watcher
}

// New returns an etcd3 implementation of storage.Interface.
func New(c *clientv3.Client, serializer runtime.Serializer, prefix string) storage.Interface {
	return newStore(c, serializer, prefix)
}

func newStore(c *clientv3.Client, serializer runtime.Serializer, prefix string) *store {
	versioner := struct{}{} // TODO
	return &store{
		client:     c,
		serializer: serializer,
		versioner:  versioner,
		pathPrefix: prefix,
		//watcher:    newWatcher(c, serializer, versioner),
	}
}

// ttlOpts returns client options based on given ttl.
// ttl: if ttl is non-zero, it will attach the key to a lease with ttl of roughly the same length
func (s *store) ttlOpts(ctx context.Context, ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	// TODO: one lease per ttl key is expensive. Based on current use case, we can have a long window to
	// put keys within into same lease. We shall benchmark this and optimize the performance.
	lcr, err := s.client.Lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(clientv3.LeaseID(lcr.ID))}, nil
}

// Versioner implements storage.Interface.Versioner.
func (s *store) Versioner() storage.Versioner {
	return s.versioner
}

// Get implements storage.Interface.Get.
func (s *store) Get(ctx context.Context, key string, out runtime.Object, ignoreNotFound bool) error {
	key = keyWithPrefix(s.pathPrefix, key)
	getResp, err := s.client.KV.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(getResp.Kvs) == 0 {
		if ignoreNotFound {
			return setZeroValue(out)
		}
		return errors.New("New key not found")
	}
	kv := getResp.Kvs[0]
	return decode(s.serializer, s.versioner, kv.Value, out, kv.ModRevision)
}

// Create implements storage.Interface.Create.
func (s *store) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	if version, err := s.versioner.ObjectResourceVersion(obj); err == nil && version != 0 {
		return errors.New("resourceVersion should not be set on objects to be created")
	}
	data, err := runtime.Encode(s.serializer, obj)
	if err != nil {
		return err
	}
	key = keyWithPrefix(s.pathPrefix, key)

	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		return err
	}

	txnResp, err := s.client.KV.Txn(ctx).If(
		notFound(key),
	).Then(
		clientv3.OpPut(key, string(data), opts...),
	).Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.New("New key exists")
	}

	if out != nil {
		putResp := txnResp.Responses[0].GetResponsePut()
		return decode(s.serializer, s.versioner, data, out, putResp.Header.Revision)
	}
	return nil
}

func keyWithPrefix(prefix, key string) string {
	if strings.HasPrefix(key, prefix) {
		return key
	}
	return path.Join(prefix, key)
}

func decode(serializer runtime.Serializer, versioner storage.Versioner, value []byte, obj runtime.Object, rev int64) error {
	if _, err := enforcePointer(obj); err != nil {
		panic("unable to convert output object to pointer")
	}
	_, err := serializer.Decode(value, obj)
	if err != nil {
		return err
	}

	// being unable to set the version does not prevent the object from being extracted
	versioner.UpdateObject(obj, uint64(rev)) // TODO: реализовать
	return nil
}

// SetZeroValue would set the object of obj to zero value of its type.
func setZeroValue(obj runtime.Object) error {
	v, err := enforcePointer(obj)
	if err != nil {
		return err
	}
	v.Set(reflect.Zero(v.Type()))
	return nil
}

// EnforcePtr ensures that obj is a pointer of some sort. Returns a reflect.Value
// of the dereferenced pointer, ensuring that it is settable/addressable.
// Returns an error if this is not possible.
func enforcePointer(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		if v.Kind() == reflect.Invalid {
			return reflect.Value{}, fmt.Errorf("expected pointer, but got invalid kind")
		}
		return reflect.Value{}, fmt.Errorf("expected pointer, but got %v type", v.Type())
	}
	if v.IsNil() {
		return reflect.Value{}, fmt.Errorf("expected pointer, but got nil")
	}
	return v.Elem(), nil
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}
