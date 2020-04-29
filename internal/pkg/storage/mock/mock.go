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

type Storage struct {
}

func (s *Storage) List(collection string, listOutPtr interface{}) error {
	return nil
}

func (s *Storage) Get(collection, key string, outPtr interface{}) error {
	return nil
}

func (s *Storage) Set(collection, key string, obj interface{}) error {
	return nil
}

func (s *Storage) Put(collection, key string, obj interface{}) error {
	return nil
}

func (s *Storage) Del(collection, key string) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func New() (*Storage, error) {
	db := new(Storage)
	return db, nil
}
