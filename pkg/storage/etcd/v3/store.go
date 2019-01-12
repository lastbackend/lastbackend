//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/coreos/etcd/mvcc/mvccpb"
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

func (s *dbstore) Info(ctx context.Context, key string) (*types.Runtime, error) {

	key = path.Join(s.pathPrefix, key)
	r := new(types.Runtime)

	log.V(logLevel).Debugf("%s:count:> key: %s with filter: %s", logPrefix, key)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:count:> request err: %v", logPrefix, err)
		return r, err
	}

	r.System.Revision = getResp.Header.Revision
	return r, nil
}

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

func (s *dbstore) Put(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {

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

	if err := setEntityRuntimeInfo(outPtr, getRuntimeFromResponse(txnResp.Header)); err != nil {
		log.V(logLevel).Errorf("%s:get:> can not set runtime info err: %v", logPrefix, err)
		return err
	}

	return nil
}

func (s *dbstore) Get(ctx context.Context, key string, outPtr interface{}, rev *int64) error {

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

	if err := setEntityRuntimeInfo(outPtr, getRuntimeFromResponse(res.Header)); err != nil {
		log.V(logLevel).Errorf("%s:get:> can not set runtime info err: %v", logPrefix, err)
		return err
	}

	return nil
}

func (s *dbstore) List(ctx context.Context, key, keyRegexFilter string, listOutPtr interface{}, rev *int64) error {

	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:list:> key: %s with filter: %s", logPrefix, key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> request err: %v", logPrefix, err)
		return err
	}

	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]*mvccpb.KeyValue)

	for _, kv := range getResp.Kvs {

		keys := strings.Split(string(kv.Key), "/")
		node := keys[len(keys)-1]

		if (keyRegexFilter != "") && !r.MatchString(string(kv.Key)) {
			continue
		}

		items[node] = kv
	}

	if err := decodeList(s.codec, items, listOutPtr); err != nil {
		log.V(logLevel).Errorf("%s:list:> decode data err: %v", logPrefix, err)
		return err
	}

	if err := setEntityRuntimeInfo(listOutPtr, getRuntimeFromResponse(getResp.Header)); err != nil {
		log.V(logLevel).Errorf("%s:get:> can not set runtime info err: %v", logPrefix, err)
		return err
	}

	return nil
}

func (s *dbstore) Map(ctx context.Context, key, keyRegexFilter string, mapOutPtr interface{}, rev *int64) error {

	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:map:> key: %s with filter: %s", logPrefix, key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("%s:map:> request err: %v", logPrefix, err)
		return err
	}

	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]*mvccpb.KeyValue, len(getResp.Kvs))

	for _, kv := range getResp.Kvs {

		if keyRegexFilter != "" && r.MatchString(string(kv.Key)) {
			keys := r.FindStringSubmatch(string(kv.Key))
			items[keys[1]] = kv
		} else {
			items[string(kv.Key)] = kv
		}
	}

	if len(items) == 0 {
		return nil
	}

	if err := decodeMap(s.codec, items, mapOutPtr); err != nil {
		log.V(logLevel).Errorf("%s:map:> decode data err: %v", logPrefix, err)
		return err
	}

	if err := setEntityRuntimeInfo(mapOutPtr, getRuntimeFromResponse(getResp.Header)); err != nil {
		log.V(logLevel).Errorf("%s:get:> can not set runtime info err: %v", logPrefix, err)
		return err
	}

	return nil
}

func (s *dbstore) Set(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64, force bool, rev *int64) error {

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

	txn := s.client.KV.Txn(ctx)

	if !force {
		rv := int64(0)
		if rev != nil {
			rv = *rev
		}

		txn = txn.If(clientv3.Compare(clientv3.ModRevision(key), "!=", rv))
	}

	txnResp, err := txn.
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

	if err := setEntityRuntimeInfo(outPtr, getRuntimeFromResponse(txnResp.Header)); err != nil {
		log.V(logLevel).Errorf("%s:get:> can not set runtime info err: %v", logPrefix, err)
		return err
	}

	return nil
}

