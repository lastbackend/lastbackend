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
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/cli/storage/db"
	"github.com/lastbackend/lastbackend/pkg/cli/storage"
)

const namespaceStorage = "mocknamespace"

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	storage.INamespace
	client *db.DB
}

// Insert namespace
func (s *NamespaceStorage) Save(namespace *n.Namespace) error {
	return s.client.Set(namespaceStorage, namespace)
}

// Get namespace
func (s *NamespaceStorage) Load() (*n.Namespace, error) {
	var ns = new(n.Namespace)
	err := s.client.Get(namespaceStorage, ns)
	return ns, err
}

// Remove namespace
func (s *NamespaceStorage) Remove() error {
	return s.client.Set(namespaceStorage, nil)
}

func newNamespaceStorage(client *db.DB) *NamespaceStorage {
	s := new(NamespaceStorage)
	s.client = client
	return s
}

