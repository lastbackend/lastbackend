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
	"github.com/lastbackend/lastbackend/pkg/log"
)

func clusterStatusState(cs *ClusterState) error {

	cs.cluster.Status.Capacity = types.ClusterResources{}
	cs.cluster.Status.Allocated = types.ClusterResources{}

	for _, n := range cs.node.list {

		cs.cluster.Status.Allocated.CPU += n.Status.Allocated.CPU
		cs.cluster.Status.Allocated.RAM += n.Status.Allocated.RAM
		cs.cluster.Status.Allocated.Storage += n.Status.Allocated.Storage
		cs.cluster.Status.Allocated.Containers += n.Status.Allocated.Containers
		cs.cluster.Status.Allocated.Pods += n.Status.Allocated.Pods

		cs.cluster.Status.Capacity.CPU += n.Status.Capacity.CPU
		cs.cluster.Status.Capacity.RAM += n.Status.Capacity.RAM
		cs.cluster.Status.Capacity.Storage += n.Status.Capacity.Storage
		cs.cluster.Status.Capacity.Containers += n.Status.Capacity.Containers
		cs.cluster.Status.Capacity.Pods += n.Status.Capacity.Pods

	}

	if err := distribution.NewClusterModel(context.Background(), envs.Get().GetStorage()).Set(cs.cluster); err != nil {
		log.Errorf("%s: cluster update status error: %s", logPrefix, err.Error())
		return err
	}

	return nil
}
