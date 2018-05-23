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

package etcd

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const (
	ipamStorage   = "utils/ipam"
	logIPAMPrefix = "storage:etcd:ipam:>"
)

type IPAMStorage struct {
	storage.IPAM
}

func (s *IPAMStorage) Get(ctx context.Context) ([]string, error) {
	log.V(logLevel).Debugf("%s get ipam context", logIPAMPrefix)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("%s create client err: %s", logIPAMPrefix, err.Error())
		return nil, err
	}
	defer destroy()

	var ips = make([]string, 0)
	if err := client.Get(ctx, ipamStorage, &ips); err != nil {
		if err.Error() == store.ErrEntityNotFound {
			return make([]string, 0), nil
		}
		log.V(logLevel).Errorf("%s get nodes list err: %s", logIPAMPrefix, err.Error())
		return nil, err
	}

	return ips, nil
}

func (s *IPAMStorage) Set(ctx context.Context, ips []string) error {

	log.V(logLevel).Debugf("%s set ipam context", logIPAMPrefix)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("%s create client err: %s", logIPAMPrefix, err.Error())
		return err
	}
	defer destroy()

	if err := client.Upsert(ctx, ipamStorage, ips, nil, 0); err != nil {
		log.V(logLevel).Errorf("%s set ipam context err: %s", logIPAMPrefix, err.Error())
		return err
	}

	return nil
}

func newIPAMStorage() *IPAMStorage {
	s := new(IPAMStorage)
	return s
}
