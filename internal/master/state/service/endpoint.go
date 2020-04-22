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

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logEndpointPrefix = "state:observer:endpoint"

// endpointEqual - validate endpoint spec
// return true if spec is valid
func endpointEqual(e *models.Endpoint, svc *models.Service) bool {

	if e == nil {
		return false
	}

	if len(e.Spec.PortMap) != len(svc.Spec.Network.Ports) {
		return false
	}

	for p, s := range svc.Spec.Network.Ports {
		if _, ok := e.Spec.PortMap[p]; !ok {
			return false
		}

		if s != e.Spec.PortMap[p] {
			return false
		}

	}

	if e.Spec.Policy != svc.Spec.Network.Policy {
		return false
	}

	return true
}

func endpointRestore(ss *ServiceState) error {

	var (
		err error
	)

	em := service.NewEndpointModel(context.Background(), ss.storage)
	ss.endpoint.endpoint, err = em.Get(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:restore:> get endpoint error: %v", logPrefix, err)
			return err
		}
	}
	if ss.endpoint.endpoint == nil {
		return nil
	}

	ss.endpoint.manifest, err = em.ManifestGet(ss.endpoint.endpoint.SelfLink().String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:restore:> get endpoint error: %v", logPrefix, err)
			return err
		}
	}

	return nil
}

func endpointProvision(ss *ServiceState, svc *models.Service) error {

	if len(svc.Spec.Network.Ports) == 0 {

		if ss.endpoint.endpoint != nil {
			if err := endpointDel(ss); err != nil {
				return err
			}
		}

		if ss.endpoint.manifest != nil {
			if err := endpointManifestDel(ss); err != nil {
				return err
			}
		}

		return nil
	}

	if ss.endpoint.endpoint != nil {
		if !endpointEqual(ss.endpoint.endpoint, svc) {
			if err := endpointSet(ss, svc); err != nil {
				return err
			}
		}
	}

	if ss.endpoint.endpoint == nil {
		if err := endpointAdd(ss, svc); err != nil {
			return err
		}
	}

	return nil
}

func endpointAdd(ss *ServiceState, svc *models.Service) error {

	var (
		err error
		em  = service.NewEndpointModel(context.Background(), ss.storage)
	)

	if svc.Spec.Network.IP == models.EmptyString {
		ip, err := ss.ipam.Lease()
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		svc.Spec.Network.IP = ip.String()
	}

	opts := models.EndpointCreateOptions{
		IP:            svc.Spec.Network.IP,
		Ports:         svc.Spec.Network.Ports,
		Policy:        svc.Spec.Network.Policy,
		BindStrategy:  svc.Spec.Network.Strategy.Bind,
		RouteStrategy: svc.Spec.Network.Strategy.Route,
		Domain:        svc.Meta.Endpoint,
	}

	ss.endpoint.endpoint, err = em.Create(svc.Meta.Namespace, svc.Meta.Name, &opts)
	if err != nil {
		log.Errorf("%s> create endpoint error: %s", logPrefix, err.Error())
		return err
	}

	if ss.endpoint.endpoint == nil {
		log.Debugf("%s> create endpoint error: %s", logPrefix, errors.New("endpoint is nil"))
		return nil
	}
	return nil

}

