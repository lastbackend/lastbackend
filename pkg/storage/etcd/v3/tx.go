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
	"github.com/coreos/etcd/clientv3"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"golang.org/x/net/context"
	"path"
)

type tx struct {
	*dbstore
	txn     clientv3.Txn
	context context.Context
	cmp     []clientv3.Cmp
	ops     []clientv3.Op
}

type TxResponse struct {
	*clientv3.TxnResponse
}

//TODO: add compare parameters as argument
func (t *tx) Put(key string, obj interface{}, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)

	log.V(logLevel).Debugf("%s:create:> key: %s, ttl: %d, val: %#v", logPrefix, key, ttl, obj)

	t.cmp = append(t.cmp, clientv3.Compare(clientv3.ModRevision(key), "=", 0))
	data, err := serializer.Encode(t.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> encode data err: %v", logPrefix, err)
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create ttl option err: %v", logPrefix, err)
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

func (t *tx) Set(key string, obj interface{}, ttl uint64, force bool) error {
	key = path.Join(t.pathPrefix, key)

	log.V(logLevel).Debugf("%s:update:> key: %s, ttl: %d, val: %#v", logPrefix, key, ttl, obj)

	t.cmp = append(t.cmp, clientv3.Compare(clientv3.ModRevision(key), "!=", 0))
	data, err := serializer.Encode(t.codec, obj)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> encode data err: %v", logPrefix, err)
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> create ttl option err: %v", logPrefix, err)
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

//TODO: add compare parameters as argument
func (t *tx) Del(key string) {
	key = path.Join(t.pathPrefix, key)

	log.V(logLevel).Debugf("%s:delete:> key: %s", logPrefix, key)

	t.ops = append(t.ops, clientv3.OpDelete(key))
}

func (t *tx) Commit() error {

	log.V(logLevel).Debugf("%s:commit:> commit transaction", logPrefix)

	_, err := t.txn.If(t.cmp...).Then(t.ops...).Commit()
	if err != nil {
		log.V(logLevel).Errorf("%s:commit:> request err: %v", logPrefix, err)
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
