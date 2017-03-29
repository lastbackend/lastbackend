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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"time"
)

const VendorTable = "vendors"

// Service User type for interface in interfaces folder
type VendorStorage struct {
	IVendor
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *VendorStorage) Insert(username, vendorUsername, vendorName, vendorHost, serviceID string, token *oauth2.Token) error {
	var (
		err error
		key = fmt.Sprintf("%s/%s/%s/%s", UserTable, username, VendorTable, vendorName)
		vm  *types.Vendor
	)

	vm.Username = vendorUsername
	vm.Vendor = vendorName
	vm.Host = vendorHost
	vm.ServiceID = serviceID
	vm.Token = token

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Create(ctx, key, vm, nil, 0); err != nil {
		return err
	}

	return nil
}

func (s *VendorStorage) Get(username, vendorName string) (*types.Vendor, error) {
	var (
		vendor = new(types.Vendor)
		key    = fmt.Sprintf("%s/%s/%s/%s", UserTable, username, VendorTable, vendorName)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, key, vendor); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return vendor, nil
}

func (s *VendorStorage) List(username string) (*types.VendorItems, error) {
	var (
		vendorItems = new(types.VendorItems)
		key         = fmt.Sprintf("%s/%s/%s", UserTable, username, VendorTable)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.List(ctx, key, ``, vendorItems); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return vendorItems, nil
}

func (s *VendorStorage) Remove(username, vendorName string) error {
	var (
		err error
		key = fmt.Sprintf("%s/%s/%s/%s", UserTable, username, VendorTable, vendorName)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Delete(ctx, key, nil); err != nil {
		return err
	}

	return nil
}

func newVendorStorage(config store.Config) *VendorStorage {
	s := new(VendorStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
