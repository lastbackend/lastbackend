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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const triggerStorage = "triggers"

// Service Trigger type for interface in interfaces folder
type TriggerStorage struct {
	storage.Trigger
}

// Get triggers by id
func (s *TriggerStorage) Get(ctx context.Context, namespace, service, name string) (*types.Trigger, error) {

	log.V(logLevel).Debugf("storage:etcd:trigger:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:trigger:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:trigger:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:trigger:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + triggerStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		trigger = new(types.Trigger)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyTrigger := keyCreate(triggerStorage, s.keyCreate(namespace, service, name))
	if err := client.Map(ctx, keyTrigger, filter, trigger); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> err: %s", name, err.Error())
		return nil, err
	}

	if trigger.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return trigger, nil
}

// Get trigger by namespace name
func (s *TriggerStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Trigger, error) {

	log.V(logLevel).Debugf("storage:etcd:trigger:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:trigger:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + triggerStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		triggers = make(map[string]*types.Trigger)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyTrigger := keyCreate(triggerStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, keyTrigger, filter, triggers); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> err: %s", namespace, err.Error())
		return nil, err
	}

	return triggers, nil
}

// Get trigger by service name
func (s *TriggerStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Trigger, error) {

	log.V(logLevel).Debugf("storage:etcd:trigger:> get list by namespace and service: %s:%s", namespace, service)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:trigger:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:trigger:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + triggerStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		triggers = make(map[string]*types.Trigger)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:>  get list by namespace and service err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyTrigger := keyCreate(triggerStorage, fmt.Sprintf("%s:%s:", namespace, service))
	if err := client.MapList(ctx, keyTrigger, filter, triggers); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> err: %s", namespace, err.Error())
		return nil, err
	}

	return triggers, nil
}

// Update trigger status
func (s *TriggerStorage) SetStatus(ctx context.Context, trigger *types.Trigger) error {

	log.V(logLevel).Debugf("storage:etcd:trigger:> update trigger status: %#v", trigger)

	if err := s.checkTriggerExists(ctx, trigger); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:>: update trigger err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(triggerStorage, s.keyGet(trigger), "status")
	if err := client.Upsert(ctx, key, trigger.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:>: update trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Update trigger status
func (s *TriggerStorage) SetSpec(ctx context.Context, trigger *types.Trigger) error {

	log.V(logLevel).Debugf("storage:etcd:trigger:> update trigger spec: %#v", trigger)

	if err := s.checkTriggerExists(ctx, trigger); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:>: update trigger err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(triggerStorage, s.keyGet(trigger), "spec")
	if err := client.Upsert(ctx, key, trigger.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:>: update trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new trigger into storage
func (s *TriggerStorage) Insert(ctx context.Context, trigger *types.Trigger) error {

	log.V(logLevel).Debugf("storage:etcd:trigger:> insert trigger: %#v", trigger)

	if err := s.checkTriggerArgument(trigger); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> insert trigger err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(triggerStorage, s.keyGet(trigger), "meta")
	if err := tx.Create(keyMeta, trigger.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> insert trigger err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(triggerStorage, s.keyGet(trigger), "status")
	if err := tx.Create(keyStatus, trigger.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> insert trigger err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(triggerStorage, s.keyGet(trigger), "spec")
	if err := tx.Create(keySpec, trigger.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> insert trigger err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> insert trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Update trigger info
func (s *TriggerStorage) Update(ctx context.Context, trigger *types.Trigger) error {

	if err := s.checkTriggerExists(ctx, trigger); err != nil {
		return err
	}

	trigger.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> update trigger err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(triggerStorage, s.keyGet(trigger), "meta")
	if err := client.Upsert(ctx, keyMeta, trigger.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> update trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Remove trigger by id from storage
func (s *TriggerStorage) Remove(ctx context.Context, trigger *types.Trigger) error {

	if err := s.checkTriggerExists(ctx, trigger); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(triggerStorage, s.keyGet(trigger))
	if err := client.DeleteDir(ctx, keyMeta); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> remove trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Watch trigger changes
func (s *TriggerStorage) Watch(ctx context.Context, trigger chan *types.Trigger) error {

	log.V(logLevel).Debug("storage:etcd:trigger:> watch trigger")

	const filter = `\b\/` + triggerStorage + `\/(.+):(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> watch trigger err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(triggerStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3]); err == nil {
			trigger <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> watch trigger err: %s", err.Error())
		return err
	}

	return nil
}

// Watch trigger spec changes
func (s *TriggerStorage) WatchSpec(ctx context.Context, trigger chan *types.Trigger) error {

	log.V(logLevel).Debug("storage:etcd:trigger:> watch trigger by spec")

	const filter = `\b\/` + triggerStorage + `\/(.+):(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> watch trigger by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(triggerStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3]); err == nil {
			trigger <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> watch trigger by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Watch trigger status changes
func (s *TriggerStorage) WatchStatus(ctx context.Context, trigger chan *types.Trigger) error {

	log.V(logLevel).Debug("storage:etcd:trigger:> watch trigger by spec")

	const filter = `\b\/` + triggerStorage + `\/(.+):(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> watch trigger by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(triggerStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2], keys[3]); err == nil {
			trigger <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> watch trigger by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Clear trigger storage
func (s *TriggerStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:trigger:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, triggerStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:trigger:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *TriggerStorage) keyCreate(namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

// keyGet util function
func (s *TriggerStorage) keyGet(t *types.Trigger) string {
	return t.SelfLink()
}

func newTriggerStorage() *TriggerStorage {
	s := new(TriggerStorage)
	return s
}

// checkTriggerArgument - check if argument is valid for manipulations
func (s *TriggerStorage) checkTriggerArgument(trigger *types.Trigger) error {

	if trigger == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if trigger.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkTriggerArgument - check if trigger exists in store
func (s *TriggerStorage) checkTriggerExists(ctx context.Context, trigger *types.Trigger) error {

	if err := s.checkTriggerArgument(trigger); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:trigger:> check trigger exists")

	if _, err := s.Get(ctx, trigger.Meta.Namespace, trigger.Meta.Service, trigger.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:trigger:> check trigger exists err: %s", err.Error())
		return err
	}

	return nil
}
