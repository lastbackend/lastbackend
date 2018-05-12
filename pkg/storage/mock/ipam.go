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
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"context"
)

type IPAMStorage struct {
	storage.IPAM
	data []string
}

func (s *IPAMStorage) Get(ctx context.Context) ([]string, error) {
	return s.data, nil
}

func (s *IPAMStorage) Set(ctx context.Context, ips []string) error {
	s.data = ips
	return nil
}

func newIPAMStorage() *IPAMStorage {
	s := new(IPAMStorage)
	s.data = make([]string, 0)
	return s
}

