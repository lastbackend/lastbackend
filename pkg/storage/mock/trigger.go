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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"fmt"
	"strings"
)

// Service Trigger type for interface in interfaces folder
type TriggerStorage struct {
	storage.Trigger
	data map[string]*types.Trigger
}

// Get trigger by id
func (s *TriggerStorage) Get(ctx context.Context, namespace, service, name string) (*types.Trigger, error) {
	if ns, ok := s.data[s.keyCreate(namespace, service, name)]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

func (s *TriggerStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Trigger, error) {
	list := make(map[string]*types.Trigger, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

func (s *TriggerStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Trigger, error) {
	list := make(map[string]*types.Trigger, 0)

	prefix := fmt.Sprintf("%s:%s:", namespace, service)

	for _, d := range s.data {
		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Insert new trigger into storage
func (s *TriggerStorage) Insert(ctx context.Context, trigger *types.Trigger) error {

	if err := s.checkTriggerArgument(trigger); err != nil {
		return err
	}

	s.data[s.keyGet(trigger)] = trigger

	return nil
}

// Update trigger info
func (s *TriggerStorage) Update(ctx context.Context, trigger *types.Trigger) error {

	if err := s.checkTriggerExists(trigger); err != nil {
		return err
	}

	s.data[s.keyGet(trigger)] = trigger

	return nil
}

// Remove trigger by from storage
func (s *TriggerStorage) Remove(ctx context.Context, trigger *types.Trigger) error {
	if err := s.checkTriggerExists(trigger); err != nil {
		return err
	}

	delete(s.data, s.keyGet(trigger))
	return nil
}

// Watch trigger changes
func (s *TriggerStorage) Watch(ctx context.Context, trigger chan *types.Trigger) error {
	return nil
}

// Watch trigger spec changes
func (s *TriggerStorage) WatchSpec(ctx context.Context, trigger chan *types.Trigger) error {
	return nil
}

// Watch trigger state changes
func (s *TriggerStorage) WatchState(ctx context.Context, trigger chan *types.Trigger) error {
	return nil
}

// Clear trigger storage
func (s *TriggerStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Trigger)
	return nil
}

// keyCreate util function
func (s *TriggerStorage) keyCreate (namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
}

// keyGet util function
func (s *TriggerStorage) keyGet(t *types.Trigger) string {
	return t.SelfLink()
}

// newTriggerStorage return new trigger storage
func newTriggerStorage() *TriggerStorage {
	s := new(TriggerStorage)
	s.data = make(map[string]*types.Trigger)
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

// checkTriggerArgument - check if deployment exists in store
func (s *TriggerStorage) checkTriggerExists(trigger *types.Trigger) error {

	if err := s.checkTriggerArgument(trigger); err != nil {
		return err
	}

	if _, ok := s.data[s.keyGet(trigger)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}