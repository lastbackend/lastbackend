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
	st "github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"golang.org/x/net/context"
)

type store struct {
	debug      bool
	client     *clientv3.Client
	opts       []clientv3.OpOption
	codec      serializer.Codec
	pathPrefix string
}

const (
	logLevel = 5
)

// Need for decode array bytes
type buffer []byte

func (s *store) WatchClose() {
	s.client.Watcher.Close()
}

func (s *store) Count(ctx context.Context, key, keyRegexFilter string) (int, error) {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Count: key: %s with filter: %s", key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Count: request err: %s", err.Error())
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
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Create: key: %s, ttl: %d, val: %#v", key, ttl, obj)

	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Create: encode data err: %s", err.Error())
		return err
	}
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Create: create ttl option err: %s", err.Error())
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, string(data), opts...)).Commit()
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Create: request err: %s", err.Error())
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrEntityExists)
	}
	if validator.IsNil(outPtr) {
		log.V(logLevel).Warn("Etcd3: Create: output struct is nil")
		return nil
	} else {
		if err := decode(s.codec, data, outPtr); err != nil {
			log.V(logLevel).Errorf("Etcd3: Create: decode data err: %s", err.Error())
			return err
		}
	}
	return nil
}

func (s *store) Get(ctx context.Context, key string, outPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Get: key: %s", key)

	res, err := s.client.KV.Get(ctx, key, s.opts...)
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Get: request err: %s", err.Error())
		return err
	}
	if len(res.Kvs) == 0 {
		return errors.New(st.ErrEntityNotFound)
	}
	if err := decode(s.codec, res.Kvs[0].Value, outPtr); err != nil {
		log.V(logLevel).Errorf("Etcd3: Get: decode data err: %s", err.Error())
		return err
	}
	return nil
}

func (s *store) List(ctx context.Context, key, keyRegexFilter string, listOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: List: key: %s with filter: %s", key, keyRegexFilter)

	if !strings.HasSuffix(key, "/") {
		key += "/"
	}

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: List: request err: %s", err.Error())
		return err
	}

	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]map[string]buffer)
	for _, kv := range getResp.Kvs {

		keys := strings.Split(string(kv.Key), "/")
		node := keys[len(keys)-2]
		field := keys[len(keys)-1]

		if (keyRegexFilter != "") && !r.MatchString(string(kv.Key)) {
			continue
		}
		if len(items[node]) == 0 {
			items[node] = make(map[string]buffer)
		}
		items[node][field] = kv.Value
	}

	if err := decodeList(s.codec, items, listOutPtr); err != nil {
		log.V(logLevel).Errorf("Etcd3: List: decode data err: %s", err.Error())
		return err
	}
	return nil
}

func (s *store) Map(ctx context.Context, key, keyRegexFilter string, mapOutPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Map: key: %s with filter: %s", key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Map: request err: %s", err.Error())
		return err
	}
	r, _ := regexp.Compile(keyRegexFilter)
	items := make(map[string]buffer, len(getResp.Kvs))
	for _, kv := range getResp.Kvs {
		if (keyRegexFilter == "") || r.MatchString(string(kv.Key)) {
			keys := r.FindStringSubmatch(string(kv.Key))
			items[keys[1]] = buffer(kv.Value)
		}
	}

	if len(items) == 0 {
		return errors.New(st.ErrEntityNotFound)
	}

	if err := decodeMap(s.codec, items, mapOutPtr); err != nil {
		log.V(logLevel).Errorf("Etcd3: Map: decode data err: %s", err.Error())
		return err
	}

	return nil
}

func (s *store) MapList(ctx context.Context, key string, keyRegexFilter string, mapOutPtr interface{}) error {

	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: MapList: key: %s with filter: %s", key, keyRegexFilter)

	getResp, err := s.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: MapList: request err: %s", err.Error())
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
		log.V(logLevel).Errorf("Etcd3: MapList: decode data err: %s", err.Error())
		return err
	}
	return nil
}

