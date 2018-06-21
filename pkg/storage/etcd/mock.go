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

package etcd

import (
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/mock"
)

type StorageV3Mock struct {
}

func NewV3Mock() (*StorageV3, error) {

	log.Debug("%s:> define v3 mock storage", logPrefix)

	var (
		err error
		s   = new(StorageV3)
	)

	s.client = new(clientV3)

	if s.client.store, s.client.dfunc, err = mock.GetMockClient(); err != nil {
		log.Errorf("%s:> store mock initialize err: %v", logPrefix, err)
		return nil, err
	}

	return s, nil
}
