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

package v3

import (
	"errors"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/store"
	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"golang.org/x/net/context"
)

type dbstore struct {
	store.Store

	debug      bool
	client     *clientv3.Client
	opts       []clientv3.OpOption
	codec      serializer.Codec
	pathPrefix string
	watcher    *watcher
}

// Need for decode array bytes
type buffer []byte

func (s *dbstore) Count(ctx context.Context, key, keyRegexFilter string) (int, error) {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:count:> key: %s with filter: %s", logPrefix, key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:count:> request err: %v", logPrefix, err)
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

func (s *dbstore) Create(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:create:> key: %s, ttl: %d, val: %#v", logPrefix, key, ttl, obj)

	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> encode data err: %v", logPrefix, err)
		return err
	}
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create ttl option err: %v", logPrefix, err)
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, string(data), opts...)).Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> request err: %v", logPrefix, err)
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(types.ErrEntityExists)
	}
	if validator.IsNil(outPtr) {
		log.V(logLevel).Warn("%s:Create: output struct is nil")
		return nil
	} else {
		if err := decode(s.codec, data, outPtr); err != nil {
			log.V(logLevel).Errorf("%s:create:> decode data err: %v", logPrefix, err)
			return err
		}
	}
	return nil
}

func (s *dbstore) Get(ctx context.Context, key string, outPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:get:> key: %s", key)

	res, err := s.client.KV.Get(ctx, key, s.opts...)
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> request err: %v", logPrefix, err)
		return err
	}
	if len(res.Kvs) == 0 {
		return errors.New(types.ErrEntityNotFound)
	}
	if err := decode(s.codec, res.Kvs[0].Value, outPtr); err != nil {
		log.V(logLevel).Errorf("%s:get:> decode data err: %v", logPrefix, err)
		return err
	}
	return nil
}

func (s *dbstore) List(ctx context.Context, key, keyRegexFilter string, listOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:list:> key: %s with filter: %s", logPrefix, key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> request err: %v", logPrefix, err)
		return err
	}

	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]buffer)

	for _, kv := range getResp.Kvs {

		keys := strings.Split(string(kv.Key), "/")
		node := keys[len(keys)-1]

		if (keyRegexFilter != "") && !r.MatchString(string(kv.Key)) {
			continue
		}

		items[node] = kv.Value
	}

	if err := decodeList(s.codec, items, listOutPtr); err != nil {
		log.V(logLevel).Errorf("%s:list:> decode data err: %v", logPrefix, err)
		return err
	}
	return nil
}

func (s *dbstore) Map(ctx context.Context, key, keyRegexFilter string, mapOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:map:> key: %s with filter: %s", logPrefix, key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:map:> request err: %v", logPrefix, err)
		return err
	}

	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]buffer, len(getResp.Kvs))

	for _, kv := range getResp.Kvs {
		if keyRegexFilter != "" && r.MatchString(string(kv.Key)) {
			keys := r.FindStringSubmatch(string(kv.Key))
			items[keys[1]] = buffer(kv.Value)
		} else {
			items[string(kv.Key)] = buffer(kv.Value)
		}
	}

	if len(items) == 0 {
		return errors.New(types.ErrEntityNotFound)
	}

	if err := decodeMap(s.codec, items, mapOutPtr); err != nil {
		log.V(logLevel).Errorf("%s:map:> decode data err: %v", logPrefix, err)
		return err
	}

	return nil
}

func (s *dbstore) MapList(ctx context.Context, key string, keyRegexFilter string, mapOutPtr interface{}) error {

	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:maplist:> key: %s with filter: %s", logPrefix, key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:maplist:> request err: %v", logPrefix, err)
		return err
	}

	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]map[string]buffer)
	for _, kv := range getResp.Kvs {

		if (keyRegexFilter != "") && !r.MatchString(string(kv.Key)) {
			continue
		}

		keys := r.FindStringSubmatch(string(kv.Key))
		field := keys[len(keys)-1]
		node := keys[len(keys)-2]

		if len(items[node]) == 0 {
			items[node] = make(map[string]buffer)
		}

		items[node][field] = kv.Value
	}

	if err := decodeMapList(s.codec, items, mapOutPtr); err != nil {
		log.V(logLevel).Errorf("%s:maplist:> decode data err: %v", logPrefix, err)
		return err
	}
	return nil
}

