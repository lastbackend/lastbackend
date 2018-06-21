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

package mock

import (
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type tx struct {
	*mockstore
	txn     clientv3.Txn
	context context.Context
	cmp     []clientv3.Cmp
	ops     []clientv3.Op
}

type TxResponse struct {
	*clientv3.TxnResponse
}

func (t *tx) Create(key string, obj interface{}, ttl uint64) error {
	return nil
}

func (t *tx) Update(key string, obj interface{}, ttl uint64) error {
	return nil
}

func (t *tx) Upsert(key string, obj interface{}, ttl uint64) error {
	return nil
}

func (t *tx) Delete(key string) {
}

func (t *tx) DeleteDir(key string) {
}

func (t *tx) Commit() error {
	return nil
}