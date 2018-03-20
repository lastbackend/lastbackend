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

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type NamespaceClient struct {
	interfaces.Namespace
	client http.Interface
}

func (s *NamespaceClient) Service(namespace string) *ServiceClient {
	return newServiceClient(s.client, namespace)
}

func (s *NamespaceClient) List(ctx context.Context) (*vv1.NamespaceList, error) {

	var (
		nl *vv1.NamespaceList
	)

	result := s.client.Get(fmt.Sprintf("/namespace")).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := result.Raw()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &nl); err != nil {
		return nil, err
	}

	return nl, nil
}

func (s *NamespaceClient) Create(ctx context.Context, opts rv1.NamespaceCreateOptions) (*vv1.Namespace, error) {

	var (
		ns *vv1.Namespace
	)

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	result := s.client.Post("/namespace").
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := result.Raw()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NamespaceClient) Get(ctx context.Context, name string) (*vv1.Namespace, error) {
	var (
		ns *vv1.Namespace
	)

	result := s.client.Get(fmt.Sprintf("/namespace/%s", name)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := result.Raw()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NamespaceClient) Update(ctx context.Context, name string, opts rv1.NamespaceUpdateOptions) (*vv1.Namespace, error) {
	var (
		ns *vv1.Namespace
	)

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	result := s.client.Put(fmt.Sprintf("/namespace/%s", name)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := result.Raw()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NamespaceClient) Remove(ctx context.Context, name string, opts rv1.NamespaceRemoveOptions) error {

	s.client.Delete(fmt.Sprintf("/namespace/%s", name)).
		AddHeader("Content-Type", "application/json").
		Do()

	return nil
}

func newNamespaceClient(client http.Interface) *NamespaceClient {
	return &NamespaceClient{client: client}
}
