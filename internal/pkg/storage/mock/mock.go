//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"fmt"
	"github.com/lastbackend/lastbackend/internal/pkg/storage/types"
	"strings"
	"sync"

	"reflect"

	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/converter"
)

type Storage struct {
	root     string
	lock     sync.RWMutex
	store    map[string]map[string][]byte
	watchers map[chan *types.WatcherEvent]string
}

func (s *Storage) Info(ctx context.Context, collection string, name string) (*types.System, error) {
	return new(types.System), nil
}

func (s *Storage) Get(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {
	s.check(collection)

	s.lock.Lock()
	defer s.lock.Unlock()

	collection = fmt.Sprintf("%s/%s", s.root, collection)

	if _, ok := s.store[collection][name]; !ok {
		return errors.New(types.ErrEntityNotFound)
	}

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	if err := json.Unmarshal(s.store[collection][name], obj); err != nil {
		return err
	}

	return nil
}

func (s *Storage) List(ctx context.Context, collection string, q string, obj interface{}, opts *types.Opts) error {
	s.check(collection)

	s.lock.Lock()
	defer s.lock.Unlock()

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	v, err := converter.EnforcePtr(obj)
	if err != nil {
		return errors.New(types.ErrStructOutIsNotPointer)
	}

	collection = fmt.Sprintf("%s/%s", s.root, collection)
	buffer := []byte("[")
	current := 0

	for k, item := range s.store[collection] {
		if strings.HasPrefix(k, q) {

			if current > 0 {
				buffer = append(buffer, []byte(",")...)
			}

			buffer = append(buffer, item...)
			current++

		}
	}

	buffer = append(buffer, []byte("]")...)

	f := v.FieldByName("Items")
	if f.Kind() != reflect.Slice {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	if !f.IsValid() {
		return nil
	}

	if !f.CanSet() {
		return nil
	}

	items := reflect.New(f.Type()).Interface().(interface{})
	if err := json.Unmarshal(buffer, items); err != nil {
		return err
	}

	f.Set(reflect.ValueOf(items).Elem())
	return nil
}

func (s *Storage) Map(ctx context.Context, collection string, q string, obj interface{}, opts *types.Opts) error {
	s.check(collection)

	s.lock.Lock()
	defer s.lock.Unlock()

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	v, err := converter.EnforcePtr(obj)
	if err != nil {
		return errors.New(types.ErrStructOutIsNotPointer)
	}

	collection = fmt.Sprintf("%s/%s", s.root, collection)

	buffer := []byte("{")
	current := 0

	for k, item := range s.store[collection] {
		if strings.HasPrefix(k, q) {

			ks := strings.Split(k, "/")

			if current > 0 {
				buffer = append(buffer, []byte(",")...)
			}

			buffer = append(buffer, []byte("\"")...)
			buffer = append(buffer, []byte(ks[len(ks)-1])...)
			buffer = append(buffer, []byte("\":")...)
			buffer = append(buffer, item...)
			current++

		}
	}

	buffer = append(buffer, []byte("}")...)

	f := v.FieldByName("Items")
	if f.Kind() != reflect.Map {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	if !f.IsValid() {
		return nil
	}

	if !f.CanSet() {
		return nil
	}

	items := reflect.New(f.Type()).Interface().(interface{})
	if err := json.Unmarshal(buffer, items); err != nil {
		return err
	}

	f.Set(reflect.ValueOf(items).Elem())

	return nil
}

func (s *Storage) Put(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {
	s.check(collection)

	s.lock.Lock()
	defer s.lock.Unlock()

	collection = fmt.Sprintf("%s/%s", s.root, collection)

	if _, ok := s.store[collection][name]; ok {

		if opts == nil {
			return errors.New(types.ErrEntityExists)
		}

		if !opts.Force {
			return errors.New(types.ErrEntityExists)
		}
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	s.store[collection][name] = b

	s.dispatch(collection, name, types.STORAGECREATEEVENT, b)
	return nil
}

func (s *Storage) Set(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {

	s.check(collection)

	s.lock.Lock()
	defer s.lock.Unlock()

	collection = fmt.Sprintf("%s/%s", s.root, collection)

	if _, ok := s.store[collection][name]; !ok {
		if opts != nil && !opts.Force {
			return errors.New(types.ErrEntityNotFound)
		}
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	s.store[collection][name] = b

	s.dispatch(collection, name, types.STORAGEUPDATEEVENT, b)
	return nil
}

func (s *Storage) Del(ctx context.Context, collection string, name string) error {
	s.check(collection)

	s.lock.Lock()
	defer s.lock.Unlock()

	collection = fmt.Sprintf("%s/%s", s.root, collection)
	if name == "" {
		s.store[collection] = make(map[string][]byte)
		return nil
	}

	bt := s.store[collection][name]
	delete(s.store[collection], name)

	s.dispatch(collection, name, types.STORAGEDELETEEVENT, bt)

	return nil
}

func (s *Storage) Watch(ctx context.Context, collection string, event chan *types.WatcherEvent, opts *types.Opts) error {

	s.check(collection)
	s.lock.Lock()
	s.watchers[event] = collection
	s.lock.Unlock()

	defer delete(s.watchers, event)
	<-ctx.Done()
	return nil
}

func (s Storage) Filter() types.Filter {
	return new(Filter)
}

func (s Storage) Key() types.Key {
	return new(Key)
}

func (s Storage) Collection() types.Collection {
	return new(Collection)
}

func (s *Storage) dispatch(collection, name, action string, b []byte) {

	for w, c := range s.watchers {

		if c == collection || c == s.Collection().Root() {

			e := new(types.WatcherEvent)
			e.Action = action
			e.SelfLink = name
			e.Storage.Key = fmt.Sprintf("%s/%s", strings.TrimPrefix(collection, s.root), name)
			e.Storage.Revision = 0
			e.Data = b

			match := strings.Split(name, ":")

			if len(match) > 0 {
				e.Name = match[len(match)-1]
			} else {
				e.Name = name
			}

			w <- e
		}
	}
}

func (s *Storage) check(kind string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	collection := fmt.Sprintf("%s/%s", s.root, kind)

	if _, ok := s.store[collection]; !ok {
		s.store[collection] = make(map[string][]byte)
	}
}

func New() (*Storage, error) {
	db := new(Storage)

	db.root = "lastbackend"
	db.store = make(map[string]map[string][]byte)
	db.watchers = make(map[chan *types.WatcherEvent]string, 0)

	return db, nil
}
