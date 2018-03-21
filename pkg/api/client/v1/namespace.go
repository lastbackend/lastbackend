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

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type NamespaceClient struct {
	interfaces.Namespace
	client http.Interface
	name   string
}

func (s *NamespaceClient) Service(name ...string) *ServiceClient {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return newServiceClient(s.client, s.name, n)
}

func (s *NamespaceClient) Secret(name ...string) *SecretClient {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return newSecretClient(s.client, s.name, n)
}

func (s *NamespaceClient) Route(name ...string) *RouteClient {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return newRouteClient(s.client, s.name, n)
}

func (s *NamespaceClient) Volume(name ...string) *VolumeClient {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return newVolumeClient(s.client, s.name, n)
}

func (s *NamespaceClient) List(ctx context.Context) (*vv1.NamespaceList, error) {

	req := s.client.Get(fmt.Sprintf("/namespace")).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	var nl *vv1.NamespaceList

	if err := json.Unmarshal(buf, &nl); err != nil {
		return nil, err
	}

	return nl, nil
}

func (s *NamespaceClient) Create(ctx context.Context, opts *rv1.NamespaceCreateOptions) (*vv1.Namespace, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Post("/namespace").
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	if err := req.Error(); err != nil {
		return nil, err
	}

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ns = new(vv1.Namespace)

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NamespaceClient) Get(ctx context.Context) (*vv1.Namespace, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s", s.name)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ns *vv1.Namespace

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NamespaceClient) Update(ctx context.Context, opts *rv1.NamespaceUpdateOptions) (*vv1.Namespace, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Put(fmt.Sprintf("/namespace/%s", s.name)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ns *vv1.Namespace

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NamespaceClient) Remove(ctx context.Context, opts *rv1.NamespaceRemoveOptions) error {

	req := s.client.Delete(fmt.Sprintf("/namespace/%s", s.name)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func newNamespaceClient(client http.Interface, name string) *NamespaceClient {
	return &NamespaceClient{client: client, name: name}
}
