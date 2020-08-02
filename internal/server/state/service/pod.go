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
//
//import (
//	"context"
//	"strings"
//	"time"
//
//	"github.com/lastbackend/lastbackend/internal/pkg/errors"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/internal/pkg/storage"
//	"github.com/lastbackend/lastbackend/internal/util/generator"
//	"github.com/lastbackend/lastbackend/tools/log"
//)
//
//const logPodPrefix = "state:observer:pod"
//
//// PodObserve function manages pod handlers based on pod state
//func PodObserve(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> observe start: %s > state %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//	// Call pod state manager methods
//	switch p.Status.State {
//	case models.StateCreated:
//		if err := handlePodStateCreated(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state created err: %s", logPodPrefix, err.Error())
//			return err
//		}
//	case models.StateProvision:
//		if err := handlePodStateProvision(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state provision err: %s", logPodPrefix, err.Error())
//			return err
//		}
//	case models.StateReady:
//		if err := handlePodStateReady(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state ready err: %s", logPodPrefix, err.Error())
//			return err
//		}
//	case models.StateError:
//		if err := handlePodStateError(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state error err: %s", logPodPrefix, err.Error())
//			return err
//		}
//		break
//	case models.StateDegradation:
//		if err := handlePodStateDegradation(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state degradation err: %s", logPodPrefix, err.Error())
//			return err
//		}
//		break
//	case models.StateDestroy:
//		if err := handlePodStateDestroy(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state destroy err: %s", logPodPrefix, err.Error())
//			return err
//		}
//		break
//	case models.StateDestroyed:
//		if err := handlePodStateDestroyed(ss, p); err != nil {
//			log.Errorf("%s:> handle pod state destroyed err: %s", logPodPrefix, err.Error())
//			return err
//		}
//		break
//
//	}
//
//	log.Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink().String())
//
//	_, sl := p.SelfLink().Parent()
//
//	if sl == nil {
//		return nil
//	}
//
//	pl, ok := ss.pod.list[sl.String()]
//	if ok && p.Status.State != models.StateDestroyed {
//		pl[p.SelfLink().String()] = p
//	}
//
//	d, ok := ss.deployment.list[sl.String()]
//	if !ok {
//		log.Debugf("%s:> deployment node found: %s", logPodPrefix, sl.String())
//		return nil
//	}
//
//	log.Debugf("%s:> observe finish: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	if err := deploymentStatusState(ss.storage, d, pl); err != nil {
//		return err
//	}
//
//	if ss.deployment.active != nil {
//		if ss.deployment.active.SelfLink().String() == d.SelfLink().String() && d.Status.State == models.StateReady {
//			if err := endpointCheck(ss); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
//
//func handlePodStateCreated(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateCreated: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	if err := podProvision(ss, p); err != nil {
//		return err
//	}
//	log.Debugf("%s handle pod create state finish: %s : %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//	return nil
//}
//
//func handlePodStateProvision(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateProvision: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	if err := podProvision(ss, p); err != nil {
//		return err
//	}
//	return nil
//}
//
//func handlePodStateReady(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateReady: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	return nil
//}
//
//func handlePodStateError(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateError: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	return nil
//}
//
//func handlePodStateDegradation(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateDegradation: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	return nil
//}
//
//func handlePodStateDestroy(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateDestroy: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	if err := podDestroy(ss, p); err != nil {
//		log.Errorf("%s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func handlePodStateDestroyed(ss *ServiceState, p *models.Pod) error {
//
//	log.Debugf("%s:> handlePodStateDestroyed: %s > %s", logPodPrefix, p.SelfLink().String(), p.Status.State)
//
//	if err := podRemove(ss, p); err != nil {
//		log.Errorf("%s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//// podCreate function creates new pod based on deployment spec
//func podCreate(stg storage.IStorage, d *models.Deployment) (*models.Pod, error) {
//	dm := service.NewPodModel(context.Background(), stg)
//
//	pod := models.NewPod()
//	pod.Meta.SetDefault()
//	pod.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
//	pod.Meta.Namespace = d.Meta.Namespace
//	sl, _ := models.NewPodSelfLink(models.KindDeployment, d.SelfLink().String(), pod.Meta.Name)
//	pod.Meta.SelfLink = *sl
//	pod.Status.SetCreated()
//
//	pod.Spec.SetSpecTemplate(pod.SelfLink().String(), d.Spec.Template)
//	pod.Spec.Selector = d.Spec.Selector
//
//	return dm.Put(pod)
//}
//
//// podDestroy function marks pod as provision
//func podProvision(ss *ServiceState, p *models.Pod) (err error) {
//
//	t := p.Meta.Updated
//
//	defer func() {
//
//		if err == nil {
//			err = podUpdate(ss.storage, p, t)
//		}
//
//	}()
//
//	if p.Meta.Node == models.EmptyString {
//
//		var node *models.Node
//
//		node, err = ss.cluster.PodLease(p)
//		if err != nil {
//			log.Errorf("%s:> pod node lease err: %s", logPrefix, err.Error())
//			return err
//		}
//
//		if node == nil {
//			p.Status.State = models.StateError
//			p.Status.Message = errors.NodeNotFound
//			p.Meta.Updated = time.Now()
//			return nil
//		}
//
//		p.Meta.Node = node.SelfLink().String()
//		p.Meta.Updated = time.Now()
//	}
//
//	if err = podManifestPut(ss.storage, p); err != nil {
//		log.Errorf("%s:> pod manifest create err: %s", logPrefix, err.Error())
//		return err
//	}
//
//	if p.Status.State != models.StateProvision {
//		p.Status.State = models.StateProvision
//		p.Meta.Updated = time.Now()
//	}
//
//	return nil
//}
//
//// podDestroy function marks pod spec as destroy
//func podDestroy(ss *ServiceState, p *models.Pod) (err error) {
//
//	t := p.Meta.Updated
//	defer func() {
//		if err == nil {
//			err = podUpdate(ss.storage, p, t)
//		}
//	}()
//
//	if p.Spec.State.Destroy {
//
//		if p.Meta.Node == models.EmptyString {
//			p.Status.State = models.StateDestroyed
//			p.Meta.Updated = time.Now()
//			return nil
//		}
//
//		if p.Status.State != models.StateDestroy {
//			p.Status.State = models.StateDestroy
//			p.Meta.Updated = time.Now()
//		}
//
//		return nil
//	}
//
//	p.Spec.State.Destroy = true
//	if err = podManifestSet(ss.storage, p); err != nil {
//		if errors.Storage().IsErrEntityNotFound(err) {
//			if p.Meta.Node != models.EmptyString {
//				if _, err := ss.cluster.PodRelease(p); err != nil {
//					if !errors.Storage().IsErrEntityNotFound(err) {
//						return err
//					}
//				}
//			}
//
//			p.Status.State = models.StateDestroyed
//			p.Meta.Updated = time.Now()
//			return nil
//		}
//
//		return err
//	}
//
//	p.Status.State = models.StateDestroy
//
//	if p.Meta.Node == models.EmptyString {
//		p.Status.State = models.StateDestroyed
//	}
//
//	p.Meta.Updated = time.Now()
//	return nil
//}
//
//// podRemove function removes pod from storage if node is released
//func podRemove(ss *ServiceState, p *models.Pod) (err error) {
//
//	pm := service.NewPodModel(context.Background(), ss.storage)
//	if _, err = ss.cluster.PodRelease(p); err != nil {
//		return err
//	}
//
//	if err = podManifestDel(ss.storage, p); err != nil {
//		return err
//	}
//
//	p.Meta.Node = models.EmptyString
//	p.Meta.Updated = time.Now()
//
//	if err = pm.Remove(p); err != nil {
//		log.Errorf("%s", err.Error())
//		return err
//	}
//
//	ss.DelPod(p)
//	return nil
//}
//
//func podUpdate(stg storage.IStorage, p *models.Pod, timestamp time.Time) error {
//
//	if timestamp.Before(p.Meta.Updated) {
//		pm := service.NewPodModel(context.Background(), stg)
//		if err := pm.Update(p); err != nil {
//			log.Errorf("%s", err.Error())
//			return err
//		}
//	}
//
//	return nil
//}
//
//func podManifestPut(stg storage.IStorage, p *models.Pod) error {
//
//	mm := service.NewPodModel(context.Background(), stg)
//	m, err := mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
//	if err != nil {
//		if !errors.Storage().IsErrEntityNotFound(err) {
//			log.Errorf("%s", err.Error())
//			return err
//		}
//	}
//
//	if m == nil {
//		pm := models.PodManifest(p.Spec)
//
//		if err := mm.ManifestAdd(p.Meta.Node, p.Meta.SelfLink.String(), &pm); err != nil {
//			log.Errorf("%s", err.Error())
//			return err
//		}
//	}
//
//	return nil
//}
//
//func podManifestSet(stg storage.IStorage, p *models.Pod) error {
//
//	var (
//		m   *models.PodManifest
//		err error
//	)
//
//	mm := service.NewPodModel(context.Background(), stg)
//	m, err = mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
//	if err != nil {
//		if !errors.Storage().IsErrEntityNotFound(err) {
//			log.Errorf("%s", err.Error())
//			return err
//		}
//	}
//
//	// Update manifest
//	if m == nil {
//		ms := models.PodManifest(p.Spec)
//		m = &ms
//	} else {
//		*m = models.PodManifest(p.Spec)
//	}
//
//	if err := mm.ManifestSet(p.Meta.Node, p.SelfLink().String(), m); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func podManifestDel(stg storage.IStorage, p *models.Pod) error {
//
//	if p.Meta.Node == models.EmptyString {
//		return nil
//	}
//
//	// Remove manifest
//	mm := service.NewPodModel(context.Background(), stg)
//	err := mm.ManifestDel(p.Meta.Node, p.SelfLink().String())
//	if err != nil {
//		if !errors.Storage().IsErrEntityNotFound(err) {
//			return err
//		}
//	}
//
//	return nil
//}
