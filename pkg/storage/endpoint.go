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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
)

const endpointStorage = "endpoints"

// Endpoint Service type for interface in interfaces folder
type EndpointStorage struct {
	IEndpoint
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get endpoints by domain name
func (s *EndpointStorage) Get(ctx context.Context, name string) ([]string, error) {

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	endpoints := []string{}
	key := keyCreate(endpointStorage, name)
	if err := client.Map(ctx, key, "", endpoints); err != nil {
		return nil, err
	}

	return endpoints, nil
}

// Insert new endpoint into storage
func (s *EndpointStorage) Insert(ctx context.Context, name string, ips []string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, name)
	if err := client.Create(ctx, key, ips, nil, 0); err != nil {
		return err
	}

	return nil
}

// Update endpoint model
func (s *EndpointStorage) Update(ctx context.Context, name string, ips []string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, name)
	if err := client.Update(ctx, key, ips, nil, 0); err != nil {
		return err
	}

	return nil
}

// Remove endpoint model
func (s *EndpointStorage) Remove(ctx context.Context, name string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, name)
	client.DeleteDir(ctx, key)

	return nil
}

// Watch endpoint model
func (s *EndpointStorage) Watch(ctx context.Context, ips chan []string) error {
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
		if len(keys) < 3 {
			return
		}

		if s, err := s.Get(ctx, keys[1]); err == nil {
			ips <- s
		} else {
			fmt.Println(err)
		}
	}

	client.Watch(ctx, key, filter, cb)
	return nil
}

func newEndpointStorage(config store.Config, util IUtil) *EndpointStorage {
	s := new(EndpointStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
