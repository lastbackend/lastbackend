//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"fmt"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logDeploymentPrefix = "state:observer:deployment"

func deploymentObserve(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> observe start: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if _, ok := ss.pod.list[d.SelfLink().String()]; !ok {
		ss.pod.list[d.SelfLink().String()] = make(map[string]*models.Pod)
	}

	switch d.Status.State {
	case models.StateCreated:
		if err := handleDeploymentStateCreated(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state create err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case models.StateProvision:
		if err := handleDeploymentStateProvision(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state provision err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case models.StateReady:
		if err := handleDeploymentStateReady(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state ready err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case models.StateError:
		if err := handleDeploymentStateError(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state error err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case models.StateDegradation:
		if err := handleDeploymentStateDegradation(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state degradation err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case models.StateDestroy:
		if err := handleDeploymentStateDestroy(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state destroy err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	case models.StateDestroyed:
		if err := handleDeploymentStateDestroyed(ss, d); err != nil {
			log.Errorf("%s:> handle deployment state destroyed err: %s", logDeploymentPrefix, err.Error())
			return err
		}
		break
	}

	if d.Status.State == models.StateDestroyed {
		delete(ss.deployment.list, d.SelfLink().String())
	} else {
		ss.deployment.list[d.SelfLink().String()] = d
	}

	log.Debugf("%s:> observe state: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if err := endpointCheck(ss); err != nil {
		return err
	}

	if err := serviceStatusState(ss); err != nil {
		return err
	}

	log.Debugf("%s:> observe finish: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	return nil
}

func handleDeploymentStateCreated(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateCreated: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	ss.deployment.provision = d

	check, err := deploymentCheckDependencies(ss, d)
	if err != nil {
		log.Errorf("%s:> handle deployment check deps: %s, err: %s", logDeploymentPrefix, d.SelfLink(), err.Error())
		return err
	}

	if !check {
		d.Status.State = models.StateWaiting
		dm := service.NewDeploymentModel(context.Background(), ss.storage)
		if err := dm.Update(d); err != nil {
			log.Errorf("%s:> handle deployment create, deps update: %s, err: %s", logDeploymentPrefix, d.SelfLink(), err.Error())
			return err
		}
		return nil
	}

	if err := deploymentCheckSelectors(ss, d); err != nil {
		d.Status.State = models.StateError
		d.Status.Message = err.Error()
		dm := service.NewDeploymentModel(context.Background(), ss.storage)
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

func handleDeploymentStateProvision(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateProvision: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.provision != nil {
		if ss.deployment.provision.Spec.Template.Updated.After(d.Spec.Template.Updated) {
			d.Status.State = models.StateCanceled
			return nil
		}
	}

	ss.deployment.provision = d
	if err := deploymentPodProvision(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleDeploymentStateReady(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateReady: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.active != nil {
		if ss.deployment.active.SelfLink().String() != d.SelfLink().String() {
			if err := deploymentDestroy(ss, ss.deployment.active); err != nil {
				log.Errorf("%s", err.Error())
				return err
			}
		}
	}

	ss.deployment.active = d

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink().String() == d.SelfLink().String() {
			ss.deployment.provision = nil
		}
	}

	if ss.deployment.active.SelfLink().String() != d.SelfLink().String() {
		return deploymentDestroy(ss, d)
	}

	return nil
}

func handleDeploymentStateError(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateError: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

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

func handleDeploymentStateDegradation(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateDegradation: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

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

func handleDeploymentStateDestroy(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateDestroy: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

	if err := deploymentDestroy(ss, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	if d.Status.State == models.StateDestroyed {
		return handleDeploymentStateDestroyed(ss, d)
	}

	return nil
}

func handleDeploymentStateDestroyed(ss *ServiceState, d *models.Deployment) error {

	log.Debugf("%s:> handleDeploymentStateDestroyed: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink().String() == d.SelfLink().String() {
			ss.deployment.provision = nil
		}
	}

	link := d.SelfLink().String()

	if _, ok := ss.pod.list[link]; ok && len(ss.pod.list[link]) > 0 {

		if err := deploymentDestroy(ss, d); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}

		d.Status.State = models.StateDestroy
		dm := service.NewDeploymentModel(context.Background(), ss.storage)
		return dm.Update(d)
	}

	if err := deploymentRemove(ss.storage, d); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	ss.DelDeployment(d)
	return nil
}

func deploymentSpecValidate(d *models.Deployment, svc *models.Service) bool {
	return d.Spec.Template.Updated.Equal(svc.Spec.Template.Updated) && d.Spec.Selector.Updated.Equal(svc.Spec.Selector.Updated)
}

// serviceCheckDependencies function - check if service can provisioned or should wait for dependencies
func deploymentCheckDependencies(ss *ServiceState, d *models.Deployment) (bool, error) {

	var (
		ctx  = context.Background()
		stg  = ss.storage
		vm   = service.NewVolumeModel(ctx, stg)
		sm   = service.NewSecretModel(ctx, stg)
		cm   = service.NewConfigModel(ctx, stg)
		deps = models.StatusDependencies{
			Volumes: make(map[string]models.StatusDependency, 0),
			Secrets: make(map[string]models.StatusDependency, 0),
			Configs: make(map[string]models.StatusDependency, 0),
		}
	)

	volumesRequiredList := make(map[string]bool, 0)
	secretsRequiredList := make(map[string]bool, 0)
	configsRequiredList := make(map[string]bool, 0)
	for _, v := range d.Spec.Template.Volumes {
		if v.Volume.Name != models.EmptyString {
			volumesRequiredList[v.Volume.Name] = true
		}
		if v.Secret.Name != models.EmptyString {
			secretsRequiredList[v.Secret.Name] = true
		}
		if v.Config.Name != models.EmptyString {
			configsRequiredList[v.Config.Name] = true
		}
	}

	for _, c := range d.Spec.Template.Containers {
		for _, e := range c.EnvVars {
			if e.Secret.Name != models.EmptyString {
				secretsRequiredList[e.Secret.Name] = true
			}

			if e.Config.Name != models.EmptyString {
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
					deps.Volumes[vr] = models.StatusDependency{
						Name:   vr,
						Type:   models.KindVolume,
						Status: v.Status.State,
					}
				}
			}

			if !f {
				deps.Volumes[vr] = models.StatusDependency{
					Name:   vr,
					Type:   models.KindVolume,
					Status: models.StateNotReady,
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
					deps.Secrets[sr] = models.StatusDependency{
						Name:   sr,
						Type:   models.KindSecret,
						Status: models.StateReady,
					}
				}
			}

			if !f {
				deps.Secrets[sr] = models.StatusDependency{
					Name:   sr,
					Type:   models.KindSecret,
					Status: models.StateNotReady,
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
					deps.Configs[cr] = models.StatusDependency{
						Name:   cr,
						Type:   models.KindConfig,
						Status: models.StateReady,
					}
				}
			}

			if !f {
				deps.Configs[cr] = models.StatusDependency{
					Name:   cr,
					Type:   models.KindConfig,
					Status: models.StateNotReady,
				}
			}
		}
	}

	d.Status.Dependencies = deps
	if !d.Status.CheckDeps() {
		d.Status.State = models.StateWaiting
		return false, nil
	}

	return true, nil
}

// deploymentCheckSelectors function - handles provided selectors to match nodes
func deploymentCheckSelectors(ss *ServiceState, d *models.Deployment) (err error) {

	var (
		ctx = context.Background()
		vm  = service.NewVolumeModel(ctx, ss.storage)
		vc  = make(map[string]string, 0)
	)

	for _, v := range d.Spec.Template.Volumes {
		if v.Volume.Name != models.EmptyString {
			vc[v.Volume.Name] = v.Name
		}
	}

	if len(vc) > 0 {

		var node string

		vl, err := vm.ListByNamespace(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:create:> create deployment, volume list err: %s", logPrefix, err.Error())
			return err
		}

		for name := range vc {

			var f = false

			for _, v := range vl.Items {

				if v.Meta.Name != name {
					continue
				}

				f = true

				if v.Status.State != models.StateReady {
					log.Errorf("%s:create:> create deployment err: volume is not ready yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotReady(v.Meta.Name)
				}

				if v.Meta.Node == models.EmptyString {
					log.Errorf("%s:create:> create deployment err: volume is not provisioned yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotProvisioned(v.Meta.Name)
				}

				if node == models.EmptyString {
					node = v.Meta.Node
				} else {
					if node != v.Meta.Node {
						return errors.New(v.Meta.Name).Volume().DifferentNodes()
					}
				}
			}

			if !f {
				log.Errorf("%s:create:> create deployment err: volume is not found: %s", logPrefix, name)
				return errors.New(name).Volume().NotFound(name)
			}
		}

		if node != models.EmptyString {

			if d.Spec.Selector.Node != models.EmptyString {
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
func deploymentPodProvision(ss *ServiceState, d *models.Deployment) (err error) {

	t := d.Meta.Updated

	var (
		provision = false
	)

	defer func() {
		if err == nil {
			err = deploymentUpdate(ss.storage, d, t)
		}
	}()

	var (
		st = []string{
			models.StateError,
			models.StateWarning,
			models.StateCreated,
			models.StateProvision,
			models.StateReady,
		}
		pm = service.NewPodModel(context.Background(), ss.storage)
	)

	pods, ok := ss.pod.list[d.SelfLink().String()]
	if !ok {
		pods = make(map[string]*models.Pod, 0)
	}

	for {

		var (
			total int
			state = make(map[string][]*models.Pod)
		)

		for _, p := range pods {

			if p.Status.State != models.StateDestroy && p.Status.State != models.StateDestroyed {

				if p.Meta.Node != models.EmptyString {

					m, e := pm.ManifestGet(p.Meta.Node, p.SelfLink().String())
					if err != nil {
						err = e
						return e
					}

					if m == nil {
						if err = podManifestPut(ss.storage, p); err != nil {
							return err
						}
					}

				}

				if p.Meta.Node == models.EmptyString {
					if err := podProvision(ss, p); err != nil {
						return err
					}
				}

				total++
			}

			if _, ok := state[p.Status.State]; !ok {
				state[p.Status.State] = make([]*models.Pod, 0)
			}

			state[p.Status.State] = append(state[p.Status.State], p)
		}

		if d.Spec.Replicas == total {
			break
		}

		if d.Spec.Replicas > total {
			log.Debugf("create additional replica: %d -> %d", total, d.Spec.Replicas)
			p, err := podCreate(ss.storage, d)
			if err != nil {
				log.Errorf("%s", err.Error())
				return err
			}
			pods[p.SelfLink().String()] = p
			provision = true
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

					provision = true
					break
				}
			}
		}

	}

	if provision {
		if d.Status.State != models.StateProvision {
			d.Status.State = models.StateProvision
			d.Meta.Updated = time.Now()
		}
	}

	return nil
}

func deploymentCreate(stg storage.IStorage, svc *models.Service, version int) (*models.Deployment, error) {

	dm := service.NewDeploymentModel(context.Background(), stg)
	d, err := dm.Create(svc, fmt.Sprintf("v%d", version))
	if err != nil {
		return nil, err
	}

	return d, nil
}

func deploymentUpdate(stg storage.IStorage, d *models.Deployment, timestamp time.Time) error {
	if timestamp.Before(d.Meta.Updated) {
		dm := service.NewDeploymentModel(context.Background(), stg)
		if err := dm.Update(d); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func deploymentDestroy(ss *ServiceState, d *models.Deployment) (err error) {

	t := d.Meta.Updated

	defer func() {
		if err == nil {
			err = deploymentUpdate(ss.storage, d, t)
		}
	}()

	if d.Status.State != models.StateDestroy {
		d.Status.State = models.StateDestroy
		d.Meta.Updated = time.Now()
	}

	pl, ok := ss.pod.list[d.SelfLink().String()]
	if !ok {
		d.Status.State = models.StateDestroyed
		d.Meta.Updated = time.Now()
		return nil
	}

	for _, p := range pl {

		if p.Status.State != models.StateDestroy {
			if err := podDestroy(ss, p); err != nil {
				return err
			}
		}

		if p.Status.State == models.StateDestroyed {
			if err := podRemove(ss, p); err != nil {
				return err
			}
		}
	}

	if len(pl) == 0 {
		d.Status.State = models.StateDestroyed
		d.Meta.Updated = time.Now()
		return nil
	}

	return nil
}

func deploymentRemove(stg storage.IStorage, d *models.Deployment) error {
	dm := service.NewDeploymentModel(context.Background(), stg)
	if err := dm.Remove(d); err != nil {
		return err
	}

	return nil
}

func deploymentScale(stg storage.IStorage, d *models.Deployment, replicas int) error {
	d.Status.State = models.StateProvision
	d.Spec.Replicas = replicas
	dm := service.NewDeploymentModel(context.Background(), stg)
	return dm.Update(d)
}

func deploymentStatusState(stg storage.IStorage, d *models.Deployment, pl map[string]*models.Pod) (err error) {

	log.Debugf("%s:> deploymentStatusState: start: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	t := d.Meta.Updated
	defer func() {
		if err == nil {
			err = deploymentUpdate(stg, d, t)
		}
	}()

	var (
		state   = make(map[string]int)
		message string
		running int
	)

	for _, p := range pl {
		state[p.Status.State]++
		if p.Status.State == models.StateError {
			message = p.Status.Message
		}

		if p.Status.Running {
			running++
		}
	}

	switch d.Status.State {
	case models.StateCreated:
		break
	case models.StateProvision:

		if _, ok := state[models.StateReady]; ok && running == len(pl) {
			d.Status.State = models.StateReady
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[models.StateError]; ok && running == 0 {
			d.Status.State = models.StateError
			d.Status.Message = message
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[models.StateProvision]; ok {
			d.Status.State = models.StateProvision
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		d.Status.State = models.StateDegradation
		d.Status.Message = models.EmptyString
		d.Meta.Updated = time.Now()
		break
	case models.StateReady:

		if _, ok := state[models.StateReady]; ok && running == len(pl) {
			break
		}

		if _, ok := state[models.StateProvision]; ok {
			d.Status.State = models.StateProvision
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[models.StateError]; ok && running == 0 {
			d.Status.State = models.StateError
			d.Status.Message = message
			d.Meta.Updated = time.Now()
			break
		}

		d.Status.State = models.StateDegradation
		d.Status.Message = models.EmptyString
		d.Meta.Updated = time.Now()

		break
	case models.StateError:

		if _, ok := state[models.StateReady]; ok && running == len(pl) {
			d.Status.State = models.StateReady
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[models.StateError]; ok && running == 0 {
			break
		}

		if _, ok := state[models.StateProvision]; ok {
			d.Status.State = models.StateProvision
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		d.Status.State = models.StateDegradation
		d.Status.Message = models.EmptyString
		d.Meta.Updated = time.Now()

		break
	case models.StateDegradation:

		if _, ok := state[models.StateReady]; ok && running == len(pl) {
			d.Status.State = models.StateReady
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[models.StateProvision]; ok {
			d.Status.State = models.StateProvision
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
			break
		}

		if _, ok := state[models.StateError]; ok && running == 0 {
			d.Status.State = models.StateError
			d.Status.Message = message
			d.Meta.Updated = time.Now()
			break
		}

		break
	case models.StateDestroy:
		if len(pl) == 0 {
			d.Status.State = models.StateDestroyed
			d.Status.Message = models.EmptyString
			d.Meta.Updated = time.Now()
		}
		break
	case models.StateDestroyed:
		break
	}

	log.Debugf("%s:> deploymentStatusState: finish: %s > %s", logDeploymentPrefix, d.SelfLink(), d.Status.State)

	return nil
}
