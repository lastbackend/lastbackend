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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
	"golang.org/x/oauth2"
)

const vendorStorage = "vendors"

// Service vendor type for interface in interfaces folder
type VendorStorage struct {
	IVendor
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *VendorStorage) Insert(ctx context.Context, owner, name, host, serviceID string, token *oauth2.Token) error {

	vm := new(types.Vendor)
	vm.Username = owner
	vm.Vendor = name
	vm.Host = host
	vm.ServiceID = serviceID
	vm.Token = token

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, vendorStorage, name)
	err = client.Get(ctx, key, vm)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		return err
	}
	if err != nil && err.Error() == store.ErrKeyNotFound {
		return client.Create(ctx, key, vm, nil, 0)
	}

	return client.Update(ctx, key, vm, nil, 0)
}

func (s *VendorStorage) Get(ctx context.Context, vendorName string) (types.Vendor, error) {

	vendor := types.Vendor{}
	client, destroy, err := s.Client()
	if err != nil {
		return vendor, err
	}
	defer destroy()

	key := s.util.Key(ctx, vendorStorage, vendorName)

	if err := client.Get(ctx, key, &vendor); err != nil {
		return vendor, err
	}

	return vendor, nil
}

func (s *VendorStorage) List(ctx context.Context) (map[string]types.Vendor, error) {

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, vendorStorage)
	vendors := make(map[string]types.Vendor)
	if err := client.Map(ctx, key, ``, vendors); err != nil {
		return vendors, err
	}

	return vendors, nil
}

func (s *VendorStorage) Remove(ctx context.Context, vendorName string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, vendorStorage, vendorName)
	err = client.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func newVendorStorage(config store.Config, util IUtil) *VendorStorage {
	s := new(VendorStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
