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

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logEndpointPrefix = "state:observer:endpoint"

// endpointEqual - validate endpoint spec
// return true if spec is valid
func endpointEqual(e *types.Endpoint, svc *types.Service) bool {

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
		stg = envs.Get().GetStorage()
	)

	em := distribution.NewEndpointModel(context.Background(), stg)
	ss.endpoint.endpoint, err = em.Get(ss.service.Meta.Namespace, ss.service.Meta.Name)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:restore:> get endpoint error: %v", logPrefix, err)
			return err
		}
	}

	ss.endpoint.manifest, err = em.ManifestGet(em.ManifestGetName(ss.service.Meta.Namespace, ss.service.Meta.Name))
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s:restore:> get endpoint error: %v", logPrefix, err)
			return err
		}
	}

	return nil
}

func endpointProvision(ss *ServiceState, svc *types.Service) error {

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

func endpointAdd(ss *ServiceState, svc *types.Service) error {

	var (
		err error
		em  = distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())
	)

	if svc.Spec.Network.IP == types.EmptyString {
		ip, err := envs.Get().GetIPAM().Lease()
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		svc.Spec.Network.IP = ip.String()
	}

	opts := types.EndpointCreateOptions{
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
		log.V(logLevel).Debugf("%s> create endpoint error: %s", logPrefix, errors.New("endpoint is nil"))
		return nil
	}
	return nil

}

func endpointSet(ss *ServiceState, svc *types.Service) error {

	var (
		err error
		em  = distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())
	)

	opts := types.EndpointUpdateOptions{
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

	em := distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())
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
		if ss.deployment.active.Status.State == types.StateReady {
			if err := endpointManifestProvision(ss); err != nil {
				return err
			}
		}
	}

	return nil
}

func endpointManifestSpecEqual(e *types.Endpoint, m *types.EndpointManifest) bool {

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

func endpointManifestUpstreamsEqual(m *types.EndpointManifest, pl map[string]*types.Pod) bool {

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

		var pl = make(map[string]*types.Pod)

		if ss.deployment.active != nil {
			if _, ok := ss.pod.list[ss.deployment.active.SelfLink()]; ok {
				pl = ss.pod.list[ss.deployment.active.SelfLink()]
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
		em  = distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())
		pl  = make(map[string]*types.Pod)
	)

	if ss.endpoint.endpoint == nil {
		log.V(logLevel).Debugf("%s> create endpoint manifest error: %s", logPrefix, errors.New("endpoint is nil"))
		return nil
	}

	if ss.deployment.active != nil {
		if _, ok := ss.pod.list[ss.deployment.active.SelfLink()]; ok {
			pl = ss.pod.list[ss.deployment.active.SelfLink()]
		}
	}

	epm, err := em.ManifestGet(ss.endpoint.endpoint.SelfLink())
	if err != nil {
		return err
	}

	if epm == nil {
		ss.endpoint.manifest = &types.EndpointManifest{}
		ss.endpoint.manifest.EndpointSpec = ss.endpoint.endpoint.Spec
		ss.endpoint.manifest.Upstreams = endpointManifestGetUpstreams(pl)

		if err = em.ManifestAdd(ss.endpoint.endpoint.SelfLink(), ss.endpoint.manifest); err != nil {
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

	if err = em.ManifestSet(ss.endpoint.endpoint.SelfLink(), epm); err != nil {
		log.Errorf("%s> update endpoint manifest error: %s", logPrefix, err.Error())
		return err
	}

	ss.endpoint.manifest = epm
	return nil
}

func endpointManifestSet(ss *ServiceState) error {

	var (
		err error
		em  = distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())
		pl  = make(map[string]*types.Pod)
	)

	if ss.endpoint.endpoint == nil {
		log.V(logLevel).Debugf("%s> update endpoint manifest error: %s", logPrefix, errors.New("endpoint is nil"))
		return nil
	}

	if ss.deployment.active != nil {
		if _, ok := ss.pod.list[ss.deployment.active.SelfLink()]; ok {
			pl = ss.pod.list[ss.deployment.active.SelfLink()]
		}
	}

	ss.endpoint.manifest.EndpointSpec = ss.endpoint.endpoint.Spec
	ss.endpoint.manifest.Upstreams = endpointManifestGetUpstreams(pl)

	if err = em.ManifestSet(ss.endpoint.endpoint.SelfLink(), ss.endpoint.manifest); err != nil {
		log.Errorf("%s> update endpoint manifest error: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

func endpointManifestDel(ss *ServiceState) error {

	em := distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())

	if ss.endpoint.manifest != nil {
		if err := em.ManifestDel(em.ManifestGetName(ss.service.Meta.Namespace, ss.service.Meta.Name)); err != nil {
			log.Errorf("%s> del endpoint manifest error: %s", logEndpointPrefix, err.Error())
			return err
		}
	}

	ss.endpoint.manifest = nil

	return nil
}

func endpointManifestGetUpstreams(pl map[string]*types.Pod) []string {

	ips := make([]string, 0)

	for _, p := range pl {
		if p.Status.State == types.StateReady && p.Status.Network.PodIP != types.EmptyString {
			ips = append(ips, p.Status.Network.PodIP)
		}
	}

	return ips
}
