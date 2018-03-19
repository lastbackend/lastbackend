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

	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type TriggerClient struct {
	interfaces.Trigger
}

func (s *TriggerClient) Create(ctx context.Context, namespace, service string, opts rv1.TriggerCreateOptions) (*vv1.Trigger, error) {
	return nil, nil
}

func (s *TriggerClient) List(ctx context.Context, namespace, service string) (*vv1.TriggerList, error) {
	return nil, nil
}

func (s *TriggerClient) Get(ctx context.Context, namespace, service, name string) (*vv1.Trigger, error) {
	return nil, nil
}

func (s *TriggerClient) Update(ctx context.Context, namespace, service, name string, opts rv1.TriggerUpdateOptions) (*vv1.Trigger, error) {
	return nil, nil
}

func (s *TriggerClient) Remove(ctx context.Context, namespace, service, name string, opts rv1.TriggerRemoveOptions) error {
	return nil
}

func newTriggerClient() *TriggerClient {
	s := new(TriggerClient)
	return s
}