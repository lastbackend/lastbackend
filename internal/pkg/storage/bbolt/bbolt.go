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

package bbolt

import (
	"errors"
	"reflect"

	"github.com/lastbackend/lastbackend/internal/pkg/storage/types"
	"github.com/lastbackend/lastbackend/internal/util/converter"
	"github.com/lastbackend/lastbackend/internal/util/serializer"
	"github.com/lastbackend/lastbackend/internal/util/serializer/json"
	bolt "go.etcd.io/bbolt"
)

type Storage struct {
	db    *bolt.DB
	codec serializer.Codec
}

type Options struct {
	// Path of the DB file.
	// Optional ("lb.db" by default).
	Path string
	// Encoding format.
	// Optional (encoding.JSON by default).
	Codec serializer.Codec
}

// DefaultOptions is an Options object with default values.
// BucketName: "default", Path: "lb.db", Codec: encoding.JSON
var DefaultOptions = Options{
	Path:  "lb.db",
	Codec: serializer.NewSerializer(json.Encoder{}, json.Decoder{}),
}

func New(options Options) (*Storage, error) {

	if options.Path == "" {
		options.Path = DefaultOptions.Path
	}
	if options.Codec == nil {
		options.Codec = DefaultOptions.Codec
	}

	db, err := bolt.Open(options.Path, 0600, nil)
	if err != nil {
		return nil, err
	}

	s := new(Storage)
	s.db = db
	s.codec = options.Codec

	return s, nil
}

func (s Storage) List(collection string, listOutPtr interface{}) error {
	if err := checkCollection(collection); err != nil {
		return err
	}

	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return errors.New(types.ErrCollectionNotExists)
		}

		items := make(map[string][]byte)
		err := b.ForEach(func(k, v []byte) error {
			items[string(k)] = v
			return nil
		})
		if err != nil {
			return err
		}

		return decodeList(s.codec, items, listOutPtr)
	})

}

func (s Storage) Get(collection, key string, outPtr interface{}) error {
	if err := checkCollection(collection); err != nil {
		return err
	}
	if err := checkKey(key); err != nil {
		return err
	}

	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return errors.New(types.ErrCollectionNotExists)
		}
		data := b.Get([]byte(key))
		if data != nil {
			return nil
		}
		return decode(s.codec, data, outPtr)
	})

}

func (s Storage) Set(collection, key string, obj interface{}) error {
	if err := checkCollection(collection); err != nil {
		return err
	}
	if err := checkKey(key); err != nil {
		return err
	}
	if err := checkValue(obj); err != nil {
		return err
	}

	buf, err := serializer.Encode(s.codec, obj)
	if err != nil {
		return err
	}

	if err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), buf)
	}); err != nil {
		return err
	}

	return nil
}

func (s Storage) Put(collection, key string, obj interface{}) error {
	if err := checkCollection(collection); err != nil {
		return err
	}
	if err := checkKey(key); err != nil {
		return err
	}
	if err := checkValue(obj); err != nil {
		return err
	}

	buf, err := serializer.Encode(s.codec, obj)
	if err != nil {
		return err
	}

	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return errors.New(types.ErrCollectionNotExists)
		}
		return b.Put([]byte(key), buf)
	}); err != nil {
		return err
	}

	return nil
}

func (s Storage) Del(collection, key string) error {
	if err := checkCollection(key); err != nil {
		return err
	}
	if err := checkKey(key); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(collection))
		if b == nil {
			return errors.New(types.ErrCollectionNotExists)
		}
		return b.Delete([]byte(key))
	})
}

func (s Storage) Close() error {
	return s.db.Close()
}

func checkValue(v interface{}) error {
	if v == nil {
		return errors.New(types.ErrValueIsNil)
	}
	return nil
}

func checkCollection(k string) error {
	if k == "" {
		return errors.New(types.ErrCollectionIsEmpty)
	}
	return nil
}

func checkKey(k string) error {
	if k == "" {
		return errors.New(types.ErrKeyIsEmpty)
	}
	return nil
}

func decode(s serializer.Codec, value []byte, out interface{}) error {
	if _, err := converter.EnforcePtr(out); err != nil {
		return errors.New(types.ErrStructOutIsNil)
	}
	return serializer.Decode(s, value, out)
}

func decodeList(codec serializer.Codec, items map[string][]byte, listOut interface{}) error {
	v, err := converter.EnforcePtr(listOut)
	if err != nil {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	if v.Kind() != reflect.Slice {
		return errors.New(types.ErrStructOutIsInvalid)
	}

	if !v.IsValid() {
		return nil
	}

	if !v.CanSet() {
		return nil
	}

	for _, item := range items {

		var obj = reflect.New(v.Type().Elem()).Interface().(interface{})
		err := serializer.Decode(codec, item, obj)
		if err != nil {
			return err
		}

		v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
	}

	return nil
}
