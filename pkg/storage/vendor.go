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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
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

func (s *VendorStorage) Insert(ctx context.Context, username, vendorUsername, vendorName, vendorHost, serviceID string, token *oauth2.Token) error {
	var (
		err error
		// Key example: /users/<username>/vendors/<vendor>
		key = fmt.Sprintf("%s/%s/%s/%s", UserTable, username, VendorTable, vendorName)
		vm  = new(types.Vendor)
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

	err = client.Get(ctx, key, vm)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		return err
	}
	if err != nil && err.Error() == store.ErrKeyNotFound {
		return client.Create(ctx, key, vm, nil, 0)
	}

	return client.Update(ctx, key, vm, nil, 0)
}

func (s *VendorStorage) Get(ctx context.Context, username, vendorName string) (*types.Vendor, error) {
	var (
		vendor = new(types.Vendor)
		// Key example: /users/<username>/vendors/<vendor>
		key = fmt.Sprintf("%s/%s/%s/%s", UserTable, username, VendorTable, vendorName)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, key, &vendor); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return vendor, nil
}

func (s *VendorStorage) List(ctx context.Context, username string) (map[string]*types.Vendor, error) {
	var (
		vendors = make(map[string]*types.Vendor)
		// Key example: /users/<username>/vendors
		key = fmt.Sprintf("%s/%s/%s", UserTable, username, VendorTable)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Map(ctx, key, ``, vendors); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return vendors, nil
}

func (s *VendorStorage) Remove(ctx context.Context, username, vendorName string) error {
	var (
		err error
		// Key example: /users/<username>/vendors/<vendor>
		key = fmt.Sprintf("%s/%s/%s/%s", UserTable, username, VendorTable, vendorName)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Delete(ctx, key, nil)
	if err != nil && err.Error() == store.ErrKeyNotFound {
		return nil
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
