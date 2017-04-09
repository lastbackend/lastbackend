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
	"fmt"
	"github.com/coreos/etcd/clientv3"
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
func (t *tx) Create(key string, objPtr interface{}, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)
	fmt.Println("Create", key)
	t.cmp = append(t.cmp, clientv3.Compare(clientv3.ModRevision(key), "=", 0))
	data, err := serializer.Encode(t.codec, objPtr)
	if err != nil {
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

func (t *tx) Update(key string, objPtr interface{}, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)
	fmt.Println("Update", key)
	t.cmp = append(t.cmp, clientv3.Compare(clientv3.ModRevision(key), "!=", 0))
	data, err := serializer.Encode(t.codec, objPtr)
	if err != nil {
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	if err != nil {
		return err
	}
	t.ops = append(t.ops, clientv3.OpPut(key, string(data), opts...))
	return nil
}

//TODO: add compare parameters as argument
func (t *tx) Delete(key string) {
	key = path.Join(t.pathPrefix, key)
	fmt.Println("Delete", key)
	t.ops = append(t.ops, clientv3.OpDelete(key))
}

//TODO: add compare parameters as argument
func (t *tx) DeleteDir(key string) {
	key = path.Join(t.pathPrefix, key)
	fmt.Println("Delete", key)
	t.ops = append(t.ops, clientv3.OpDelete(key, clientv3.WithPrefix()))
}

func (t *tx) Commit() error {
	_, err := t.txn.If(t.cmp...).Then(t.ops...).Commit()
	return err
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
