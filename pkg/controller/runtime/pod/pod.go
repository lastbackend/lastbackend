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

package pod

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"sort"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/v3/store"
)

func Provision(p *types.Pod) error {

	log.Debugf("scheduler:pod:controller:provision: start pod provision: %s", p.Meta.SelfLink)

	var (
		stg = envs.Get().GetStorage()

		memory = int64(0)
		node   *types.Node
	)

	pm := distribution.NewPodModel(context.Background(), stg)
	if d, err := pm.Get(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment, p.Meta.Name); d == nil || err != nil {
		if d == nil {
			return errors.New(store.ErrEntityNotFound)
		}
		log.Errorf("scheduler:pod:controller:provision: get pod error: %s", err.Error())
		return err
	}

	nm := distribution.NewNodeModel(context.Background(), stg)

	if p.Meta.Node != "" {
		log.Debugf("Node: [%s]: handle scheduled pos: %s", p.Meta.Node, p.SelfLink())
		n, err := nm.Get(p.Meta.Node)
		if err != nil {
			log.Errorf("Node: find node err: %s", err.Error())
			return err
		}

		if n == nil {
			log.Errorf("Node: not attached to pod")
			return errors.New(errors.NodeNotFound)
		}

		if p.Spec.State.Destroy {
			log.Debugf("Node: [%s]: update pod spec: %s", n.Meta.Name, p.SelfLink())
			if err := nm.UpdatePod(n, p); err != nil {
				log.Errorf("Node: update pod spec err: %s", err.Error())
				return err
			}
			return nil
		}

		if err := nm.InsertPod(n, p); err != nil {
			log.Errorf("Node: update pod spec err: %s", err.Error())
			return err
		}
	}

	if p.Spec.State.Destroy {
		return nil
	}

	log.Debugf("Allocate node for pod: %s", p.Meta.Name)

	nodes, err := nm.List()
	if err != nil {
		log.Errorf("Node: allocate: get nodes error: %s", err.Error())
		return err
	}

	for _, c := range p.Spec.Template.Containers {
		memory += c.Resources.Request.RAM
	}

	var nl []*types.Node
	for _, n := range nodes {
		nl = append(nl, n)
	}

	sort.Slice(nl, func(i, j int) bool {
		n1 := nl[i].Status.Capacity.Memory - nl[i].Status.Allocated.Memory
		n2 := nl[j].Status.Capacity.Memory - nl[j].Status.Allocated.Memory
		return n2 < n1
	})

	for _, n := range nl {

		if !n.Online {
			continue
		}

		ram := n.Status.Capacity.Memory - n.Status.Allocated.Memory
		pds := n.Status.Capacity.Pods - n.Status.Allocated.Pods
		cns := n.Status.Capacity.Containers - n.Status.Allocated.Containers

		if ram <= memory {
			continue
		}

		if pds == 0 {
			continue
		}

		if cns <= len(p.Spec.Template.Containers) {
			continue
		}

		node = n
		break
	}

	if node == nil {

		log.Debug("Node: Allocate: Available node not found")

		if err := distribution.NewPodModel(context.Background(), stg).SetStatus(p, &types.PodStatus{
			Stage:   types.StateError,
			Message: errors.NodeNotFound,
		}); err != nil {
			log.Errorf("set pod state error: %s", err.Error())
			return err
		}

		return errors.New(errors.NodeNotFound)
	}

	p.Meta.Node = node.Meta.Name

	if err := pm.SetNode(p, node); err != nil {
		log.Errorf("Pod: attach node to spec err: %s", err.Error())
		return err
	}

	if err := nm.InsertPod(node, p); err != nil {
		log.Errorf("Node: Pod spec add: insert spec to node err: %s", err.Error())
		return err
	}

	return nil
}

func HandleStatus(p *types.Pod) error {

	var (
		stg       = envs.Get().GetStorage()
		logPrefix = "controller:pod:controller:status>"
		status    = make(map[string]int)
		message   string
	)

	pm := distribution.NewPodModel(context.Background(), stg)
	dm := distribution.NewDeploymentModel(context.Background(), stg)

	d, err := dm.Get(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
	if err != nil {
		log.Errorf("%s get deployment err: %v", logPrefix, err)
		return err
	}
	if d == nil {
		log.Warnf("%s: deployment already removed: [%s:%s:%s]", logPrefix, p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
		if err := pm.Remove(context.Background(), p); err != nil {
			log.Errorf("%s: remove pod [%s] err: %v", logPrefix, p.SelfLink())
			return err
		}
		return nil
	}

	pl, err := pm.ListByDeployment(p.Meta.Namespace, p.Meta.Service, p.Meta.Deployment)
	if err != nil {
		log.Errorf("%s> get pod list err: %v", logPrefix, err)
		return err
	}

	for _, ps := range pl {

		switch ps.Status.Stage {

		case types.StateError:
			status[types.StateError] += 1
			// TODO: check if many pods contains different errors: create an error map
			message = ps.Status.Message
			break
		case types.StateProvision:
			status[types.StateProvision] += 1
			break
		case types.StatePull:
			status[types.StateProvision] += 1
			break
		case types.StateCreated:
			status[types.StateProvision] += 1
			break
		case types.StateStarting:
			status[types.StateProvision] += 1
			break
		case types.StateStarted:
			status[types.StateRunning] += 1
			break
		case types.StateRunning:
			status[types.StateRunning] += 1
			break
		case types.StateStopped:
			status[types.StateStopped] += 1
			break
		case types.StateExited:
			status[types.StateStopped] += 1
			break
		case types.StateDestroy:
			status[types.StateDestroy] += 1
			break
		case types.StateDestroyed:
			status[types.StateDestroyed] += 1
			break
		case types.StateWarning:
			status[types.StateWarning] += 1
			break
		}
	}

	opts := new(request.DeploymentUpdateOptions)

	switch true {
	case status[types.StateError] > 0:
		opts.Status.State = types.StateError
		opts.Status.Message = message
		break
	case status[types.StateProvision] > 0:
		opts.Status.State = types.StateProvision
		opts.Status.Message = ""
		break
	case status[types.StateDestroy] > 0:
		opts.Status.State = types.StateDestroy
		opts.Status.Message = ""
		break
	case status[types.StateWarning] > 0 && d.Status.State != types.StateDestroy:
		opts.Status.State = types.StateWarning
		opts.Status.Message = p.Status.Message
		break
	case status[types.StateRunning] == d.Spec.Replicas:
		opts.Status.State = types.StateRunning
		break
	case status[types.StateStopped] == d.Spec.Replicas:
		opts.Status.State = types.StateStopped
		break
	case status[types.StateDestroyed] == 1 && p.Status.Stage == types.StateDestroyed:
		opts.Status.State = types.StateDestroyed
		break
	}

	if p.Status.Stage == types.StateDestroyed {
		if err := pm.Remove(context.Background(), p); err != nil {
			log.Errorf("%s: remove pod [%s] err: %v", logPrefix, p.SelfLink(), err)
			return err
		}
	}

	if err := dm.Update(d, opts); err != nil {
		log.Errorf("%s> update deployment err: %v", logPrefix, err)
		return err
	}

	return nil
}
