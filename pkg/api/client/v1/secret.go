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

	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type SecretClient struct {
	interfaces.Secret
	client    http.Interface
	namespace string
	name      string
}

func (s *SecretClient) Create(ctx context.Context, opts *rv1.SecretCreateOptions) (*vv1.Secret, error) {
	return nil, nil
}

func (s *SecretClient) List(ctx context.Context) (*vv1.SecretList, error) {
	return nil, nil
}

func (s *SecretClient) Get(ctx context.Context) (*vv1.Secret, error) {
	return nil, nil
}

func (s *SecretClient) Update(ctx context.Context, opts *rv1.SecretUpdateOptions) (*vv1.Secret, error) {
	return nil, nil
}

func (s *SecretClient) Remove(ctx context.Context, opts *rv1.SecretRemoveOptions) error {
	return nil
}

func newSecretClient(client http.Interface, namespace, name string) *SecretClient {
	s := new(SecretClient)
	s.client = client
	s.namespace = namespace
	s.name = name
	return s
}
