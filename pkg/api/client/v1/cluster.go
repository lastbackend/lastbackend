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
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type ClusterClient struct {
	interfaces.Cluster
	client http.Interface
}

func (s *ClusterClient) Node(hostname ...string) *NodeClient {
	hst := ""
	if len(hostname) > 0 {
		hst = hostname[0]
	}
	return newNodeClient(s.client, hst)
}

func (s *ClusterClient) Ingress(name ...string) *IngressClient {
	hst := ""
	if len(name) > 0 {
		hst = name[0]
	}
	return newIngressClient(s.client, hst)
}

func (s *ClusterClient) Get(ctx context.Context) (*vv1.Cluster, error) {

	res := s.client.Get("/cluster").
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

	var cl *vv1.Cluster

	if err := json.Unmarshal(buf, &cl); err != nil {
		return nil, err
	}

	return cl, nil
}

func (s *ClusterClient) Update(ctx context.Context, opts *rv1.ClusterUpdateOptions) (*vv1.Cluster, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	res := s.client.Put("/cluster").
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

	var cl *vv1.Cluster

	if err := json.Unmarshal(buf, &cl); err != nil {
		return nil, err
	}

	return cl, nil
}

func newClusterClient(req http.Interface) *ClusterClient {
	return &ClusterClient{client: req}
}
