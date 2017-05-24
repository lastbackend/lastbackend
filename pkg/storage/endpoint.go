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
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
)

const endpointStorage = "endpoints"

// Endpoint Service type for interface in interfaces folder
type EndpointStorage struct {
	IEndpoint
	log    logger.ILogger
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get endpoints by domain name
func (s *EndpointStorage) Get(ctx context.Context, name string) ([]string, error) {

	s.log.V(logLevel).Debugf("Storage: Endpoint: get endpoint by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be nil")
		s.log.V(logLevel).Errorf("Storage: Endpoint: get endpoint by name err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	endpoints := []string{}
	key := keyCreate(endpointStorage, name)
	if err := client.Get(ctx, key, &endpoints); err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: get endpoint err: %s", err.Error())
		return nil, err
	}

	return endpoints, nil
}

// Update endpoint model
func (s *EndpointStorage) Upsert(ctx context.Context, name string, ips []string) error {

	s.log.V(logLevel).Debugf("Storage: Endpoint: update endpoint by name: %s with ips: %#v", name, ips)

	if len(name) == 0 {
		err := errors.New("name can not be nil")
		s.log.V(logLevel).Errorf("Storage: Endpoint: update endpoint err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, name)
	if err := client.Upsert(ctx, key, ips, nil, 0); err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: upsert endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Remove endpoint model
func (s *EndpointStorage) Remove(ctx context.Context, name string) error {

	s.log.V(logLevel).Debugf("Storage: Endpoint: remove endpoint by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be nil")
		s.log.V(logLevel).Errorf("Storage: Endpoint: remove endpoint err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, name)
	if err := client.DeleteDir(ctx, key); err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: delete dir endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// WatchSetvice endpoint model
func (s *EndpointStorage) Watch(ctx context.Context, endpoint chan string) error {

	s.log.V(logLevel).Debug("Storage: Endpoint: watch endpoint")

	const filter = `\b.+` + endpointStorage + `\/(.+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(endpointStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}
		endpoint <- keys[1]
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		s.log.V(logLevel).Errorf("Storage: Endpoint: watch endpoint err: %s", err.Error())
		return err
	}

	return nil
}

func newEndpointStorage(config store.Config, log logger.ILogger, util IUtil) *EndpointStorage {
	s := new(EndpointStorage)
	s.log = log
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
