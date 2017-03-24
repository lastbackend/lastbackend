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
	"github.com/lastbackend/lastbackend/pkg/serializer"
	"golang.org/x/net/context"
	"path"
)

type tx struct {
	*store
	txn     clientv3.Txn
	context context.Context
}

type TxResponse struct {
	*clientv3.TxnResponse
}

// Commit transaction context
func (t *tx) Create(key string, obj, ttl uint64) error {
	key = path.Join(t.pathPrefix, key)
	data, err := serializer.Encode(t.codec, obj)
	if err != nil {
		return err
	}
	opts, err := t.ttlOpts(int64(ttl))
	t.txn = t.txn.Then(clientv3.OpPut(key, string(data), opts...))
	return nil
}

// Delete key transaction context
func (t *tx) Delete(key string) {
	key = path.Join(t.pathPrefix, key)
	t.txn = t.txn.Then(clientv3.OpDelete(key))
}

// Commit transaction context
func (t *tx) Commit() (*TxResponse, error) {
	resp, err := t.txn.Commit()
	if err != nil {
		return nil, err
	}
	return &TxResponse{resp}, nil
}

// ttlOpts returns client options based on given ttl.
// ttl: if ttl is non-zero, it will attach the key to a lease with ttl of roughly the same length
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
