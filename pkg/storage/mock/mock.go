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
	"errors"

	"strings"

	"reflect"

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/storage/types"
	"github.com/lastbackend/lastbackend/pkg/util/converter"
)

type Storage struct {
	store map[types.Kind]map[string][]byte
}

func (s *Storage) Get(ctx context.Context, kind types.Kind, name string, obj interface{}) error {
	s.check(kind)

	if _, ok := s.store[kind][name]; !ok {
		return errors.New(types.ErrEntityNotFound)
	}

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	if err := json.Unmarshal(s.store[kind][name], obj); err != nil {
		return err
	}

	return nil
}

func (s *Storage) List(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	s.check(kind)

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	v, err := converter.EnforcePtr(obj)
	if err != nil {
		return errors.New(types.ErrStructOutIsNotPointer)
	}

	if v.Kind() != reflect.Slice {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	buffer := []byte("[")
	current := 0
	for k, item := range s.store[kind] {
		if strings.HasPrefix(k, q) {

			if current > 0 {
				buffer = append(buffer, []byte(",")...)
			}

			buffer = append(buffer, item...)
			current++

		}
	}

	buffer = append(buffer, []byte("]")...)

	if err := json.Unmarshal(buffer, obj); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Map(ctx context.Context, kind types.Kind, q string, obj interface{}) error {
	s.check(kind)

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	v, err := converter.EnforcePtr(obj)
	if err != nil {
		return errors.New(types.ErrStructOutIsNotPointer)
	}

	if v.Kind() != reflect.Map {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	buffer := []byte("{")
	current := 0
	for k, item := range s.store[kind] {
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

	if err := json.Unmarshal(buffer, obj); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Put(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	s.check(kind)

	if _, ok := s.store[kind][name]; ok {

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

	s.store[kind][name] = b
	return nil
}

func (s *Storage) Set(ctx context.Context, kind types.Kind, name string, obj interface{}, opts *types.Opts) error {
	s.check(kind)

	if _, ok := s.store[kind][name]; !ok {
		if opts != nil && !opts.Force {
			return errors.New(types.ErrEntityNotFound)
		}
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	s.store[kind][name] = b

	return nil
}

func (s *Storage) Del(ctx context.Context, kind types.Kind, name string) error {
	s.check(kind)
	if name == "" {
		s.store[kind] = make(map[string][]byte)
		return nil
	}
	delete(s.store[kind], name)
	return nil
}

func (s *Storage) Watch(ctx context.Context, kind types.Kind, event chan *types.WatcherEvent) error {
	s.check(kind)
	return nil
}

func (s Storage) Filter() types.Filter {
	return new(Filter)
}

func (s Storage) Key() types.Key {
	return new(Key)
}

func (s *Storage) check(kind types.Kind) {
	if _, ok := s.store[kind]; !ok {
		s.store[kind] = make(map[string][]byte)
	}
}

func New() (*Storage, error) {
	db := new(Storage)
	db.store = make(map[types.Kind]map[string][]byte)
	return db, nil
}
