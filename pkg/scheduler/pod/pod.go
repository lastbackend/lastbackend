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
		ctx = context.Get().Background()

		node   *types.Node
		memory = int64(0)
	)

	log.Debugf("Allocate node for pod: %s", p.Meta.Name)

	nodes, err := stg.Node().List(ctx)
	if err != nil {
		log.Errorf("Node: allocate: get nodes error: %s", err.Error())
		return err
	}

	for _, c := range p.Spec.Containers {
		memory += c.Quota.Memory
	}

	for _, n := range nodes {
		log.Debugf("Node: Allocate: available memory %d", n.State.Capacity)
		if n.State.Capacity.Memory > memory && n.Alive {
			node = n
			break
		}
	}

	if node == nil {
		log.Error("Node: Allocate: Available node not found")
		return errors.New(errors.NodeNotFound)
	}

	spec := &types.PodNodeSpec{
		Meta:  p.Meta,
		State: p.State,
		Spec:  p.Spec,
	}

	if err := stg.Node().InsertPod(ctx, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec add: insert spec to node err: %s", err.Error())
		return err
	}

	return nil
}

func Update(p *types.Pod) error {
	var (
		stg = context.Get().GetStorage()
		log = context.Get().GetLogger()
		ctx = context.Get().Background()
	)

	node, err := stg.Node().Get(ctx, p.Meta.Hostname)
	if err != nil {
		log.Errorf("Node: Pod spec update: find node err: %s", err.Error())
		return err
	}

	spec := &types.PodNodeSpec{
		Meta:  p.Meta,
		State: p.State,
		Spec:  p.Spec,
	}

	log.Debug("Update pod spec from node")
	if err := stg.Node().UpdatePod(ctx, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec update: update pod spec err: %s", err.Error())
		return err
	}

	return nil
}

func Remove(p *types.Pod) error {
	var (
		stg = context.Get().GetStorage()
		log = context.Get().GetLogger()
		ctx = context.Get().Background()
	)

	node, err := stg.Node().Get(ctx, p.Meta.Hostname)
	if err != nil {
		log.Errorf("Node: Pod spec remove: find node err: %s", err.Error())
		return err
	}

	spec := &types.PodNodeSpec{
		Meta:  p.Meta,
		State: p.State,
		Spec:  p.Spec,
	}

	log.Debug("Remove pod spec from node")
	if err := stg.Node().RemovePod(ctx, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	return nil
}
