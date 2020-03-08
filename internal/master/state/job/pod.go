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

package job

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logPodPrefix = "state:observer:pod"

// PodObserve function manages pod handlers based on pod state
func PodObserve(js *JobState, p *types.Pod) (err error) {

	log.V(logLevel).Debugf("%s:> observe start: %s > state %s", logPodPrefix, p.SelfLink(), p.Status.State)

	// Call pod state manager methods
	switch p.Status.State {
	case types.StateCreated:
		err = handlePodStateCreated(js, p)
	case types.StateProvision:
		err = handlePodStateProvision(js, p)
	case types.StateReady:
		err = handlePodStateReady(js, p)
	case types.StateError:
		err = handlePodStateError(js, p)
	case types.StateDegradation:
		err = handlePodStateDegradation(js, p)
	case types.StateDestroy:
		err = handlePodStateDestroy(js, p)
	case types.StateDestroyed:
		err = handlePodStateDestroyed(js, p)
	}
	if err != nil {
		log.Errorf("%s:> handle pod state %s err: %s", logPodPrefix, p.Status.State, err.Error())
		return err
	}

	log.V(logLevel).Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink())

	_, sl := p.SelfLink().Parent()
	if p.Status.State == types.StateDestroyed {
		delete(js.pod.list, sl.String())
		return nil
	} else {
		js.pod.list[sl.String()] = p

		task, ok := js.task.list[sl.String()]
		if !ok {
			log.V(logLevel).Debugf("%s:> task not found: %s", logPodPrefix, sl.String())
			return nil
		}

		log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)
		if err := taskStatusState(js, task, p); err != nil {
			return err
		}
	}

	log.V(logLevel).Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink())

	return nil
}

func handlePodStateCreated(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateCreated: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(js, p); err != nil {
		return err
	}

	return nil
}

func handlePodStateProvision(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateProvision: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(js, p); err != nil {
		return err
	}

	return nil
}

func handlePodStateReady(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateReady: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateError(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateError: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDegradation(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDegradation: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDestroy(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDestroy: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podDestroy(js, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handlePodStateDestroyed(js *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDestroyed: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podRemove(js, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

// podCreate function creates new pod based on task spec
func podCreate(stg storage.Storage, t *types.Task) (*types.Pod, error) {
	pm := model.NewPodModel(context.Background(), stg)

	pod := types.NewPod()
	pod.Meta.SetDefault()
	pod.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	pod.Meta.Namespace = t.Meta.Namespace
	sl, _ := types.NewPodSelfLink(types.KindTask, t.SelfLink().String(), pod.Meta.Name)
	pod.Meta.SelfLink = *sl
	pod.Status.SetCreated()

	pod.Spec.SetSpecRuntime(t.Spec.Runtime)
	pod.Spec.SetSpecTemplate(pod.SelfLink().String(), t.Spec.Template)
	pod.Spec.Selector = t.Spec.Selector

	return pm.Put(pod)
}

// podDestroy function marks pod as provision
func podProvision(js *JobState, p *types.Pod) (err error) {

	t := p.Meta.Updated

	defer func() {
		if err == nil {
			err = podUpdate(js.storage, p, t)
		}
	}()

	if p.Status.State != types.StateProvision {
		p.Status.State = types.StateProvision
		p.Meta.Updated = time.Now()
	}

	if p.Meta.Node == types.EmptyString {

		var node *types.Node

		node, err = js.cluster.PodLease(p)
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

		p.Meta.Node = node.SelfLink().String()
		p.Meta.Updated = time.Now()
	}

	if err = podManifestPut(js.storage, p); err != nil {
		log.Errorf("%s:> pod manifest create err: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

// podDestroy function marks pod spec as destroy
func podDestroy(js *JobState, p *types.Pod) (err error) {

	t := p.Meta.Updated
	defer func() {
		if err == nil {
			err = podUpdate(js.storage, p, t)
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

	if err = podManifestSet(js.storage, p); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			if p.Meta.Node != types.EmptyString {
				if _, err := js.cluster.PodRelease(p); err != nil {
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
func podRemove(js *JobState, p *types.Pod) (err error) {

	pm := model.NewPodModel(context.Background(), js.storage)
	if _, err = js.cluster.PodRelease(p); err != nil {
		return err
	}

	if err = podManifestDel(js.storage, p); err != nil {
		return err
	}

	p.Meta.Node = types.EmptyString
	p.Meta.Updated = time.Now()

	if err = pm.Remove(p); err != nil && !errors.Storage().IsErrEntityNotFound(err) {
		log.Errorf("pod remove %s", err.Error())
		return err
	}

	js.DelPod(p)
	return nil
}

func podUpdate(stg storage.Storage, p *types.Pod, timestamp time.Time) error {

	if timestamp.Before(p.Meta.Updated) {
		pm := model.NewPodModel(context.Background(), stg)
		if err := pm.Update(p); err != nil {
			log.Errorf("pod update %s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestPut(stg storage.Storage, p *types.Pod) error {

	mm := model.NewPodModel(context.Background(), stg)
	m, err := mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	if m == nil {
		pm := types.PodManifest(p.Spec)

		if err := mm.ManifestAdd(p.Meta.Node, p.Meta.SelfLink.String(), &pm); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestSet(stg storage.Storage, p *types.Pod) error {

	var (
		m   *types.PodManifest
		err error
	)

	mm := model.NewPodModel(context.Background(), stg)
	m, err = mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	// Update manifest
	if m == nil {
		ms := types.PodManifest(p.Spec)
		m = &ms
	} else {
		*m = types.PodManifest(p.Spec)
	}

	if err := mm.ManifestSet(p.Meta.Node, p.Meta.SelfLink.String(), m); err != nil {
		return err
	}

	return nil
}

func podManifestDel(stg storage.Storage,  p *types.Pod) error {

	if p.Meta.Node == types.EmptyString {
		return nil
	}

	// Remove manifest
	mm := model.NewPodModel(context.Background(), stg)
	err := mm.ManifestDel(p.Meta.Node, p.SelfLink().String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	return nil
}