func (s *dbstore) Del(ctx context.Context, key string) error {

	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("%s:delete:> key: %s", logPrefix, key)

	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpGet(key), clientv3.OpDelete(key, clientv3.WithPrefix())).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:delete:> request err: %v", logPrefix, err)
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

func (s *dbstore) Watch(ctx context.Context, key, keyRegexFilter string, rev *int64) (types.Watcher, error) {
	log.V(logLevel).Debugf("%s:watch:> key: %s, filter: %s", logPrefix, key, keyRegexFilter)
	key = path.Join(s.pathPrefix, key)
	return s.watcher.Watch(ctx, key, keyRegexFilter, rev)
}

func (s *dbstore) Decode(ctx context.Context, value []byte, out interface{}) error {
	return decode(s.codec, value, out)
}

func setEntityRuntimeInfo(out interface{}, runtime types.Runtime) error {

	var (
		v   reflect.Value
		err error
	)

	if reflect.TypeOf(out).Kind() == reflect.Ptr {
		v, err = converter.EnforcePtr(out)
		if err != nil {
			return errors.New("unable to convert output struct to pointer")
		}
	} else {
		v = reflect.ValueOf(out)
	}

	setValueRuntimeInfo(v, runtime)
	return nil
}

func setValueRuntimeInfo(v reflect.Value, runtime types.Runtime) error {

	if v.Kind() != reflect.Struct {
		return nil
	}

	f := v.FieldByName("Runtime")

	if !f.IsValid() {
		return nil
	}

	if !f.CanSet() {
		return nil
	}

	f.Set(reflect.ValueOf(runtime.Runtime))

	return nil
}

func getRuntimeFromResponse(res *etcdserverpb.ResponseHeader) types.Runtime {
	runtime := types.Runtime{}
	runtime.System.Revision = res.Revision
	return runtime
}

func getRuntimeFromValue(res *mvccpb.KeyValue) types.Runtime {
	runtime := types.Runtime{}
	runtime.System.Revision = res.ModRevision
	return runtime
}

func decode(s serializer.Codec, value []byte, out interface{}) error {
	if _, err := converter.EnforcePtr(out); err != nil {
		return errors.New("unable to convert output struct to pointer")
	}
	return serializer.Decode(s, value, out)
}

func decodeList(codec serializer.Codec, items map[string]*mvccpb.KeyValue, listOut interface{}) error {
	v, err := converter.EnforcePtr(listOut)
	if err != nil {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	f := v.FieldByName("Items")
	if f.Kind() != reflect.Slice {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	if !f.IsValid() {
		return nil
	}

	if !f.CanSet() {
		return nil
	}

	for _, item := range items {

		var obj = reflect.New(f.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, item.Value, obj)
		if err != nil {
			return err
		}

		setValueRuntimeInfo(reflect.Indirect(reflect.ValueOf(obj).Elem()), getRuntimeFromValue(item))
		f.Set(reflect.Append(f, reflect.ValueOf(obj).Elem()))
	}

	return nil
}

func decodeMap(codec serializer.Codec, items map[string]*mvccpb.KeyValue, mapOut interface{}) error {

	v, err := converter.EnforcePtr(mapOut)
	if err != nil {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	f := v.FieldByName("Items")
	if f.Kind() != reflect.Map {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	if !f.IsValid() {
		return nil
	}

	if !f.CanSet() {
		return nil
	}

	for key, item := range items {
		var obj = reflect.New(f.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, item.Value, obj)
		if err != nil {
			return err
		}

		setValueRuntimeInfo(reflect.Indirect(reflect.ValueOf(obj).Elem()), getRuntimeFromValue(item))
		f.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(obj).Elem())
	}

	return nil

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
