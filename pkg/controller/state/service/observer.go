//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"strconv"
	"strings"
	"sync"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logLevel  = 3
	logPrefix = "state:service"
)

type ServiceState struct {
	lock sync.Mutex

	cluster  *cluster.ClusterState
	service  *types.Service
	endpoint struct {
		endpoint *types.Endpoint
		manifest *types.EndpointManifest
	}

	deployment struct {
		index     int
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

func (ss *ServiceState) Namespace() string {
	return ss.service.Meta.Namespace
}

func (ss *ServiceState) Restore() error {

	log.V(logLevel).Debugf("%s:restore state for service: %s", logPrefix, ss.service.SelfLink())

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

		var index int

		index, err := strconv.Atoi(strings.Replace(d.Meta.Name,"v","", -1))
		if err != nil {
			fmt.Println(err)
		}

		if ss.deployment.index < index {
			ss.deployment.index = index
		}

		log.Infof("index:> %d", ss.deployment.index)

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
	case types.StateWaiting:
		fallthrough
	case types.StateDegradation:
		fallthrough
	case types.StateProvision:

		for _, d := range ss.deployment.list {
			if d.Status.State == types.StateReady {
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
			log.V(logLevel).Debugf("%s:observe:pod:> %s", logPrefix, p.SelfLink())
			if err := PodObserve(ss, p); err != nil {
				log.Errorf("%s:observe:pod err:> %s", logPrefix, err.Error())
			}
			break

		case d := <-ss.observers.deployment:
			log.V(logLevel).Debugf("%s:observe:deployment:> %s", logPrefix, d.SelfLink())
			if err := deploymentObserve(ss, d); err != nil {
				log.Errorf("%s:observe:deployment err:> %s", logPrefix, err.Error())
			}
			break

		case s := <-ss.observers.service:
			log.V(logLevel).Debugf("%s:observe:service:> %s", logPrefix, s.SelfLink())
			if err := serviceObserve(ss, s); err != nil {
				log.Errorf("%s:observe:service err:> %s", logPrefix, err.Error())
			}
			break
		}

	}
}

func (ss *ServiceState) SetService(s *types.Service) {
	ss.observers.service <- s
}

func (ss *ServiceState) SetDeployment(d *types.Deployment) {
	ss.observers.deployment <- d
}

func (ss *ServiceState) DelDeployment(d *types.Deployment) {

	if _, ok := ss.pod.list[d.SelfLink()]; !ok {
		return
	}
	delete(ss.pod.list, d.SelfLink())

	if _, ok := ss.deployment.list[d.SelfLink()]; !ok {
		return
	}

	delete(ss.deployment.list, d.SelfLink())

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

func (ss *ServiceState) SetPod(p *types.Pod) {
	ss.observers.pod <- p
}

func (ss *ServiceState) DelPod(p *types.Pod) {
	if _, ok := ss.pod.list[p.DeploymentLink()]; !ok {
		return
	}

	delete(ss.pod.list[p.DeploymentLink()], p.SelfLink())
}

func (ss *ServiceState) CheckDeps(dep types.DeploymentStatusDependency) {

	log.Debugf("%s:> check dependency: %s",logPrefix, dep.Name)

	if ss.deployment.provision == nil {
		log.Debugf("%s:> check dependency: %s: provision deployment not found",logPrefix, dep.Name)
		return
	}

	if ss.deployment.provision.Status.State == types.StateWaiting {

		switch dep.Type {
		case types.KindVolume:
			if _, ok := ss.deployment.provision.Status.Dependencies.Volumes[dep.Name]; !ok {
				return
			}

			ss.deployment.provision.Status.Dependencies.Volumes[dep.Name] = dep
			if ss.deployment.provision.Status.CheckDeps() {
				ss.deployment.provision.Status.State = types.StateCreated
				ss.observers.deployment <- ss.deployment.provision
			}
		case types.KindSecret:
			if _, ok := ss.deployment.provision.Status.Dependencies.Secrets[dep.Name]; !ok {
				return
			}

			ss.deployment.provision.Status.Dependencies.Secrets[dep.Name] = dep
			if ss.deployment.provision.Status.CheckDeps() {
				ss.deployment.provision.Status.State = types.StateCreated
				ss.observers.deployment <- ss.deployment.provision
			}

		case types.KindConfig:
			if _, ok := ss.deployment.provision.Status.Dependencies.Configs[dep.Name]; !ok {
				return
			}

			ss.deployment.provision.Status.Dependencies.Configs[dep.Name] = dep
			if ss.deployment.provision.Status.CheckDeps() {
				ss.deployment.provision.Status.State = types.StateCreated
				ss.observers.deployment <- ss.deployment.provision
			}
		}

	}
}

func NewServiceState(cs *cluster.ClusterState, s *types.Service) *ServiceState {

	var ss = new(ServiceState)

	ss.service = s
	ss.cluster = cs

	ss.observers.service = make(chan *types.Service)
	ss.observers.deployment = make(chan *types.Deployment)
	ss.observers.pod = make(chan *types.Pod)

	ss.deployment.list = make(map[string]*types.Deployment)
	ss.pod.list = make(map[string]map[string]*types.Pod)

	go ss.Observe()

	return ss
}
