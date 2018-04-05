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
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type PodClient struct {
	interfaces.Pod
	client     http.Interface
	namespace  string
	service    string
	deployment string
}

func (pc *PodClient) List(ctx context.Context) (*vv1.PodList, error) {

	var s *vv1.PodList
	var e *errors.Http

	err := pc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s/pod", pc.namespace, pc.service, pc.deployment)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.PodList, 0)
		s = &list
	}

	return s, nil
}

func (pc *PodClient) Get(ctx context.Context, name string) (*vv1.Pod, error) {

	var s *vv1.Pod
	var e *errors.Http

	err := pc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s/pod/%s", pc.namespace, pc.service, pc.deployment, name)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func newPodClient(client http.Interface, namespace, service, deployment string) *PodClient {
	return &PodClient{client: client, namespace: namespace, service: service, deployment: deployment}
}
