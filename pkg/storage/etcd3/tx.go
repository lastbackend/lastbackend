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
	"github.com/coreos/etcd/clientv3"
	st "github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"golang.org/x/net/context"
	"path"
)

type tx struct {
	*store
	txn     clientv3.Txn
	context context.Context
	cmp     []clientv3.Cmp
	ops     []clientv3.Op
}

type TxResponse struct {
	*clientv3.TxnResponse
}

//TODO: add compare parameters as argument
func (t *tx) Create(key string, obj interface{}, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)

	t.log.V(st.DebugLevel).Debugf("Etcd3: Create: key: %s, ttl: %d, val: %#v", key, ttl, obj)

	t.cmp = append(t.cmp, clientv3.Compare(clientv3.ModRevision(key), "=", 0))
	data, err := serializer.Encode(t.codec, obj)
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Create: encode data err: %s", err.Error())
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Create: create ttl option err: %s", err.Error())
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

func (t *tx) Update(key string, obj interface{}, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)

	t.log.V(st.DebugLevel).Debugf("Etcd3: Update: key: %s, ttl: %d, val: %#v", key, ttl, obj)

	t.cmp = append(t.cmp, clientv3.Compare(clientv3.ModRevision(key), "!=", 0))
	data, err := serializer.Encode(t.codec, obj)
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Update: encode data err: %s", err.Error())
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Update: create ttl option err: %s", err.Error())
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

func (t *tx) Upsert(key string, obj interface{}, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)

	t.log.V(st.DebugLevel).Debugf("Etcd3: Upsert: key: %s, val: %#v", key, obj)

	data, err := serializer.Encode(t.codec, obj)
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Upsert: encode data err: %s", err.Error())
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Upsert: create ttl option err: %s", err.Error())
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

//TODO: add compare parameters as argument
func (t *tx) Delete(key string) {
	key = path.Join(t.pathPrefix, key)

	t.log.V(st.DebugLevel).Debugf("Etcd3: Delete: key: %s", key)

	t.ops = append(t.ops, clientv3.OpDelete(key))
}

//TODO: add compare parameters as argument
func (t *tx) DeleteDir(key string) {
	key = path.Join(t.pathPrefix, key)

	t.log.V(st.DebugLevel).Debugf("Etcd3: DeleteDir: key: %s", key)

	t.ops = append(t.ops, clientv3.OpDelete(key, clientv3.WithPrefix()))
}

func (t *tx) Commit() error {

	t.log.V(st.DebugLevel).Debugf("Etcd3: Commit")

	_, err := t.txn.If(t.cmp...).Then(t.ops...).Commit()
	if err != nil {
		t.log.V(st.DebugLevel).Errorf("Etcd3: Commit: request err: %s", err.Error())
		return err
	}
	return nil
}

func (t *tx) ttlOpts(ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	lcr, err := t.client.Lease.Grant(t.context, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(clientv3.LeaseID(lcr.ID))}, nil
}
