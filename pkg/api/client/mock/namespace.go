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

	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
)

type NamespaceClient struct {
	data map[string]*v1.Namespace
	interfaces.Namespace
}

func (s *NamespaceClient) List(ctx context.Context) (*v1.NamespaceList, error) {
	list := make(v1.NamespaceList, 0)
	for _, ns := range s.data {
		list = append(list, ns)
	}
	return &list, nil
}

func newNamespaceClient() *NamespaceClient {
	s := new(NamespaceClient)
	s.data = make(map[string]*v1.Namespace)
	return s
}
