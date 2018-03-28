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

func (s *PodClient) List(ctx context.Context) (*vv1.PodList, error) {

	res := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s/pod", s.namespace, s.service, s.deployment)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	var pl *vv1.PodList

	if err := json.Unmarshal(buf, &pl); err != nil {
		return nil, err
	}

	return pl, nil
}

func (s *PodClient) Get(ctx context.Context, name string) (*vv1.Pod, error) {

	res := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s/pod/%s", s.namespace, s.service, s.deployment, name)).
		AddHeader("Content-Type", "application/json").
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

	var ps *vv1.Pod

	if err := json.Unmarshal(buf, &ps); err != nil {
		return nil, err
	}

	return ps, nil
}

func newPodClient(client http.Interface, namespace, service, deployment string) *PodClient {
	return &PodClient{client: client, namespace: namespace, service: service, deployment: deployment}
}
