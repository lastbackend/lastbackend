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

package etcd3

import (
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"github.com/lastbackend/lastbackend/pkg/serializer"
	st "github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"golang.org/x/net/context"
	"path"
	"reflect"
	"strings"
)

type store struct {
	client *clientv3.Client
	// getOpts contains additional options that should be passed to all Get() calls.
	getOps     []clientv3.OpOption
	codec      serializer.Codec
	pathPrefix string
}

type itemForDecode []byte

// Create implements store.Interface.Create.
// You can optionally set a TTL for a key to expire in a certain number of seconds.
func (s *store) Create(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		return err
	}
	key = path.Join(s.pathPrefix, key)
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		return err
	}
	fmt.Println("Create:", key, string(data))
	txnResp, err := s.client.KV.Txn(ctx).If(notFound(key)).
		Then(clientv3.OpPut(key, string(data), opts...)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrKeyExists)
	}

	if outPtr != nil {
		return decode(s.codec, data, outPtr)
	}

	return nil
}

// Get implements store.Interface.Get.
func (s *store) Get(ctx context.Context, key string, outPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	fmt.Println("Get:", key)
	res, err := s.client.KV.Get(ctx, key, s.getOps...)
	if err != nil {
		return err
	}
	if len(res.Kvs) == 0 {
		return nil
	}
	return decode(s.codec, res.Kvs[0].Value, outPtr)
}

// List implements storage.Interface.List.
func (s *store) List(ctx context.Context, key string, listOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	// We need to make sure the key ended with "/" so that we only get children "directories".
	// e.g. if we have key "/a", "/a/b", "/ab", getting keys with prefix "/a" will return all three,
	// while with prefix "/a/" will return only "/a/b" which is the correct answer.
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}

	fmt.Println("List:", key)
	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	items := make([]itemForDecode, 0, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		items = append(items, itemForDecode(kv.Value))
	}

	return decodeList(items, listOutPtr, s.codec)
}

// Delete implements store.Interface.Delete.
func (s *store) Delete(ctx context.Context, key string, outPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	// We need to do get and delete in single transaction in order to
	// know the value and revision before deleting it.
	fmt.Println("Del:", key)
	txnResp, err := s.client.KV.Txn(ctx).If().Then(
		clientv3.OpGet(key),
		clientv3.OpDelete(key),
	).Commit()
	if err != nil {
		return err
	}
	if validator.IsNil(outPtr) {
		return nil
	}

	getResp := txnResp.Responses[0].GetResponseRange()
	if len(getResp.Kvs) == 0 {
		return errors.New(st.ErrKeyNotFound)
	}
	return decode(s.codec, getResp.Kvs[0].Value, outPtr)
}

// Decode decodes value of bytes into object.
// On success, objPtr would be set to the object.
func decode(s serializer.Codec, value []byte, outPtr interface{}) error {
	if _, err := converter.EnforcePtr(outPtr); err != nil {
		panic("unable to convert output object to pointer")
	}
	return serializer.Decode(s, value, outPtr)
}

// decodeList decodes a list of values into a list of objects.
// On success, ListObjPtr would be set to the list of objects.
func decodeList(items []itemForDecode, ListOutPtr interface{}, codec serializer.Codec) error {
	v, err := converter.EnforcePtr(ListOutPtr)
	if err != nil || v.Kind() != reflect.Slice {
		panic("need ptr to slice")
	}

	for _, item := range items {
		var obj = reflect.New(v.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, item, obj)
		if err != nil {
			return err
		}
		v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
	}
	return nil
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}

// ttlOpts returns client options based on given ttl.
// ttl: if ttl is non-zero, it will attach the key to a lease with ttl of roughly the same length
func (s *store) ttlOpts(ctx context.Context, ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	lcr, err := s.client.Lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(clientv3.LeaseID(lcr.ID))}, nil
}