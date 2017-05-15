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

package pod

import (
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/scheduler/context"
)

func Provision(p *types.Pod) error {

	var (
		log = context.Get().GetLogger()
		stg = context.Get().GetStorage()

		node   *types.Node
		memory = int64(0)
	)

	log.Debugf("Allocate node for pod: %s", p.Meta.Name)

	nodes, err := stg.Node().List(context.Get().Background())
	if err != nil {
		log.Errorf("Node: allocate: get nodes error: %s", err.Error())
		return err
	}

	for _, c := range p.Spec.Containers {
		memory += c.Quota.Memory
	}

	for _, node = range nodes {
		log.Debugf("Node: Allocate: available memory %d", node.State.Capacity)
		if node.State.Capacity.Memory > memory {
			break
		}
	}

	if node == nil {
		log.Error("Node: Allocate: Available node not found")
		return errors.New(errors.NodeNotFound)
	}

	stg.Node().InsertPod(context.Get().Background(), &node.Meta, &types.PodNodeSpec{
		Meta:  p.Meta,
		Spec:  p.Spec,
		State: p.State,
	})

	return nil
}
