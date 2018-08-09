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
	"time"
)

const logDeploymentPrefix = "state:observer:deployment"

func deploymentObserve(ss *ServiceState, d *types.Deployment) error {

	if _, ok := ss.pod.list[d.SelfLink()]; !ok {
		ss.pod.list[d.SelfLink()] = make(map[string]*types.Pod)
	}

	switch d.Status.State {
	case types.StateCreated:
		if err := handleDeploymentStateCreated(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state create err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case types.StateProvision:
		if err := handleDeploymentStateProvision(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state provision err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case types.StateReady:
		if err := handleDeploymentStateReady(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state ready err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case types.StateError:
		if err := handleDeploymentStateError(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state error err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case types.StateDegradation:
		if err := handleDeploymentStateDegradation(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state degradation err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case types.StateDestroy:
		if err := handleDeploymentStateDestroy(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state destroy err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case types.StateDestroyed:
		if err := handleDeploymentStateDestroyed(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state destroyed err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	}

	if d.Status.State == types.StateDestroyed {
		delete(ss.deployment.list, d.SelfLink())
	} else {
		ss.deployment.list[d.SelfLink()] = d
	}

	serviceStatusState(ss)

	return nil
}

func handleDeploymentStateCreated(ss *ServiceState, d *types.Deployment) error {


	if err := deploymentPodProvision(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleDeploymentStateProvision(ss *ServiceState, d *types.Deployment) error {

	if err := deploymentPodProvision(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleDeploymentStateReady(ss *ServiceState, d *types.Deployment) error {

	if ss.deployment.active != nil {
		log.Info("call active deployment down")
		if err := deploymentDestroy(ss, ss.deployment.active); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	ss.deployment.active = d

	if ss.deployment.provision.SelfLink() == d.SelfLink() {
		ss.deployment.provision = nil
	}

	serviceStatusState(ss)

	return nil
}

func handleDeploymentStateError(ss *ServiceState, d *types.Deployment) error {

	if ss.deployment.active == nil {
		ss.deployment.provision = nil
		ss.deployment.active = d
	}

	serviceStatusState(ss)
	return nil
}

func handleDeploymentStateDegradation(ss *ServiceState, d *types.Deployment) error {

	if ss.deployment.active.SelfLink() == d.SelfLink() {
		serviceStatusState(ss)
	}

	return nil
}

func handleDeploymentStateDestroy(ss *ServiceState, d *types.Deployment) error {

	if err := deploymentDestroy(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	if d.Status.State == types.StateDestroyed {
		return handleDeploymentStateDestroyed(ss, d)
	}

	return nil
}

func handleDeploymentStateDestroyed(ss *ServiceState, d *types.Deployment) error {

	link := d.SelfLink()

	if _, ok := ss.pod.list[link]; ok && len(ss.pod.list[link]) > 0 {

		if err := deploymentDestroy(ss, d); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}

		d.Status.State = types.StateDestroy
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		return dm.Update(d)
	}

	if err := deploymentRemove(d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

	if ss.deployment.active != nil {
		if ss.deployment.active.SelfLink() == d.SelfLink() {
			ss.deployment.active = nil
		}
	}

	return nil
}

func deploymentSpecValidate(d *types.Deployment, spec types.SpecTemplate) bool {
	return !d.Spec.Template.Updated.Before(spec.Updated)
}

// deploymentPodProvision - handles deployment provision logic
// based on current deployment state and current pod list of provided deployment
func deploymentPodProvision(ss *ServiceState, d *types.Deployment) (err error) {

	t := d.Meta.Updated

	defer func() {
		if err != nil {
			err = deploymentUpdate(d, t)
		}
	}()

	if d.Status.State != types.StateProvision {
		d.Status.State = types.StateProvision
		d.Meta.Updated = time.Now()
	}

	var (
		st       = []string{
			types.StateError,
			types.StateWarning,
			types.StateCreated,
			types.StateProvision,
			types.StateReady,
		}
	)

	pods, ok := ss.pod.list[d.SelfLink()]
	if !ok {
		pods = make(map[string]*types.Pod, 0)
	}

	for {

		var (
			total int
			state = make(map[string][]*types.Pod)
		)

		for _, p := range pods {

			if p.Status.State != types.StateDestroy && p.Status.State != types.StateDestroyed {
				total++
			}

			if _, ok := state[p.Status.State]; !ok {
				state[p.Status.State] = make([]*types.Pod, 0)
			}

			state[p.Status.State] = append(state[p.Status.State], p)
		}

		if d.Spec.Replicas == total {
			return nil
		}

		if d.Spec.Replicas > total {
			log.Debugf("create additional replica: %d -> %d", total, d.Spec.Replicas)
			p, err := podCreate(d)
			if err != nil {
				log.Errorf("%s", err.Error())
				return err
			}
			pods[p.SelfLink()] = p
			continue
		}

		if d.Spec.Replicas < total {
			log.Debugf("remove unneeded replica: %d -> %d", total, d.Spec.Replicas)
			for _, s := range st {

				if len(state[s]) > 0 {

					p := state[s][0]

					if err := podDestroy(ss, p); err != nil {
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

func deploymentCreate(svc *types.Service) (*types.Deployment, error) {

	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())

	d, err := dm.Create(svc)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func deploymentUpdate(d *types.Deployment, timestamp time.Time) error {
	if timestamp.Before(d.Meta.Updated) {
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		if err := dm.Update(d); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func deploymentDestroy(ss *ServiceState, d *types.Deployment) (err error) {

	t := d.Meta.Updated

	defer func() {
		if err != nil {
			err = deploymentUpdate(d, t)
		}
	}()

	if d.Status.State != types.StateDestroy {
		d.Status.State = types.StateDestroy
		d.Meta.Updated = time.Now()
	}

	pl, ok := ss.pod.list[d.SelfLink()]
	if !ok {
		d.Status.State = types.StateDestroyed
		d.Meta.Updated = time.Now()
		return nil
	}

	for _, p := range pl {

		if p.Status.State == types.StateDestroyed {
			if err := podRemove(ss, p); err != nil {
				return err
			}
			continue
		}

		if p.Status.State != types.StateDestroy {
			if err := podDestroy(ss, p); err != nil {
				return err
			}
		}
	}

	if len(pl) == 0 {
		d.Status.State = types.StateDestroyed
		d.Meta.Updated = time.Now()
		return nil
	}


	return nil
}

func deploymentRemove(d *types.Deployment) error {
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	return dm.Remove(d)
}

func deploymentScale(d *types.Deployment, replicas int) error {
	d.Status.State = types.StateProvision
	d.Spec.Replicas = replicas
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	return dm.Update(d)
}

func deploymentStatusState(d *types.Deployment, pl map[string]*types.Pod) (err error) {

	t := d.Meta.Updated

	defer func() {
		if err != nil {
			err = deploymentUpdate(d, t)
		}
	}()

	var (
		state = make(map[string]int)
		message string
	)

	for _, p := range pl {
		state[p.Status.State]++
		if p.Status.State == types.StateError {
			message = p.Status.Message
		}
	}

	switch d.Status.State {
	case types.StateCreated:
		break
	case types.StateProvision:

		if _, ok := state[types.StateReady]; ok && state[types.StateReady] == len(pl) {
			d.Status.State = types.StateReady
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && state[types.StateError] == len(pl) {
			d.Status.State = types.StateError
			d.Status.Message = message
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateProvision]; ok {
			d.Status.State = types.StateProvision
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateCreated]; ok {
			d.Status.State = types.StateProvision
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		d.Status.State = types.StateDegradation
		d.Status.Message = types.EmptyString
		d.Meta.Updated = time.Now()
		break
	case types.StateReady:

		if _, ok := state[types.StateReady]; ok && state[types.StateReady] == len(pl) {
			break
		}

		if _, ok := state[types.StateProvision]; ok {
			d.Status.State = types.StateProvision
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && state[types.StateError] == len(pl) {
			d.Status.State = types.StateError
			d.Status.Message = message
			d.Meta.Updated = time.Now()
			break
		}

		d.Status.State = types.StateDegradation
		d.Status.Message = types.EmptyString
		d.Meta.Updated = time.Now()

		break
	case types.StateError:

		if _, ok := state[types.StateReady]; ok && state[types.StateReady] == len(pl) {
			d.Status.State = types.StateReady
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && state[types.StateError] == len(pl) {
			break
		}

		if _, ok := state[types.StateProvision]; ok {
			d.Status.State = types.StateProvision
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		d.Status.State = types.StateDegradation
		d.Status.Message = types.EmptyString
		d.Meta.Updated = time.Now()

		break
	case types.StateDegradation:

		if _, ok := state[types.StateReady]; ok && state[types.StateReady] == len(pl) {
			d.Status.State = types.StateReady
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateProvision]; ok {
			d.Status.State = types.StateProvision
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && state[types.StateError] == len(pl) {
			d.Status.State = types.StateError
			d.Status.Message = message
			d.Meta.Updated = time.Now()
			break
		}

		break
	case types.StateDestroy:
		if len(pl) == 0 {
			d.Status.State = types.StateDestroyed
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
		}
		break
	case types.StateDestroyed:
		break
	}

	return nil
}
