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

package storage

import (
	n "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
	"github.com/lastbackend/lastbackend/pkg/cli/storage/db"
)

const appStorage = "app"

// App Service type for interface in interfaces folder
type AppStorage struct {
	IApp
	client *db.DB
}

// Insert app
func (s *AppStorage) Save(app *n.App) error {
	return s.client.Set(appStorage, app)
}

// Get app
func (s *AppStorage) Load() (*n.App, error) {
	var ns = new(n.App)
	err := s.client.Get(appStorage, ns)
	return ns, err
}

// Remove app
func (s *AppStorage) Remove() error {
	return s.client.Set(appStorage, nil)
}

func newAppStorage(client *db.DB) *AppStorage {
	s := new(AppStorage)
	s.client = client
	return s
}
