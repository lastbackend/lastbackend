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
	"sync"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logPrefix = "state:service"

type ServiceState struct {
	lock sync.Mutex

	service  *types.Service
	endpoint *types.Endpoint

	deployment struct {
		active    *types.Deployment
		provision *types.Deployment
		list      map[string]*types.Deployment
	}
	pod struct {
		list map[string]map[string]*types.Pod
	}

	observers struct {
		service    chan *types.Service
		deployment chan *types.Deployment
		pod        chan *types.Pod
	}
}

func (ss *ServiceState) Restore() error {

	log.Debugf("%s:restore state for service: %s", logPrefix, ss.service.SelfLink())

	var (
		err error
		stg = envs.Get().GetStorage()
	)

	// Get all pods
	pm := distribution.NewPodModel(context.Background(), stg)
	pl, err := pm.ListByService(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get pod map error: %v", logPrefix, err)
		return err
	}

	for _, p := range pl.Items {
		log.Infof("%s: restore: restore pod: %s", logPrefix, p.SelfLink())

		// Check if deployment map for pod exists
		if _, ok := ss.pod.list[p.DeploymentLink()]; !ok {
			ss.pod.list[p.DeploymentLink()] = make(map[string]*types.Pod)
		}

		// put pod into map by deployment name and pod name
		ss.pod.list[p.DeploymentLink()][p.SelfLink()] = p
	}

	// Get all deployments
	dm := distribution.NewDeploymentModel(context.Background(), stg)
	dl, err := dm.ListByService(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get deployment map error: %v", logPrefix, err)
		return err
	}

	for _, d := range dl.Items {
		log.Infof("%s: restore deployment: %s", logPrefix, d.SelfLink())
		ss.deployment.list[d.SelfLink()] = d
	}

	// Set service current spec and provision spec
	switch ss.service.Status.State {
	// if service is in ready state - mark deployment with same spec as current
	case types.StateReady:
		for _, d := range ss.deployment.list {
			if d.Spec.Template.Updated.Equal(ss.service.Spec.Template.Updated) {
				ss.deployment.active = d
			}
		}
		break
	// if service is in provision state - mark deployment in ready state as current
	case types.StateProvision:
		for _, d := range ss.deployment.list {
			if d.Status.State == types.StateReady {
				ss.deployment.active = d
			}

			if ss.deployment.provision == nil {
				ss.deployment.provision = d
				continue
			}
			// Mark latest created deployment as provision deployment for current service
			if ss.deployment.provision.Spec.Template.Updated.Before(d.Spec.Template.Updated) {
				ss.deployment.provision = d
			}
		}
		break
	}

	// Get endpoint
	em := distribution.NewEndpointModel(context.Background(), stg)
	ss.endpoint, err = em.Get(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:restore:> get endpoint error: %v", logPrefix, err)
			return err
		}
	}

	// Range over pods to sync pod status
	for _, pl := range ss.pod.list {
		for _, p := range pl {
			ss.observers.pod <- p
		}
	}

	// Provision deployment only in provision state
	for _, d := range ss.deployment.list {
		if d.Status.State == types.StateProvision {
			ss.observers.deployment <- d
		}
	}

	ss.observers.service <- ss.service
	return nil
}

