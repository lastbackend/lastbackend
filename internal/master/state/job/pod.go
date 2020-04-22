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
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logPodPrefix = "state:observer:pod"

// PodObserve function manages pod handlers based on pod state
func PodObserve(js *JobState, p *models.Pod) (err error) {

	log.Debugf("%s:> observe start: %s > state %s", logPodPrefix, p.SelfLink(), p.Status.State)

	// Call pod state manager methods
	switch p.Status.State {
	case models.StateCreated:
		err = handlePodStateCreated(js, p)
	case models.StateProvision:
		err = handlePodStateProvision(js, p)
	case models.StateReady:
		err = handlePodStateReady(js, p)
	case models.StateError:
		err = handlePodStateError(js, p)
	case models.StateDegradation:
		err = handlePodStateDegradation(js, p)
	case models.StateDestroy:
		err = handlePodStateDestroy(js, p)
	case models.StateDestroyed:
		err = handlePodStateDestroyed(js, p)
	}
	if err != nil {
		log.Errorf("%s:> handle pod state %s err: %s", logPodPrefix, p.Status.State, err.Error())
		return err
	}

	log.Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink())

	_, sl := p.SelfLink().Parent()
	if p.Status.State == models.StateDestroyed {
		delete(js.pod.list, sl.String())
		return nil
	} else {
		js.pod.list[sl.String()] = p

		task, ok := js.task.list[sl.String()]
		if !ok {
			log.Debugf("%s:> task not found: %s", logPodPrefix, sl.String())
			return nil
		}

		log.Debugf("%s:> observe finish: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)
		if err := taskStatusState(js, task, p); err != nil {
			return err
		}
	}

	log.Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink())

	return nil
}

func handlePodStateCreated(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateCreated: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(js, p); err != nil {
		return err
	}

	return nil
}

func handlePodStateProvision(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateProvision: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(js, p); err != nil {
		return err
	}

	return nil
}

func handlePodStateReady(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateReady: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateError(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateError: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDegradation(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateDegradation: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDestroy(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateDestroy: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podDestroy(js, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handlePodStateDestroyed(js *JobState, p *models.Pod) error {

	log.Debugf("%s:> handlePodStateDestroyed: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podRemove(js, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

// podCreate function creates new pod based on task spec
func podCreate(stg storage.IStorage, t *models.Task) (*models.Pod, error) {
	pm := service.NewPodModel(context.Background(), stg)

	pod := models.NewPod()
	pod.Meta.SetDefault()
	pod.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	pod.Meta.Namespace = t.Meta.Namespace
	sl, _ := models.NewPodSelfLink(models.KindTask, t.SelfLink().String(), pod.Meta.Name)
	pod.Meta.SelfLink = *sl
	pod.Status.SetCreated()

	pod.Spec.SetSpecRuntime(t.Spec.Runtime)
	pod.Spec.SetSpecTemplate(pod.SelfLink().String(), t.Spec.Template)
	pod.Spec.Selector = t.Spec.Selector

	return pm.Put(pod)
}

// podDestroy function marks pod as provision
func podProvision(js *JobState, p *models.Pod) (err error) {

	t := p.Meta.Updated

	defer func() {
		if err == nil {
			err = podUpdate(js.storage, p, t)
		}
	}()

	if p.Status.State != models.StateProvision {
		p.Status.State = models.StateProvision
		p.Meta.Updated = time.Now()
	}

	if p.Meta.Node == models.EmptyString {

		var node *models.Node

		node, err = js.cluster.PodLease(p)
		if err != nil {
			log.Errorf("%s:> pod node lease err: %s", logPrefix, err.Error())
			return err
		}

		if node == nil {
			p.Status.State = models.StateError
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
func podDestroy(js *JobState, p *models.Pod) (err error) {

	t := p.Meta.Updated
	defer func() {
		if err == nil {
			err = podUpdate(js.storage, p, t)
		}
	}()

	if p.Spec.State.Destroy {

		if p.Meta.Node == models.EmptyString {
			p.Status.State = models.StateDestroyed
			p.Meta.Updated = time.Now()
			return nil
		}

		if p.Status.State != models.StateDestroy {
			p.Status.State = models.StateDestroy
			p.Meta.Updated = time.Now()
		}
		return nil
	}

	p.Spec.State.Destroy = true

	if err = podManifestSet(js.storage, p); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			if p.Meta.Node != models.EmptyString {
				if _, err := js.cluster.PodRelease(p); err != nil {
					if !errors.Storage().IsErrEntityNotFound(err) {
						return err
					}
				}
			}

			p.Status.State = models.StateDestroyed
			p.Meta.Updated = time.Now()
			return nil
		}

		return err
	}

	p.Status.State = models.StateDestroy

	if p.Meta.Node == models.EmptyString {
		p.Status.State = models.StateDestroyed
	}

	p.Meta.Updated = time.Now()
	return nil
}

// podRemove function removes pod from storage if node is released
func podRemove(js *JobState, p *models.Pod) (err error) {

	pm := service.NewPodModel(context.Background(), js.storage)
	if _, err = js.cluster.PodRelease(p); err != nil {
		return err
	}

	if err = podManifestDel(js.storage, p); err != nil {
		return err
	}

	p.Meta.Node = models.EmptyString
	p.Meta.Updated = time.Now()

	if err = pm.Remove(p); err != nil && !errors.Storage().IsErrEntityNotFound(err) {
		log.Errorf("pod remove %s", err.Error())
		return err
	}

	js.DelPod(p)
	return nil
}

func podUpdate(stg storage.IStorage, p *models.Pod, timestamp time.Time) error {

	if timestamp.Before(p.Meta.Updated) {
		pm := service.NewPodModel(context.Background(), stg)
		if err := pm.Update(p); err != nil {
			log.Errorf("pod update %s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestPut(stg storage.IStorage, p *models.Pod) error {

	mm := service.NewPodModel(context.Background(), stg)
	m, err := mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	if m == nil {
		pm := models.PodManifest(p.Spec)

		if err := mm.ManifestAdd(p.Meta.Node, p.Meta.SelfLink.String(), &pm); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestSet(stg storage.IStorage, p *models.Pod) error {

	var (
		m   *models.PodManifest
		err error
	)

	mm := service.NewPodModel(context.Background(), stg)
	m, err = mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	// Update manifest
	if m == nil {
		ms := models.PodManifest(p.Spec)
		m = &ms
	} else {
		*m = models.PodManifest(p.Spec)
	}

	if err := mm.ManifestSet(p.Meta.Node, p.Meta.SelfLink.String(), m); err != nil {
		return err
	}

	return nil
}

func podManifestDel(stg storage.IStorage, p *models.Pod) error {

	if p.Meta.Node == models.EmptyString {
		return nil
	}

	// Remove manifest
	mm := service.NewPodModel(context.Background(), stg)
	err := mm.ManifestDel(p.Meta.Node, p.SelfLink().String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	return nil
}
