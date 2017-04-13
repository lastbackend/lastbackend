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
	"github.com/coreos/etcd/clientv3"
	st "github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"golang.org/x/net/context"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type store struct {
	debug      bool
	client     *clientv3.Client
	opts       []clientv3.OpOption
	codec      serializer.Codec
	pathPrefix string
}

// Need for decode array bytes
type buffer []byte

func (s *store) Count(ctx context.Context, key, keyRegexFilter string) (int, error) {
	key = path.Join(s.pathPrefix, key)
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	printf("COUNT: %s", key)
	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	r, _ := regexp.Compile(keyRegexFilter)

	if len(keyRegexFilter) == 0 {
		return len(getResp.Kvs), nil
	}

	count := 0
	for _, kv := range getResp.Kvs {
		if r.MatchString(string(kv.Key)) {
			count++
		}
	}
	return count, nil
}

func (s *store) Create(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		return err
	}
	key = path.Join(s.pathPrefix, key)
	printf("CREATE: %s", key)
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, string(data), opts...)).Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrKeyExists)
	}
	if validator.IsNil(outPtr) {
		return nil
	}
	if outPtr != nil {
		return decode(s.codec, data, outPtr)
	}
	return nil
}

func (s *store) Get(ctx context.Context, key string, outPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	printf("GET: %s", key)
	res, err := s.client.KV.Get(ctx, key, s.opts...)
	if err != nil {
		return err
	}
	if len(res.Kvs) == 0 {
		return errors.New(st.ErrKeyNotFound)
	}
	return decode(s.codec, res.Kvs[0].Value, outPtr)
}

func (s *store) List(ctx context.Context, key, keyRegexFilter string, listOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	printf("LIST: %s", key)
	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r, _ := regexp.Compile(keyRegexFilter)
	items := make([]buffer, 0, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		if r.MatchString(string(kv.Key)) {
			items = append(items, buffer(kv.Value))
		}
	}
	return decodeList(s.codec, items, listOutPtr)
}

func (s *store) Map(ctx context.Context, key, keyRegexFilter string, mapOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	printf("MAP: %s", key)
	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]buffer, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		if (keyRegexFilter == "") || r.MatchString(string(kv.Key)) {

			var (
				key   []string
				index string
			)

			match := r.FindStringSubmatch(string(kv.Key))
			if len(match) > 1 {
				index = match[1]
			} else {
				key = strings.Split(string(kv.Key), "/")
				index = key[len(key)-1]
			}

			items[index] = buffer(kv.Value)
		}

	}
	return decodeMap(s.codec, items, mapOutPtr)
}

func (s *store) Update(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		return err
	}
	key = path.Join(s.pathPrefix, key)
	printf("UPDATE: %s", key)
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "!=", 0)).
		Then(clientv3.OpPut(key, string(data), opts...)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrKeyNotFound)
	}
	if validator.IsNil(outPtr) {
		return nil
	}
	if outPtr != nil {
		return decode(s.codec, data, outPtr)
	}
	return nil
}

func (s *store) Delete(ctx context.Context, key string) error {
	key = path.Join(s.pathPrefix, key)
	printf("DELETE: %s", key)
	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpGet(key), clientv3.OpDelete(key)).
		Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteDir(ctx context.Context, key string) error {
	key = path.Join(s.pathPrefix, key)
	printf("DELETEDIR: %s", key)
	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpDelete(key, clientv3.WithPrefix())).
		Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *store) Begin(ctx context.Context) st.ITx {
	printf("BEGIN")
	return &tx{
		store:   s,
		context: ctx,
		txn:     s.client.KV.Txn(ctx),
	}
}

func decode(s serializer.Codec, value []byte, out interface{}) error {
	if _, err := converter.EnforcePtr(out); err != nil {
		panic("Error: unable to convert output object to pointer")
	}
	return serializer.Decode(s, value, out)
}

func decodeList(codec serializer.Codec, items []buffer, listOut interface{}) error {
	v, err := converter.EnforcePtr(listOut)
	if err != nil || (v.Kind() != reflect.Slice) {
		panic("Error: need ptr slice")
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

func decodeMap(codec serializer.Codec, items map[string]buffer, mapOut interface{}) error {
	v := reflect.ValueOf(mapOut)
	if v.Kind() != reflect.Map {
		panic("Error: need map")
	}
	for key, item := range items {
		var obj = reflect.New(v.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, item, obj)
		if err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(obj).Elem())
	}
	return nil
}

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