func (s *dbstore) Update(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:update:> key: %s, ttl: %d, val: %#v", logPrefix, key, ttl, obj)

	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> encode data err: %v", logPrefix, err)
		return err
	}
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> create ttl option err: %v", logPrefix, err)
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "!=", 0)).
		Then(clientv3.OpPut(key, string(data), opts...)).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> request err: %v", logPrefix, err)
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(types.ErrEntityNotFound)
	}
	if validator.IsNil(outPtr) {
		log.V(logLevel).Warnf("%s:Update: output struct is nil", logPrefix)
		return nil
	}
	if outPtr != nil {
		if err := decode(s.codec, data, outPtr); err != nil {
			log.V(logLevel).Errorf("%s:update:> decode data err: %v", logPrefix, err)
			return err
		}
	}
	return nil
}

func (s *dbstore) Upsert(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:upsert:> key: %s, ttl: %d, val: %#v", logPrefix, key, ttl, obj)

	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("%s:upsert:> encode data err: %v", logPrefix, err)
		return err
	}
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("%s:upsert:> create ttl option err: %v", logPrefix, err)
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpPut(key, string(data), opts...)).Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:upsert:> request err: %v", logPrefix, err)
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(types.ErrOperationFailure)
	}
	if validator.IsNil(outPtr) {
		log.V(logLevel).Warn("%s:Upsert: output struct is nil")
		return nil
	}
	if outPtr != nil {
		if err := decode(s.codec, data, outPtr); err != nil {
			log.V(logLevel).Errorf("%s:upsert:> decode data err: %v", logPrefix, err)
			return err
		}
	}
	return nil
}

func (s *dbstore) Delete(ctx context.Context, key string) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:delete:> key: %s", logPrefix, key)

	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpGet(key), clientv3.OpDelete(key)).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:delete:> request err: %v", logPrefix, err)
		return err
	}
	return nil
}

func (s *dbstore) DeleteDir(ctx context.Context, key string) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:deletedir:> key: %s", logPrefix, key)

	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpDelete(key, clientv3.WithPrefix())).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:deletedir:> request err: %v", logPrefix, err)
		return err
	}
	return nil
}

func (s *dbstore) Begin(ctx context.Context) store.TX {

	log.V(logLevel).Debugf("%s:begin:> start transaction", logPrefix)

	t := new(tx)
	t.dbstore = s
	t.context = ctx
	t.txn = s.client.KV.Txn(ctx)
	return t
}

func (s *dbstore) Watch(ctx context.Context, key, keyRegexFilter string) (types.Watcher, error) {
	log.V(logLevel).Debugf("%s:watch:> key: %s, filter: %s", logPrefix, key, keyRegexFilter)
	key = path.Join(s.pathPrefix, key)
	return s.watcher.Watch(ctx, key, keyRegexFilter)
}

func (s *dbstore) Decode(ctx context.Context, value []byte, out interface{}) error {
	return decode(s.codec, value, out)
}

func decode(s serializer.Codec, value []byte, out interface{}) error {
	if _, err := converter.EnforcePtr(out); err != nil {
		return errors.New("unable to convert output struct to pointer")
	}
	return serializer.Decode(s, value, out)
}

func decodeList(codec serializer.Codec, items map[string]buffer, listOut interface{}) error {
	v, err := converter.EnforcePtr(listOut)
	if err != nil || (v.Kind() != reflect.Slice) {
		return errors.New(types.ErrStructOutIsInvalid)
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
	if v.Kind() == reflect.Map {
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

	err := serializer.Decode(codec, joinJSON(items), mapOut)
	if err != nil {
		return err
	}
	return nil
}

func decodeMapList(codec serializer.Codec, items map[string]map[string]buffer, mapOut interface{}) error {
	v := reflect.ValueOf(mapOut)
	if v.Kind() != reflect.Map {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	for key, item := range items {
		var obj = reflect.New(v.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, joinJSON(item), obj)
		if err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(obj).Elem())
		delete(items, key)
	}
	return nil
}

func joinJSON(item map[string]buffer) []byte {
	current := 0
	total := len(item)
	buffer := []byte("{")
	for field, data := range item {
		current++
		buffer = append(buffer, []byte("\"")...)
		buffer = append(buffer, []byte(field)...)
		buffer = append(buffer, []byte("\":")...)
		buffer = append(buffer, data...)
		if current != total {
			buffer = append(buffer, []byte(",")...)
		}
	}
	buffer = append(buffer, []byte("}")...)
	return buffer
}

func (s *dbstore) ttlOpts(ctx context.Context, ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	lcr, err := s.client.Lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(clientv3.LeaseID(lcr.ID))}, nil
}
