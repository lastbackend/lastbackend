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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"time"
)

const logPodPrefix = "state:observer:pod"

// PodObserve function manages pod handlers based on pod state
func PodObserve(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > state %s", logPodPrefix, p.SelfLink(), p.Status.State)
	// Call pod state manager methods
	switch p.Status.State {
	case types.StateCreated:
		if err := handlePodStateCreated(ss, p); err != nil {
			log.Errorf("%s:> handle pod state created err: %s", logPodPrefix, err.Error())
			return err
		}
	case types.StateProvision:
		if err := handlePodStateProvision(ss, p); err != nil {
			log.Errorf("%s:> handle pod state provision err: %s", logPodPrefix, err.Error())
			return err
		}
	case types.StateReady:
		if err := handlePodStateReady(ss, p); err != nil {
			log.Errorf("%s:> handle pod state ready err: %s", logPodPrefix, err.Error())
			return err
		}
	case types.StateError:
		if err := handlePodStateError(ss, p); err != nil {
			log.Errorf("%s:> handle pod state error err: %s", logPodPrefix, err.Error())
			return err
		}
		break
	case types.StateDegradation:
		if err := handlePodStateDegradation(ss, p); err != nil {
			log.Errorf("%s:> handle pod state degradation err: %s", logPodPrefix, err.Error())
			return err
		}
		break
	case types.StateDestroy:
		if err := handlePodStateDestroy(ss, p); err != nil {
			log.Errorf("%s:> handle pod state destroy err: %s", logPodPrefix, err.Error())
			return err
		}
		break
	case types.StateDestroyed:
		if err := handlePodStateDestroyed(ss, p); err != nil {
			log.Errorf("%s:> handle pod state destroyed err: %s", logPodPrefix, err.Error())
			return err
		}
		break

	}

	log.V(logLevel).Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink())

	pl, ok := ss.pod.list[p.DeploymentLink()]
	if ok && p.Status.State != types.StateDestroyed {
		pl[p.SelfLink()] = p
	}

	d, ok := ss.deployment.list[p.DeploymentLink()]
	if !ok {
		log.V(logLevel).Debugf("%s:> deployment node found: %s", logPodPrefix, p.DeploymentLink())
		return nil
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)
	return deploymentStatusState(d, pl)
}

func handlePodStateCreated(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateCreated: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(ss, p); err != nil {
		return err
	}
	log.V(logLevel).Debugf("%s handle pod create state finish: %s : %s", logPodPrefix, p.SelfLink(), p.Status.State)
	return nil
}

func handlePodStateProvision(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateProvision: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(ss, p); err != nil {
		return err
	}
	return nil
}

func handlePodStateReady(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateReady: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateError(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateError: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDegradation(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDegradation: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDestroy(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDestroy: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podDestroy(ss, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handlePodStateDestroyed(ss *ServiceState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDestroyed: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podRemove(ss, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

// podCreate function creates new pod based on deployment spec
func podCreate(d *types.Deployment) (*types.Pod, error) {
	dm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	return dm.Create(d)
}

// podDestroy function marks pod as provision
func podProvision(ss *ServiceState, p *types.Pod) (err error) {

	t := p.Meta.Updated

	defer func() {

		if err == nil {
			err = podUpdate(p, t)
		}

	}()

	if p.Meta.Node == types.EmptyString {

		var node *types.Node

		node, err = ss.cluster.PodLease(p)
		if err != nil {
			log.Errorf("%s:> pod node lease err: %s", logPrefix, err.Error())
			return err
		}

		if node == nil {
			p.Status.State = types.StateError
			p.Status.Message = errors.NodeNotFound
			p.Meta.Updated = time.Now()
			return nil
		}

		p.Meta.Node = node.SelfLink()
		p.Meta.Updated = time.Now()
	}

	if err = podManifestPut(p); err != nil {
		log.Errorf("%s:> pod manifest create err: %s", logPrefix, err.Error())
		return err
	}

	if p.Status.State != types.StateProvision {
		p.Status.State = types.StateProvision
		p.Meta.Updated = time.Now()
	}

	return nil
}

// podDestroy function marks pod spec as destroy
func podDestroy(ss *ServiceState, p *types.Pod) (err error) {

	t := p.Meta.Updated
	defer func() {
		if err == nil {
			err = podUpdate(p, t)
		}
	}()

	if p.Spec.State.Destroy {

		if p.Meta.Node == types.EmptyString {
			p.Status.State = types.StateDestroyed
			p.Meta.Updated = time.Now()
			return nil
		}

		if p.Status.State != types.StateDestroy {
			p.Status.State = types.StateDestroy
			p.Meta.Updated = time.Now()
		}

		return nil
	}

	p.Spec.State.Destroy = true
	if err = podManifestSet(p); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			if p.Meta.Node != types.EmptyString {
				if _, err := ss.cluster.PodRelease(p); err != nil {
					if !errors.Storage().IsErrEntityNotFound(err) {
						return err
					}
				}
			}

			p.Status.State = types.StateDestroyed
			p.Meta.Updated = time.Now()
			return nil
		}

		return err
	}

	p.Status.State = types.StateDestroy

	if p.Meta.Node == types.EmptyString {
		p.Status.State = types.StateDestroyed
	}

	p.Meta.Updated = time.Now()
	return nil
}

// podRemove function removes pod from storage if node is released
func podRemove(ss *ServiceState, p *types.Pod) (err error) {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	if _, err = ss.cluster.PodRelease(p); err != nil {
		return err
	}

	if err = podManifestDel(p); err != nil {
		return err
	}

	p.Meta.Node = types.EmptyString
	p.Meta.Updated = time.Now()

	if err = pm.Remove(p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	ss.DelPod(p)
	return nil
}

func podUpdate(p *types.Pod, timestamp time.Time) error {

	if timestamp.Before(p.Meta.Updated) {
		pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
		if err := pm.Update(p); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestPut(p *types.Pod) error {

	mm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	m, err := mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	if m == nil {
		pm := types.PodManifest(p.Spec)

		if err := mm.ManifestAdd(p.Meta.Node, p.Meta.SelfLink, &pm); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestSet(p *types.Pod) error {

	var (
		m   *types.PodManifest
		err error
	)

	mm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	m, err = mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink)
	if err != nil {
		return err
	}

	// Update manifest
	if m == nil {
		ms := types.PodManifest(p.Spec)
		m = &ms
	} else {
		*m = types.PodManifest(p.Spec)
	}

	if err := mm.ManifestSet(p.Meta.Node, p.Meta.SelfLink, m); err != nil {
		return err
	}

	return nil
}

func podManifestDel(p *types.Pod) error {

	if p.Meta.Node == types.EmptyString {
		return nil
	}

	// Remove manifest
	mm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	err := mm.ManifestDel(p.Meta.Node, p.SelfLink())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	return nil
}
