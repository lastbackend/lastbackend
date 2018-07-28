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

package cluster

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

func NodeLease(req types.NodeLease, nodes map[string]*types.Node) (*types.Node, error) {

	for _, n := range nodes {
		if n.Status.Capacity.Memory < *req.Request.Memory {

			n.Status.Allocated.Pods++
			n.Status.Allocated.Memory += *req.Request.Memory
			nm := distribution.NewNodeModel(context.Background(), envs.Get().GetStorage())
			nm.SetStatus(n, n.Status)

			return n, nil
		}
	}

	return nil, nil
}

func NodeRelease() {

}
