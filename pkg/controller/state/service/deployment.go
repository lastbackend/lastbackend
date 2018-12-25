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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"

	"time"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logDeploymentPrefix = "state:observer:deployment"

func deploymentObserve(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

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

	log.V(logLevel).Debugf("%s:> observe state: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if err := endpointCheck(ss); err != nil {
		return err
	}

	if err := serviceStatusState(ss); err != nil {
		return err
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	return nil
}

func handleDeploymentStateCreated(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> handleDeploymentStateCreated: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	check, err := deploymentCheckDependencies(ss, d)
	if err != nil {
		log.Errorf("%s:> handle deployment check deps: %s, err: %s", logDeploymentPrefix, d.SelfLink(), err.Error())
		return err
	}

	if !check {
		d.Status.State = types.StateWaiting
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		if err := dm.Update(d); err != nil {
			log.Errorf("%s:> handle deployment create, deps update: %s, err: %s", logDeploymentPrefix, d.SelfLink(), err.Error())
			return err
		}
		return nil
	}

	if err := deploymentCheckSelectors(ss, d); err != nil {
		d.Status.State = types.StateError
		d.Status.Message = err.Error()
		dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
		if err := dm.Update(d); err != nil {
			log.Errorf("%s:> handle deployment create, deps update: %s, err: %s", logDeploymentPrefix, d.SelfLink(), err.Error())
			return err
		}
		return nil
	}

	if err := deploymentPodProvision(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleDeploymentStateProvision(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> handleDeploymentStateProvision: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if err := deploymentPodProvision(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleDeploymentStateReady(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> handleDeploymentStateReady: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.active != nil {
		if ss.deployment.active.SelfLink() != d.SelfLink() {
			if err := deploymentDestroy(ss, ss.deployment.active); err != nil {
				log.Errorf("%s", err.Error())
				return err
			}
		}
	}

	ss.deployment.active = d

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

	//if ss.deployment.active.SelfLink() != d.SelfLink() {
	//	return deploymentDestroy(ss, d)
	//}

	return nil
}

func handleDeploymentStateError(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> handleDeploymentStateError: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.active == nil {
		ss.deployment.provision = nil
		ss.deployment.active = d
	}

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

	//if ss.deployment.active.SelfLink() != d.SelfLink() {
	//	return deploymentDestroy(ss, d)
	//}

	return nil
}

func handleDeploymentStateDegradation(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> handleDeploymentStateDegradation: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if err := deploymentPodProvision(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	if ss.deployment.active == nil {
		ss.deployment.provision = nil
		ss.deployment.active = d
		return nil
	}

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

	//if ss.deployment.active.SelfLink() != d.SelfLink() {
	//	return deploymentDestroy(ss, d)
	//}

	return nil
}

func handleDeploymentStateDestroy(ss *ServiceState, d *types.Deployment) error {

	log.V(logLevel).Debugf("%s:> handleDeploymentStateDestroy: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

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

	log.V(logLevel).Debugf("%s:> handleDeploymentStateDestroyed: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

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

	ss.DelDeployment(d)
	return nil
}

func deploymentSpecValidate(d *types.Deployment, svc *types.Service) bool {
	return d.Spec.Template.Updated.Equal(svc.Spec.Template.Updated) && d.Spec.Selector.Updated.Equal(svc.Spec.Selector.Updated)
}

// serviceCheckDependencies function - check if service can provisioned or should wait for dependencies
func deploymentCheckDependencies(ss *ServiceState, d *types.Deployment) (bool, error) {

	var (
		ctx  = context.Background()
		stg  = envs.Get().GetStorage()
		vm   = distribution.NewVolumeModel(ctx, stg)
		sm   = distribution.NewSecretModel(ctx, stg)
		cm   = distribution.NewConfigModel(ctx, stg)
		deps = types.DeploymentStatusDependencies{
			Volumes: make(map[string]types.DeploymentStatusDependency, 0),
			Secrets: make(map[string]types.DeploymentStatusDependency, 0),
			Configs: make(map[string]types.DeploymentStatusDependency, 0),
		}
	)

	volumesRequiredList := make(map[string]bool, 0)
	secretsRequiredList := make(map[string]bool, 0)
	configsRequiredList := make(map[string]bool, 0)
	for _, v := range d.Spec.Template.Volumes {
		if v.Volume.Name != types.EmptyString {
			volumesRequiredList[v.Volume.Name] = true
		}
		if v.Secret.Name != types.EmptyString {
			secretsRequiredList[v.Secret.Name] = true
		}
		if v.Config.Name != types.EmptyString {
			configsRequiredList[v.Config.Name] = true
		}
	}

	for _, c := range d.Spec.Template.Containers {
		for _, e := range c.EnvVars {
			if e.Secret.Name != types.EmptyString {
				secretsRequiredList[e.Secret.Name] = true
			}

			if e.Config.Name != types.EmptyString {
				configsRequiredList[e.Config.Name] = true
			}
		}
	}

	if len(volumesRequiredList) != 0 {

		vl, err := vm.ListByNamespace(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:> service check deps err: %s", logServicePrefix, err.Error())
			return false, err
		}

		for vr := range volumesRequiredList {
			var f = false

			for _, v := range vl.Items {
				if vr == v.Meta.Name {
					f = true
					deps.Volumes[vr] = types.DeploymentStatusDependency{
						Name:   vr,
						Type:   types.KindVolume,
						Status: v.Status.State,
					}
				}
			}

			if !f {
				deps.Volumes[vr] = types.DeploymentStatusDependency{
					Name:   vr,
					Type:   types.KindVolume,
					Status: types.StateNotReady,
				}
			}
		}
	}

	if len(secretsRequiredList) != 0 {

		sl, err := sm.List(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:> service check deps err: %s", logServicePrefix, err.Error())
			return false, err
		}

		for sr := range secretsRequiredList {
			var f = false

			for _, s := range sl.Items {
				if sr == s.Meta.Name {
					f = true
					deps.Secrets[sr] = types.DeploymentStatusDependency{
						Name:   sr,
						Type:   types.KindSecret,
						Status: types.StateReady,
					}
				}
			}

			if !f {
				deps.Secrets[sr] = types.DeploymentStatusDependency{
					Name:   sr,
					Type:   types.KindSecret,
					Status: types.StateNotReady,
				}
			}
		}
	}

	if len(configsRequiredList) != 0 {

		cl, err := cm.List(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:> service check deps err: %s", logServicePrefix, err.Error())
			return false, err
		}

		for cr := range configsRequiredList {
			var f = false

			for _, c := range cl.Items {
				if cr == c.Meta.Name {
					f = true
					deps.Configs[cr] = types.DeploymentStatusDependency{
						Name:   cr,
						Type:   types.KindConfig,
						Status: types.StateReady,
					}
				}
			}

			if !f {
				deps.Configs[cr] = types.DeploymentStatusDependency{
					Name:   cr,
					Type:   types.KindConfig,
					Status: types.StateNotReady,
				}
			}
		}
	}

	d.Status.Dependencies = deps
	if !d.Status.CheckDeps() {
		d.Status.State = types.StateWaiting
		return false, nil
	}

	return true, nil
}

// deploymentCheckSelectors function - handles provided selectors to match nodes
func deploymentCheckSelectors(ss *ServiceState, d *types.Deployment) (err error) {

	var (
		ctx = context.Background()
		stg = envs.Get().GetStorage()
		vm  = distribution.NewVolumeModel(ctx, stg)
		vc  = make(map[string]string, 0)
	)

	for _, v := range d.Spec.Template.Volumes {
		if v.Volume.Name != types.EmptyString {
			vc[v.Volume.Name] = v.Name
		}
	}

	if len(vc) > 0 {

		var node string

		vl, err := vm.ListByNamespace(d.Meta.Namespace)
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> create deployment, volume list err: %s", logPrefix, err.Error())
			return err
		}

		for name := range vc {

			var f = false

			for _, v := range vl.Items {

				if v.Meta.Name != name {
					continue
				}

				f = true

				if v.Status.State != types.StateReady {
					log.V(logLevel).Errorf("%s:create:> create deployment err: volume is not ready yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotReady(v.Meta.Name)
				}

				if v.Meta.Node == types.EmptyString {
					log.V(logLevel).Errorf("%s:create:> create deployment err: volume is not provisioned yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotProvisioned(v.Meta.Name)
				}

				if node == types.EmptyString {
					node = v.Meta.Node
				} else {
					if node != v.Meta.Node {
						return errors.New(v.Meta.Name).Volume().DifferentNodes()
					}
				}
			}

			if !f {
				log.V(logLevel).Errorf("%s:create:> create deployment err: volume is not found: %s", logPrefix, name)
				return errors.New(name).Volume().NotFound(name)
			}
		}

		if node != types.EmptyString {

			if d.Spec.Selector.Node != types.EmptyString {
				if d.Spec.Selector.Node != node {
					return errors.New("spec.selector.node not matched with attached volumes")
				}

				return nil
			}

			d.Spec.Selector.Node = node
		}

	}

	return nil
}

// deploymentPodProvision - handles deployment provision logic
// based on current deployment state and current pod list of provided deployment
func deploymentPodProvision(ss *ServiceState, d *types.Deployment) (err error) {

	t := d.Meta.Updated

	var (
		provision = false
	)

	defer func() {
		if err == nil {
			err = deploymentUpdate(d, t)
		}
	}()

	var (
		st = []string{
			types.StateError,
			types.StateWarning,
			types.StateCreated,
			types.StateProvision,
			types.StateReady,
		}
		pm = distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
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

				if p.Meta.Node != types.EmptyString {

					m, e := pm.ManifestGet(p.Meta.Node, p.SelfLink())
					if err != nil {
						err = e
						return e
					}

					if m == nil {
						if err = podManifestPut(p); err != nil {
							return err
						}
					}

				}

				total++
			}

			if _, ok := state[p.Status.State]; !ok {
				state[p.Status.State] = make([]*types.Pod, 0)
			}

			state[p.Status.State] = append(state[p.Status.State], p)
		}

		if d.Spec.Replicas == total {
			break
		}

		if d.Spec.Replicas > total {
			log.V(logLevel).Debugf("create additional replica: %d -> %d", total, d.Spec.Replicas)
			p, err := podCreate(d)
			if err != nil {
				log.Errorf("%s", err.Error())
				return err
			}
			pods[p.SelfLink()] = p
			provision = true
			continue
		}

		if d.Spec.Replicas < total {
			log.V(logLevel).Debugf("remove unneeded replica: %d -> %d", total, d.Spec.Replicas)
			for _, s := range st {

				if len(state[s]) > 0 {

					p := state[s][0]

					if err := podDestroy(ss, p); err != nil {
						log.Errorf("%s", err.Error())
						return err
					}

					provision = true
					break
				}
			}
		}

	}

	if provision {
		if d.Status.State != types.StateProvision {
			d.Status.State = types.StateProvision
			d.Meta.Updated = time.Now()
		}
	}

	return nil
}

func deploymentCreate(svc *types.Service, version int) (*types.Deployment, error) {

	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())

	d, err := dm.Create(svc, version)
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
		if err == nil {
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

		if p.Status.State != types.StateDestroy {
			if err := podDestroy(ss, p); err != nil {
				return err
			}
		}

		if p.Status.State == types.StateDestroyed {
			if err := podRemove(ss, p); err != nil {
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
	if err := dm.Remove(d); err != nil {
		return err
	}

	return nil
}

func deploymentScale(d *types.Deployment, replicas int) error {
	d.Status.State = types.StateProvision
	d.Spec.Replicas = replicas
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	return dm.Update(d)
}

func deploymentStatusState(d *types.Deployment, pl map[string]*types.Pod) (err error) {

	log.V(logLevel).Debugf("%s:> deploymentStatusState: start: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	t := d.Meta.Updated
	defer func() {
		if err == nil {
			err = deploymentUpdate(d, t)
		}
	}()

	var (
		state   = make(map[string]int)
		message string
		running int
	)

	for _, p := range pl {
		state[p.Status.State]++
		if p.Status.State == types.StateError {
			message = p.Status.Message
		}

		if p.Status.Running {
			running++
		}
	}

	switch d.Status.State {
	case types.StateCreated:
		break
	case types.StateProvision:

		if _, ok := state[types.StateReady]; ok && running == len(pl) {
			d.Status.State = types.StateReady
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && running == 0 {
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

		d.Status.State = types.StateDegradation
		d.Status.Message = types.EmptyString
		d.Meta.Updated = time.Now()
		break
	case types.StateReady:

		if _, ok := state[types.StateReady]; ok && running == len(pl) {
			break
		}

		if _, ok := state[types.StateProvision]; ok {
			d.Status.State = types.StateProvision
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && running == 0 {
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

		if _, ok := state[types.StateReady]; ok && running == len(pl) {
			d.Status.State = types.StateReady
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[types.StateError]; ok && running == 0 {
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

		if _, ok := state[types.StateReady]; ok && running == len(pl) {
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

		if _, ok := state[types.StateError]; ok && running == 0 {
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

	log.V(logLevel).Debugf("%s:> deploymentStatusState: finish: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	return nil
}
