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
	"github.com/lastbackend/lastbackend/pkg/serializer"
	st "github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"golang.org/x/net/context"
	"path"
)

type store struct {
	client *clientv3.Client
	// getOpts contains additional options that should be passed to all Get() calls.
	getOps     []clientv3.OpOption
	codec      serializer.Codec
	pathPrefix string
}

// Create implements store.Interface.Create.
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
func (s *store) Get(ctx context.Context, key string, objPtr interface{}) error {
	key = path.Join(s.pathPrefix, key)
	res, err := s.client.KV.Get(ctx, key, s.getOps...)
	if err != nil {
		return err
	}
	if len(res.Kvs) == 0 {
		return nil
	}
	return decode(s.codec, res.Kvs[0].Value, objPtr)
}

// Delete implements store.Interface.Delete.
func (s *store) Delete(ctx context.Context, key string, out interface{}) error {
	key = path.Join(s.pathPrefix, key)
	txnResp, err := s.client.KV.Txn(ctx).If(notFound(key)).
		Then(clientv3.OpDelete(key)).
		Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return errors.New(st.ErrOperationFailure)
	}
	return nil
}

// Decode decodes value of bytes into object.
// On success, objPtr would be set to the object.
func decode(s serializer.Codec, value []byte, objPtr interface{}) error {
	if _, err := converter.EnforcePtr(objPtr); err != nil {
		panic("unable to convert output object to pointer")
	}
	return serializer.Decode(s, value, objPtr)
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
