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
	"context"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const endpointStorage = "endpoints"

// EndpointStorage type for interface in interfaces folder
type EndpointStorage struct {
	storage.Endpoint
}

// Get endpoints by domain name
func (s *EndpointStorage) Get(ctx context.Context, name string) ([]string, error) {
	return make([]string, 0), nil
}

// Upsert endpoint model
func (s *EndpointStorage) Upsert(ctx context.Context, name string, ips []string) error {
	return nil
}

// Remove endpoint model
func (s *EndpointStorage) Remove(ctx context.Context, name string) error {
	return nil
}

// Watch endpoint model
func (s *EndpointStorage) Watch(ctx context.Context, endpoint chan string) error {
	return nil
}

func newEndpointStorage() *EndpointStorage {
	s := new(EndpointStorage)
	return s
}
