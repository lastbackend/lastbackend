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

package service

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

// DeploymentProvision - handles deployment provision logic
// based on current deployment state and current pod list of provided deployment
func DeploymentProvision(d *types.Deployment, pods map[string]*types.Pod) error {

	if d.Status.State != types.StateProvision {
		d.Status.State = types.StateProvision
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		return dm.Update(d)
	}

	var (
		replicas = d.Spec.Replicas
		st       = []string{types.StateError, types.StateWarning, types.StateCreated, types.StateProvision, types.StateReady}
	)

	for {

		var (
			total int
			state = make(map[string][]*types.Pod)
		)

		log.Debugf(">> len pods: %d", len(pods))
		for _, p := range pods {

			if p.Status.State != types.StateDestroy && p.Status.State != types.StateDestroyed {
				total++
			}

			if _, ok := state[p.Status.State]; !ok {
				state[p.Status.State] = make([]*types.Pod, 0)
			}

			state[p.Status.State] = append(state[p.Status.State], p)
		}

		d.Spec.Replicas = total
		log.Debugf(">> %d", d.Spec.Replicas)

		if d.Spec.Replicas == replicas {
			return nil
		}

		if d.Spec.Replicas < replicas {
			log.Debugf("create additional replica: %d -> %d", d.Spec.Replicas, replicas)
			p, err := PodCreate(d)
			if err != nil {
				log.Errorf("%s", err.Error())
				return err
			}
			pods[p.Meta.Name] = p
			continue
		}

		if d.Spec.Replicas > replicas {
			log.Debugf("remove unneeded replica: %d -> %d", d.Spec.Replicas, replicas)
			for _, s := range st {

				if len(state[s]) > 0 {

					p := state[s][0]

					if err := PodDestroy(p); err != nil {
						log.Errorf("%s", err.Error())
						return err
					}

					break
				}
			}
		}

	}

	return nil
}

func DeploymentCreate(svc *types.Service) (*types.Deployment, error) {

	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())

	d, err := dm.Create(svc)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func DeploymentDestroy(_ *types.Deployment, pl map[string]*types.Pod) error {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())

	for _, p := range pl {

		if p.Status.State == types.StateDestroyed {
			continue
		}

		if p.Status.State != types.StateDestroy {
			if err := pm.Destroy(p); err != nil {
				return err
			}
		}
	}

	return nil
}

func DeploymentCancel(_ *types.Deployment, pl map[string]*types.Pod) error {
	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())

	for _, p := range pl {

		if p.Status.State == types.StateDestroyed {
			continue
		}

		if p.Status.State != types.StateDestroy {
			if err := pm.Destroy(p); err != nil {
				return err
			}
		}
	}

	return nil
}

func DeploymentRemove(d *types.Deployment) error {
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	return dm.Remove(d)
}

func DeploymentScale(d *types.Deployment, replicas int) error {
	d.Status.State = types.StateProvision
	d.Spec.Replicas = replicas
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	return dm.Update(d)
}

func DeploymentSync() {

}
