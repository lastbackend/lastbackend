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

func (sc *SecretClient) Create(ctx context.Context, opts *rv1.SecretCreateOptions) (*vv1.Secret, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Secret
	var e *errors.Http

	err = sc.client.Post(fmt.Sprintf("/namespace/%s/secret", sc.namespace)).
		AddHeader("Content-Entity", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (sc *SecretClient) List(ctx context.Context) (*vv1.SecretList, error) {

	var s *vv1.SecretList
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/secret", sc.namespace)).
		AddHeader("Content-Entity", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.SecretList, 0)
		s = &list
	}

	return s, nil
}

func (sc *SecretClient) Update(ctx context.Context, opts *rv1.SecretUpdateOptions) (*vv1.Secret, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Secret
	var e *errors.Http

	err = sc.client.Put(fmt.Sprintf("/namespace/%s/secret/%s", sc.namespace, sc.name)).
		AddHeader("Content-Entity", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (sc *SecretClient) Remove(ctx context.Context, opts *rv1.SecretRemoveOptions) error {

	req := sc.client.Delete(fmt.Sprintf("/namespace/%s/secret/%s", sc.namespace, sc.name)).
		AddHeader("Content-Entity", "application/json")

	if opts != nil {
		if opts.Force {
			req.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	var e *errors.Http

	if err := req.JSON(nil, &e); err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func newSecretClient(client http.Interface, namespace, name string) *SecretClient {
	return &SecretClient{client: client, namespace: namespace, name: name}
}
