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
		"strings"

	"reflect"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type Storage struct {
	store map[string]map[string][]byte
}

func (s *Storage) Info(ctx context.Context, collection string, name string) (*types.Runtime, error) {
	return new(types.Runtime), nil
}

func (s *Storage) Get(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {
	s.check(collection)

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

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	v, err := converter.EnforcePtr(obj)
	if err != nil {
		return errors.New(types.ErrStructOutIsNotPointer)
	}

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

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	v, err := converter.EnforcePtr(obj)
	if err != nil {
		return errors.New(types.ErrStructOutIsNotPointer)
	}

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
	return nil
}

func (s *Storage) Set(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {
	s.check(collection)

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

	return nil
}

func (s *Storage) Del(ctx context.Context, collection string, name string)  error {
	s.check(collection)
	if name == "" {
		s.store[collection] = make(map[string][]byte)
		return nil
	}
	delete(s.store[collection], name)
	return nil
}

func (s *Storage) Watch(ctx context.Context, collection string, event chan *types.WatcherEvent, opts *types.Opts)  error {
	s.check(collection)
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



func (s *Storage) check(kind string) {
	if _, ok := s.store[kind]; !ok {
		s.store[kind] = make(map[string][]byte)
	}
}

func New() (*Storage, error) {
	db := new(Storage)
	db.store = make(map[string]map[string][]byte)
	return db, nil
}