func (ss *ServiceState) Observe() {
	for {
		select {

		case p := <-ss.observers.pod:
			log.Debugf("%s:observe:pod:> %v", logPrefix, p)

			// Call pod state manager methods
			switch p.Status.State {
			case types.StateProvision:
			case types.StateCreated:
			case types.StateReady:
			case types.StateError:
				break
			case types.StateDestroyed:
				if err := PodRemove(p); err != nil {
					log.Errorf("%s", err.Error())
					break
				}
				delete(ss.pod.list[p.DeploymentLink()], p.SelfLink())
				break
			case types.StateDestroy:
				if err := PodDestroy(p); err != nil {
					log.Errorf("%s", err.Error())
					break
				}
				break
			}

			// Sync deployment state after pod changes
			DeploymentSync()

			break

		case d := <-ss.observers.deployment:
			log.Debugf("%s:observe:deployment:> %v", logPrefix, d)

			switch d.Status.State {
			case types.StateCreated:
			case types.StateProvision:


				link := d.SelfLink()
				log.Debugf("0:> %s:> %#v", link, ss.pod.list)

				_, ok := ss.pod.list[link]
				if !ok {
					log.Debugf("1:> %#v", ss.pod.list[link])
					ss.pod.list[link] = make(map[string]*types.Pod)
				}

				log.Debugf("2:> %#v", ss.pod.list[link])
				if err := DeploymentProvision(d, ss.pod.list[link]); err != nil {
					log.Errorf("%s", err.Error())
					break
				}
				break
			case types.StateReady:
				if ss.deployment.active != nil {
					if err := DeploymentCancel(ss.deployment.active, ss.pod.list[ss.deployment.active.SelfLink()]); err != nil {
						log.Errorf("%s", err.Error())
						break
					}
				}

				ss.deployment.active = d
				ss.deployment.provision = nil

				break
			case types.StateError:
				break
			case types.StateDestroy:
				if err := DeploymentDestroy(d, ss.pod.list[d.SelfLink()]); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			case types.StateDestroyed:
				if err := DeploymentRemove(d); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			}

			ss.deployment.list[d.SelfLink()] = d

			if ss.deployment.active != nil {
				if err := ServiceSync(ss.service, ss.deployment.active); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			}

			if ss.deployment.provision != nil {
				if err := ServiceSync(ss.service, ss.deployment.provision); err != nil {
					log.Errorf("%s", err.Error())
				}
				break
			}

			if err := ServiceSync(ss.service, nil); err != nil {
				log.Errorf("%s", err.Error())
			}

			break

		case s := <-ss.observers.service:
			log.Debugf("%s:observe:service:> %v", logPrefix, s)

			switch s.Status.State {
			case types.StateCreated:

				if len(s.Spec.Network.Ports) > 0 && ss.endpoint == nil {

					e, err := EndpointCreate(s.Meta.Namespace, s.Meta.Name, s.Meta.Endpoint, &s.Spec.Network)
					if err != nil {
						log.Errorf("%s", err.Error())
						break
					}

					ss.endpoint = e
				}

				d, err := ServiceProvision(s)
				if err != nil {
					log.Errorf("%s", err.Error())
					continue
				}

				ss.deployment.list[d.SelfLink()] = d
				ss.deployment.provision = d

				if ss.endpoint != nil {
					e, err := EndpointUpdate(ss.endpoint, &s.Spec.Network)
					if err != nil {
						log.Errorf("%s", err.Error())
						break
					}
					ss.endpoint = e
				}

				break

			case types.StateProvision:

				// service template spec updated time check
				if (ss.deployment.active == nil && ss.deployment.provision == nil) || ss.service.Spec.Template.Updated.Before(s.Spec.Template.Updated) {

					d, err := ServiceProvision(s)
					if err != nil {
						log.Errorf("%s", err.Error())
						continue
					}

					ss.deployment.list[d.SelfLink()] = d
					ss.deployment.provision = d
				}

				// check expose spec and update endpoint if spec changed
				if ss.service.Spec.Network.Updated.Before(s.Spec.Network.Updated) {

					if len(s.Spec.Network.Ports) == 0 {
						if err := EndpointRemove(ss.endpoint); err != nil {
							log.Errorf("%s", err.Error())
						}
						ss.endpoint = nil
					} else {
						e, err := EndpointUpdate(ss.endpoint, &s.Spec.Network)
						if err != nil {
							log.Errorf("%s", err.Error())
						}
						ss.endpoint = e
					}
				}

				// check replicas changes
				if ss.service.Spec.Replicas != s.Spec.Replicas {

					var (
						err error
					)

					err = DeploymentScale(ss.deployment.active, s.Spec.Replicas)
					if err != nil {
						log.Errorf("%s", err.Error())
					}
				}

				break

			// Check service ready/error state triggers
			case types.StateReady:
			case types.StateError:
				// Generate events if needed
				break

			// Run service destroy process
			case types.StateDestroy:
				ss.service = s
				if err := ServiceDestroy(ss.service, ss.deployment.list); err != nil {
					log.Errorf("%s")
				}
				break

			// Remove service from storage if it is already destroyed
			case types.StateDestroyed:
				if err := ServiceRemove(ss.service); err != nil {
					log.Errorf("%s")
				}
				break
			}

			ss.service = s
			break
		}
	}
}

func (ss *ServiceState) SetService(s *types.Service) {
	ss.observers.service <- s
}

func (ss *ServiceState) SetPod(p *types.Pod) {
	ss.observers.pod <- p
}

func (ss *ServiceState) DelPod(p *types.Pod) {
	if _, ok := ss.pod.list[p.DeploymentLink()]; !ok {
		return
	}

	delete(ss.pod.list[p.DeploymentLink()], p.SelfLink())

	if _, ok := ss.deployment.list[p.DeploymentLink()]; !ok {
		return
	}

	d := ss.deployment.list[p.DeploymentLink()]
	ss.observers.deployment <- d

}

func NewServiceState(s *types.Service) *ServiceState {

	var ss = new(ServiceState)

	ss.service = s

	ss.observers.service = make(chan *types.Service)
	ss.observers.deployment = make(chan *types.Deployment)
	ss.observers.pod = make(chan *types.Pod)

	ss.deployment.list = make(map[string]*types.Deployment)
	ss.pod.list = make(map[string]map[string]*types.Pod)

	go ss.Observe()

	return ss
}
