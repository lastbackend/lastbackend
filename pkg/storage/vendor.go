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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/oauth2"
)

const vendorStorage = "vendors"

// Service vendor type for interface in interfaces folder
type VendorStorage struct {
	IVendor
	log    logger.ILogger
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *VendorStorage) Insert(ctx context.Context, owner, name, host, serviceID string, token *oauth2.Token) error {

	s.log.V(logLevel).Debugf("Storage: Vendor: insert vendor owner: %s, name: %s, host: %s, serviceID: %s, token: %#v",
		owner, name, host, serviceID, token)

	vm := new(types.Vendor)
	vm.Username = owner
	vm.Vendor = name
	vm.Host = host
	vm.ServiceID = serviceID
	vm.Token = token

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(vendorStorage, name)
	if err := client.Upsert(ctx, key, vm, nil, 0); err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: upsert vendor info err: %s", err.Error())
		return err
	}
	return nil
}

func (s *VendorStorage) Get(ctx context.Context, name string) (*types.Vendor, error) {

	s.log.V(logLevel).Debugf("Storage: Vendor: get by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		s.log.V(logLevel).Errorf("Storage: Vendor: get name err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	vendor := new(types.Vendor)
	key := keyCreate(vendorStorage, name)
	if err := client.Get(ctx, key, vendor); err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: get vendor info err: %s", err.Error())
		return nil, err
	}

	return vendor, nil
}

func (s *VendorStorage) List(ctx context.Context) (map[string]*types.Vendor, error) {

	s.log.V(logLevel).Debug("Storage: Vendor: get vendors list")

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(vendorStorage)
	vendors := make(map[string]*types.Vendor)
	if err := client.Map(ctx, key, ``, vendors); err != nil && err.Error() != store.ErrKeyNotFound {
		s.log.V(logLevel).Errorf("Storage: Vendor: map vendors err: %s", err.Error())
		return nil, err
	}

	return vendors, nil
}

func (s *VendorStorage) Remove(ctx context.Context, name string) error {

	s.log.V(logLevel).Debugf("Storage: Vendor: remove vendor by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		s.log.V(logLevel).Errorf("Storage: Vendor: remove vendor by name err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(vendorStorage, name)
	if err := client.Delete(ctx, key); err != nil {
		s.log.V(logLevel).Errorf("Storage: Vendor: delete vendor err: %s", err.Error())
		return err
	}

	return nil
}

func newVendorStorage(config store.Config, log logger.ILogger) *VendorStorage {
	s := new(VendorStorage)
	s.log = log
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
