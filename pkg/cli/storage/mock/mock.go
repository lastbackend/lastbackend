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

package mock

import (
	"github.com/lastbackend/lastbackend/pkg/cli/storage"
	"github.com/lastbackend/lastbackend/pkg/cli/storage/db"
)

type Storage struct {
	*AppStorage
}

func (s *Storage) App() storage.IApp {
	if s == nil {
		return nil
	}
	return s.AppStorage
}

func Get() (*Storage, error) {

	client, err := db.Init()
	if err != nil {
		return nil, err
	}

	store := new(Storage)
	store.AppStorage = newAppStorage(client)

	return store, nil
}
