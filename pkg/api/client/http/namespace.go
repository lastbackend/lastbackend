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

package http

import (
	"context"
	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
)

type NamespaceClient struct {
	interfaces.Namespace
}

func (s *NamespaceClient) List(ctx context.Context) (*v1.NamespaceList, error) {

	var (
		r  = NewRequest(ctx)
		nl *v1.NamespaceList
	)

	body, err := r.Get("namespace")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, nl); err != nil {
		return nil, err
	}

	return nl, nil
}

func newNamespaceClient() *NamespaceClient {
	s := new(NamespaceClient)
	return s
}
