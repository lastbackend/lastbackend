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
	"strconv"
	"strings"
	"sync"

	"github.com/lastbackend/lastbackend/internal/master/ipam"
	"github.com/lastbackend/lastbackend/internal/master/state/cluster"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 3
	logPrefix = "state:service"
)

type ServiceState struct {
	lock sync.Mutex

	storage storage.IStorage
	ipam    ipam.IPAM

	cluster  *cluster.ClusterState
	service  *models.Service
	endpoint struct {
		endpoint *models.Endpoint
		manifest *models.EndpointManifest
	}

	deployment struct {
		index     int
		active    *models.Deployment
		provision *models.Deployment
		list      map[string]*models.Deployment
	}
	pod struct {
		list map[string]map[string]*models.Pod
	}

	observers struct {
		service    chan *models.Service
		deployment chan *models.Deployment
		pod        chan *models.Pod
	}
}

func (ss *ServiceState) Namespace() string {
	return ss.service.Meta.Namespace
}

func (ss *ServiceState) Restore() error {

	log.Debugf("%s:restore state for service: %s", logPrefix, ss.service.SelfLink())

	var (
		err error
	)

	// Get all pods
	pm := service.NewPodModel(context.Background(), ss.storage)
	pl, err := pm.ListByService(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get pod map error: %v", logPrefix, err)
		return err
	}

	for _, p := range pl.Items {
		log.Infof("%s: restore: restore pod: %s", logPrefix, p.SelfLink().String())

		// Check if deployment map for pod exists
		_, sl := p.SelfLink().Parent()
		if _, ok := ss.pod.list[sl.String()]; !ok {
			ss.pod.list[sl.String()] = make(map[string]*models.Pod)
		}

		// put pod into map by deployment name and pod name
		ss.pod.list[sl.String()][p.SelfLink().String()] = p
	}

	// Get all deployments
	dm := service.NewDeploymentModel(context.Background(), ss.storage)
	dl, err := dm.ListByService(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get deployment map error: %v", logPrefix, err)
		return err
	}

	for _, d := range dl.Items {
		log.Infof("%s: restore deployment: %s", logPrefix, d.SelfLink())

		var index int

		index, err := strconv.Atoi(strings.Replace(d.Meta.Name, "v", "", -1))
		if err != nil {
			log.Errorf("%s:> get deployment index err: %s", logPrefix, err.Error())
		}

		if ss.deployment.index < index {
			ss.deployment.index = index
		}

		log.Infof("index:> %d", ss.deployment.index)

		ss.deployment.list[d.SelfLink().String()] = d
	}

	// Set service current spec and provision spec
	switch ss.service.Status.State {
	// if service is in ready state - mark deployment with same spec as current
	case models.StateReady:
		for _, d := range ss.deployment.list {
			if d.Spec.Template.Updated.Equal(ss.service.Spec.Template.Updated) {
				ss.deployment.active = d
			}
		}
		break
	// if service is in provision state - mark deployment in ready state as current
	case models.StateWaiting:
		fallthrough
	case models.StateDegradation:
		fallthrough
	case models.StateProvision:

		for _, d := range ss.deployment.list {
			if d.Status.State == models.StateReady {
				ss.deployment.active = d
				continue
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

		if ss.deployment.active != nil && ss.deployment.provision != nil {
			if ss.deployment.provision.Spec.Template.Updated.Equal(ss.deployment.active.Spec.Template.Updated) {
				ss.deployment.provision = nil
			}
		}

		break
	}

	// Get endpoint
	if err := endpointRestore(ss); err != nil {
		log.Errorf("%s: restore endpoint: %s", logPrefix, err.Error())
		return err
	}

	// Range over pods to sync pod status
	for _, pl := range ss.pod.list {
		for _, p := range pl {
			ss.observers.pod <- p
		}
	}

	// Provision deployment only in provision state
	for _, d := range ss.deployment.list {
		ss.observers.deployment <- d
	}

	ss.observers.service <- ss.service

	return nil
}

func (ss *ServiceState) Observe() {
	for {
		select {

		case p := <-ss.observers.pod:
			log.Debugf("%s:observe:pod:start> %s", logPrefix, p.SelfLink())
			if err := PodObserve(ss, p); err != nil {
				log.Errorf("%s:observe:pod err:> %s", logPrefix, err.Error())
			}
			log.Debugf("%s:observe:pod:finish> %s", logPrefix, p.SelfLink())
			break

		case d := <-ss.observers.deployment:
			log.Debugf("%s:observe:deployment:start> %s", logPrefix, d.SelfLink())
			if err := deploymentObserve(ss, d); err != nil {
				log.Errorf("%s:observe:deployment err:> %s", logPrefix, err.Error())
			}
			log.Debugf("%s:observe:deployment:finish> %s", logPrefix, d.SelfLink())
			break

		case s := <-ss.observers.service:
			log.Debugf("%s:observe:service:start> %s", logPrefix, s.SelfLink())
			if err := serviceObserve(ss, s); err != nil {
				log.Errorf("%s:observe:service err:> %s", logPrefix, err.Error())
			}
			log.Debugf("%s:observe:service:finish> %s", logPrefix, s.SelfLink())
			break
		}

	}
}

func (ss *ServiceState) SetService(s *models.Service) {
	ss.observers.service <- s
}

func (ss *ServiceState) SetDeployment(d *models.Deployment) {
	ss.observers.deployment <- d
}

func (ss *ServiceState) DelDeployment(d *models.Deployment) {

	if _, ok := ss.pod.list[d.SelfLink().String()]; !ok {
		return
	}
	delete(ss.pod.list, d.SelfLink().String())

	if _, ok := ss.deployment.list[d.SelfLink().String()]; !ok {
		return
	}

	delete(ss.deployment.list, d.SelfLink().String())

	if ss.deployment.active != nil {
		if ss.deployment.active.SelfLink() == d.SelfLink() {
			ss.deployment.active = nil
		}
	}

	if ss.deployment.provision != nil {
		if ss.deployment.provision.SelfLink() == d.SelfLink() {
			ss.deployment.provision = nil
		}
	}

}

func (ss *ServiceState) SetPod(p *models.Pod) {
	ss.observers.pod <- p
}

func (ss *ServiceState) DelPod(p *models.Pod) {
	_, sl := p.SelfLink().Parent()

	if sl == nil {
		return
	}

	if _, ok := ss.pod.list[sl.String()]; !ok {
		return
	}

	delete(ss.pod.list[sl.String()], p.SelfLink().String())
}

func (ss *ServiceState) CheckDeps(dep models.StatusDependency) {

	log.Debugf("%s:> check dependency: %s", logPrefix, dep.Name)

	if ss.deployment.provision == nil {
		log.Debugf("%s:> check dependency: %s: provision deployment not found", logPrefix, dep.Name)
		return
	}

	if ss.deployment.provision.Status.State == models.StateWaiting {

		switch dep.Type {
		case models.KindVolume:
			if _, ok := ss.deployment.provision.Status.Dependencies.Volumes[dep.Name]; !ok {
				return
			}

			ss.deployment.provision.Status.Dependencies.Volumes[dep.Name] = dep
			if ss.deployment.provision.Status.CheckDeps() {
				ss.deployment.provision.Status.State = models.StateCreated
				ss.observers.deployment <- ss.deployment.provision
			}
		case models.KindSecret:
			if _, ok := ss.deployment.provision.Status.Dependencies.Secrets[dep.Name]; !ok {
				return
			}

			ss.deployment.provision.Status.Dependencies.Secrets[dep.Name] = dep
			if ss.deployment.provision.Status.CheckDeps() {
				ss.deployment.provision.Status.State = models.StateCreated
				ss.observers.deployment <- ss.deployment.provision
			}

		case models.KindConfig:
			if _, ok := ss.deployment.provision.Status.Dependencies.Configs[dep.Name]; !ok {
				return
			}

			ss.deployment.provision.Status.Dependencies.Configs[dep.Name] = dep
			if ss.deployment.provision.Status.CheckDeps() {
				ss.deployment.provision.Status.State = models.StateCreated
				ss.observers.deployment <- ss.deployment.provision
			}
		}

	}
}

func NewServiceState(stg storage.IStorage, ipam ipam.IPAM, cs *cluster.ClusterState, s *models.Service) *ServiceState {

	var ss = new(ServiceState)

	ss.storage = stg
	ss.ipam = ipam
	ss.service = s
	ss.cluster = cs

	ss.observers.service = make(chan *models.Service)
	ss.observers.deployment = make(chan *models.Deployment)
	ss.observers.pod = make(chan *models.Pod)

	ss.deployment.list = make(map[string]*models.Deployment)
	ss.pod.list = make(map[string]map[string]*models.Pod)

	go ss.Observe()

	return ss
}
