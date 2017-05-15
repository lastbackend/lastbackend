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
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const systemStorage = "system"
const systemLeadTTL = 11

// Namespace Service type for interface in interfaces folder
type SystemStorage struct {
	ISystem
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *SystemStorage) ProcessSet(ctx context.Context, process *types.Process) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMember := keyCreate(systemStorage, process.Meta.Kind, "process", process.Meta.Hostname)
	err = client.Upsert(ctx, keyMember, process, nil, systemLeadTTL)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemStorage) Elect(ctx context.Context, process *types.Process) (bool, error) {

	var (
		id   string
		err  error
		lead bool
	)

	client, destroy, err := s.Client()
	if err != nil {
		return lead, err
	}
	defer destroy()

	key := keyCreate(systemStorage, process.Meta.Kind, "lead")
	err = client.Get(ctx, key, &id)
	if err != nil && (err.Error() != store.ErrKeyNotFound) {
		return lead, err
	}

	if id != "" {
		return lead, nil
	}

	if err.Error() == store.ErrKeyNotFound {
		err = client.Create(ctx, key, &process.Meta.ID, nil, systemLeadTTL)
		if err != nil && err.Error() != store.ErrKeyExists {
			return lead, err
		}
		lead = true
	}

	return lead, nil
}

func (s *SystemStorage) ElectUpdate(ctx context.Context, process *types.Process) error {

	var (
		id  string
		err error
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := keyCreate(systemStorage, process.Meta.Kind, "lead")
	err = client.Get(ctx, key, &id)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		return err
	}

	if id != process.Meta.ID {
		return nil
	}

	err = client.Update(ctx, key, &process.Meta.ID, nil, systemLeadTTL)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemStorage) ElectWait(ctx context.Context, process *types.Process, lead chan bool) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := keyCreate(systemStorage, process.Meta.Kind, "lead")

	cb := func(action, key string, val []byte) {

		if action == "PUT" {

			var id string
			if err := json.Unmarshal(val, &id); err != nil {
				//TODO: return error and start loop over
			}

			if id == process.Meta.ID {
				lead <- true
			} else {
				lead <- false
			}
		}

		if action == "DELETE" {
			s.Elect(ctx, process)
		}
	}

	client.Watch(ctx, key, "", cb)

	return nil
}

func newSystemStorage(config store.Config, util IUtil) *SystemStorage {
	s := new(SystemStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
