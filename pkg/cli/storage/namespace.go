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
	"github.com/lastbackend/lastbackend/pkg/cli/storage/db"
	v "github.com/lastbackend/lastbackend/pkg/cli/view"
)

const nsStorage = "namespace"

// Namespace type for interface in interfaces folder
type NamespaceStorage struct {
	INamespace
	client *db.DB
}

// Insert namespace
func (s *NamespaceStorage) Save(data *v.Namespace) error {
	return s.client.Set(nsStorage, data)
}

// Get namespace
func (s *NamespaceStorage) Load() (*v.Namespace, error) {
	var data = new(v.Namespace)
	if err := s.client.Get(nsStorage, data); err != nil {
		return nil, err
	}
	return data, nil
}

// Destroy namespace
func (s *NamespaceStorage) Remove() error {
	return s.client.Set(nsStorage, nil)
}

func newNamespaceStorage(client *db.DB) *NamespaceStorage {
	s := new(NamespaceStorage)
	s.client = client
	return s
}