func endpointSet(ss *ServiceState, svc *models.Service) error {

	var (
		err error
		em  = service.NewEndpointModel(context.Background(), ss.storage)
	)

	opts := models.EndpointUpdateOptions{
		Ports:         svc.Spec.Network.Ports,
		Policy:        svc.Spec.Network.Policy,
		BindStrategy:  svc.Spec.Network.Strategy.Bind,
		RouteStrategy: svc.Spec.Network.Strategy.Route,
	}

	ss.endpoint.endpoint, err = em.Update(ss.endpoint.endpoint, &opts)
	if err != nil {
		log.Errorf("%s> set endpoint error: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

func endpointDel(ss *ServiceState) error {

	em := service.NewEndpointModel(context.Background(), ss.storage)
	if ss.endpoint.endpoint != nil {
		if err := em.Remove(ss.endpoint.endpoint); err != nil {
			log.Errorf("%s> del endpoint error: %s", logEndpointPrefix, err.Error())
			return err
		}
	}

	ss.endpoint.endpoint = nil

	return nil
}

func endpointCheck(ss *ServiceState) error {

	if ss.deployment.active != nil {
		if ss.deployment.active.Status.State == models.StateReady {
			if err := endpointManifestProvision(ss); err != nil {
				return err
			}
		}
	}

	return nil
}

func endpointManifestSpecEqual(e *models.Endpoint, m *models.EndpointManifest) bool {

	if e.Spec.IP != m.IP {
		return false
	}

	if e.Spec.Domain != m.Domain {
		return false
	}

	if e.Spec.Policy != m.Policy {
		return false
	}

	if e.Spec.Strategy.Route != m.Strategy.Route {
		return false
	}

	if e.Spec.Strategy.Bind != m.Strategy.Bind {
		return false
	}

	for p, mp := range e.Spec.PortMap {
		if _, ok := m.PortMap[p]; !ok {
			return false
		}

		if mp != m.PortMap[p] {
			return false
		}
	}

	if len(e.Spec.Upstreams) != len(m.Upstreams) {
		return false
	}

	for _, u := range e.Spec.Upstreams {
		var f = false
		for _, mu := range m.Upstreams {
			if u == mu {
				f = true
				break
			}
		}
		if !f {
			return false
		}
	}

	return true
}

func endpointManifestUpstreamsEqual(m *models.EndpointManifest, pl map[string]*models.Pod) bool {

	var (
		ips = endpointManifestGetUpstreams(pl)
		ups = make(map[string]bool)
	)

	if len(m.Upstreams) != len(ips) {
		return false
	}

	for _, ip := range m.Upstreams {
		ups[ip] = true
	}

	for _, ip := range ips {
		if _, ok := ups[ip]; !ok {
			return false
		}
	}

	return true
}

func endpointManifestProvision(ss *ServiceState) error {

	if ss.endpoint.endpoint == nil {
		return nil
	}

	if ss.endpoint.manifest != nil {

		var pl = make(map[string]*models.Pod)

		if ss.deployment.active != nil {
			if _, ok := ss.pod.list[ss.deployment.active.SelfLink().String()]; ok {
				pl = ss.pod.list[ss.deployment.active.SelfLink().String()]
			}
		}

		if !endpointManifestSpecEqual(ss.endpoint.endpoint, ss.endpoint.manifest) || !endpointManifestUpstreamsEqual(ss.endpoint.manifest, pl) {
			if err := endpointManifestSet(ss); err != nil {
				return err
			}
		}

	}

	if ss.endpoint.manifest == nil {
		if err := endpointManifestAdd(ss); err != nil {
			return err
		}
	}

	return nil
}

func endpointManifestAdd(ss *ServiceState) error {

	var (
		err error
		em  = service.NewEndpointModel(context.Background(), ss.storage)
		pl  = make(map[string]*models.Pod)
	)

	if ss.endpoint.endpoint == nil {
		log.Debugf("%s> create endpoint manifest error: %s", logPrefix, errors.New("endpoint is nil"))
		return nil
	}

	if ss.deployment.active != nil {
		if _, ok := ss.pod.list[ss.deployment.active.SelfLink().String()]; ok {
			pl = ss.pod.list[ss.deployment.active.SelfLink().String()]
		}
	}

	epm, err := em.ManifestGet(ss.endpoint.endpoint.SelfLink().String())
	if err != nil {
		return err
	}

	if epm == nil {
		ss.endpoint.manifest = &models.EndpointManifest{}
		ss.endpoint.manifest.EndpointSpec = ss.endpoint.endpoint.Spec
		ss.endpoint.manifest.Upstreams = endpointManifestGetUpstreams(pl)

		if err = em.ManifestAdd(ss.endpoint.endpoint.SelfLink().String(), ss.endpoint.manifest); err != nil {
			log.Errorf("%s> add endpoint manifest error: %s", logPrefix, err.Error())
			return err
		}

		return nil
	}

	if endpointManifestSpecEqual(ss.endpoint.endpoint, epm) {
		return nil
	}

	epm.EndpointSpec = ss.endpoint.endpoint.Spec
	epm.Upstreams = endpointManifestGetUpstreams(pl)

	if err = em.ManifestSet(ss.endpoint.endpoint.SelfLink().String(), epm); err != nil {
		log.Errorf("%s> update endpoint manifest error: %s", logPrefix, err.Error())
		return err
	}

	ss.endpoint.manifest = epm
	return nil
}

func endpointManifestSet(ss *ServiceState) error {

	var (
		err error
		em  = service.NewEndpointModel(context.Background(), ss.storage)
		pl  = make(map[string]*models.Pod)
	)

	if ss.endpoint.endpoint == nil {
		log.Debugf("%s> update endpoint manifest error: %s", logPrefix, errors.New("endpoint is nil"))
		return nil
	}

	if ss.deployment.active != nil {
		if _, ok := ss.pod.list[ss.deployment.active.SelfLink().String()]; ok {
			pl = ss.pod.list[ss.deployment.active.SelfLink().String()]
		}
	}

	ss.endpoint.manifest.EndpointSpec = ss.endpoint.endpoint.Spec
	ss.endpoint.manifest.Upstreams = endpointManifestGetUpstreams(pl)

	if err = em.ManifestSet(ss.endpoint.endpoint.SelfLink().String(), ss.endpoint.manifest); err != nil {
		log.Errorf("%s> update endpoint manifest error: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

func endpointManifestDel(ss *ServiceState) error {

	em := service.NewEndpointModel(context.Background(), ss.storage)

	if ss.endpoint.manifest != nil {
		if err := em.ManifestDel(em.ManifestGetSelfLink(ss.service.Meta.Namespace, ss.service.Meta.Name)); err != nil {
			log.Errorf("%s> del endpoint manifest error: %s", logEndpointPrefix, err.Error())
			return err
		}
	}

	ss.endpoint.manifest = nil

	return nil
}

func endpointManifestGetUpstreams(pl map[string]*models.Pod) []string {

	ips := make([]string, 0)

	for _, p := range pl {
		if p.Status.State == models.StateReady && p.Status.Network.PodIP != models.EmptyString {
			ips = append(ips, p.Status.Network.PodIP)
		}
	}

	return ips
}
