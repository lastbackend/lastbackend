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
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type PodClient struct {
	interfaces.Pod
	req        http.Interface
	namespace  string
	service    string
	deployment string
}

func (s *PodClient) List(ctx context.Context) (*vv1.PodList, error) {
	return nil, nil
}

func (s *PodClient) Get(ctx context.Context, na string) (*vv1.Pod, error) {
	return nil, nil
}

func newPodClient(req http.Interface, namespace, service, deployment string) *PodClient {
	return &PodClient{req: req, namespace: namespace, service: service, deployment: deployment}
}
