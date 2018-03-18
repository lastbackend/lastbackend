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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const triggerStorage string = "triggers"

// Service Trigger type for interface in interfaces folder
type TriggerStorage struct {
	storage.Trigger
}

// Get triggers by id
func (s *TriggerStorage) Get(ctx context.Context, namespace, service, name string) (*types.Trigger, error) {

	log.V(logLevel).Debugf("Storage: Trigger: get trigger by name: %s", name)

	if len(name) == 0 {
		err := errors.New("id can not be empty")
		log.V(logLevel).Errorf("Storage: Trigger: get trigger by id err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Trigger: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	trigger := new(types.Trigger)
	keyMeta := keyCreate(triggerStorage, name)
	if err := client.Get(ctx, keyMeta, &trigger); err != nil {
		log.V(logLevel).Errorf("Storage: Trigger: get trigger meta err: %s", err.Error())
		return nil, err
	}

	return trigger, nil
}

// Insert new trigger into storage
func (s *TriggerStorage) Insert(ctx context.Context, trigger *types.Trigger) error {

	log.V(logLevel).Debugf("Storage: Trigger: create trigger: %#v", trigger)

	if trigger == nil {
		err := errors.New("trigger can not be nil")
		log.V(logLevel).Errorf("Storage: Trigger: create trigger err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Trigger: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(triggerStorage, trigger.Meta.Name)
	if err := client.Create(ctx, key, trigger, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Trigger: create trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Remove trigger by id from storage
func (s *TriggerStorage) Remove(ctx context.Context, trigger *types.Trigger) error {

	log.V(logLevel).Debugf("Storage: Trigger: remove trigger by id %#v", trigger.Meta.Name)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Trigger: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(triggerStorage, trigger.Meta.Name)
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("Storage: Trigger: remove trigger err: %s", err.Error())
		return err
	}
	return nil
}

func newTriggerStorage() *TriggerStorage {
	s := new(TriggerStorage)
	return s
}