func (s *store) Update(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Update: key: %s, ttl: %d, val: %#v", key, ttl, obj)

	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Update: encode data err: %s", err.Error())
		return err
	}
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Update: create ttl option err: %s", err.Error())
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "!=", 0)).
		Then(clientv3.OpPut(key, string(data), opts...)).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Update: request err: %s", err.Error())
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrEntityNotFound)
	}
	if validator.IsNil(outPtr) {
		log.V(logLevel).Warn("Etcd3: Update: output struct is nil")
		return nil
	}
	if outPtr != nil {
		if err := decode(s.codec, data, outPtr); err != nil {
			log.V(logLevel).Errorf("Etcd3: Update: decode data err: %s", err.Error())
			return err
		}
	}
	return nil
}

func (s *store) Upsert(ctx context.Context, key string, obj, outPtr interface{}, ttl uint64) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Upsert: key: %s, ttl: %d, val: %#v", key, ttl, obj)

	data, err := serializer.Encode(s.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Upsert: encode data err: %s", err.Error())
		return err
	}
	opts, err := s.ttlOpts(ctx, int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Upsert: create ttl option err: %s", err.Error())
		return err
	}
	txnResp, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpPut(key, string(data), opts...)).Commit()
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Upsert: request err: %s", err.Error())
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrEntityExists)
	}
	if validator.IsNil(outPtr) {
		log.V(logLevel).Warn("Etcd3: Upsert: output struct is nil")
		return nil
	}
	if outPtr != nil {
		if err := decode(s.codec, data, outPtr); err != nil {
			log.V(logLevel).Errorf("Etcd3: Upsert: decode data err: %s", err.Error())
			return err
		}
	}
	return nil
}

func (s *store) Delete(ctx context.Context, key string) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: Delete: key: %s", key)

	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpGet(key), clientv3.OpDelete(key)).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: Delete: request err: %s", err.Error())
		return err
	}
	return nil
}

func (s *store) DeleteDir(ctx context.Context, key string) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: DeleteDir: key: %s", key)

	_, err := s.client.KV.Txn(ctx).
		Then(clientv3.OpDelete(key, clientv3.WithPrefix())).
		Commit()
	if err != nil {
		log.V(logLevel).Errorf("Etcd3: DeleteDir: request err: %s", err.Error())
		return err
	}
	return nil
}

func (s *store) Begin(ctx context.Context) st.TX {

	log.V(logLevel).Debugf("Etcd3: Begin")

	t := new(tx)
	t.store = s
	t.context = ctx
	t.txn = s.client.KV.Txn(ctx)
	return t
}

func (s *store) Watch(ctx context.Context, key, keyRegexFilter string, f func(string, string, []byte)) error {
	key = path.Join(s.pathPrefix, key)

	log.V(logLevel).Debugf("Etcd3: WatchService: key: %s, filter: %s", key, keyRegexFilter)

	r, _ := regexp.Compile(keyRegexFilter)
	rch := s.client.Watch(context.Background(), key, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if r.MatchString(string(ev.Kv.Key)) {

				action := EventTypeCreate

				if ev.Type.String() == "PUT" && wresp.Header.Revision != ev.Kv.CreateRevision {
					action = EventTypeUpdate
				}

				if ev.Type.String() == "DEL" {
					action = EventTypeDelete
				}

				go f(action, string(ev.Kv.Key), ev.Kv.Value)
			}
		}
	}

	return nil
}

func (s *store) Decode(ctx context.Context, value []byte, out interface{}) error {
	return decode(s.codec, value, out)
}

func decode(s serializer.Codec, value []byte, out interface{}) error {
	if _, err := converter.EnforcePtr(out); err != nil {
		panic("Error: unable to convert output struct to pointer")
	}
	return serializer.Decode(s, value, out)
}

func decodeList(codec serializer.Codec, items map[string]map[string]buffer, listOut interface{}) error {
	v, err := converter.EnforcePtr(listOut)
	if err != nil || (v.Kind() != reflect.Slice) {
		panic("Error: need ptr slice")
	}

	for _, item := range items {
		var obj = reflect.New(v.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, joinJSON(item), obj)
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
		panic("Error: need map")
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

func GetStore(client *clientv3.Client, opts []clientv3.OpOption, codec serializer.Codec, pathPrefix string, debug bool) *store {
	return &store{
		client:     client,
		opts:       opts,
		codec:      codec,
		pathPrefix: pathPrefix,
		debug:      debug,
	}
}
