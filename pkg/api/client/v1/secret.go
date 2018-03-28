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
	"strconv"
)

type SecretClient struct {
	interfaces.Secret
	client    http.Interface
	namespace string
	name      string
}

func (s *SecretClient) Create(ctx context.Context, opts *rv1.SecretCreateOptions) (*vv1.Secret, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	res := s.client.Post(fmt.Sprintf("/namespace/%s/secret", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	if err := res.Error(); err != nil {
		return nil, err
	}

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ss = new(vv1.Secret)

	if err := json.Unmarshal(buf, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (s *SecretClient) List(ctx context.Context) (*vv1.SecretList, error) {

	res := s.client.Get(fmt.Sprintf("/namespace/%s/secret", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	var sl *vv1.SecretList

	if len(buf) == 0 {
		list := make(vv1.SecretList, 0)
		return &list, nil
	}

	if err := json.Unmarshal(buf, &sl); err != nil {
		return nil, err
	}

	return sl, nil
}

func (s *SecretClient) Update(ctx context.Context, opts *rv1.SecretUpdateOptions) (*vv1.Secret, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	res := s.client.Put(fmt.Sprintf("/namespace/%s/secret/%s", s.namespace, s.name)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ss *vv1.Secret

	if err := json.Unmarshal(buf, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (s *SecretClient) Remove(ctx context.Context, opts *rv1.SecretRemoveOptions) error {

	res := s.client.Delete(fmt.Sprintf("/namespace/%s/secret/%s", s.namespace, s.name)).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Force {
			res.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	req := res.Do()

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

func newSecretClient(client http.Interface, namespace, name string) *SecretClient {
	s := new(SecretClient)
	s.client = client
	s.namespace = namespace
	s.name = name
	return s
}
