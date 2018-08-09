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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

// EndpointValidate - validate endpoint spec
// return true if spec is valid
func EndpointValidate(e *types.Endpoint, spec types.SpecNetwork) bool {

	if e == nil {
		return false
	}

	if len(e.Spec.PortMap) != len(spec.Ports) {
		return false
	}

	for p, s := range spec.Ports {
		if _, ok := e.Spec.PortMap[p]; !ok {
			return false
		}

		if s != e.Spec.PortMap[p] {
			return false
		}

	}

	if e.Spec.Policy != spec.Policy {
		return false
	}

	return true
}

func EndpointCreate(namespace, service, domain string, spec types.SpecNetwork) (*types.Endpoint, error) {
	em := distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())

	if spec.IP == types.EmptyString {
		ip, err := envs.Get().GetIPAM().Lease()
		if err != nil {
			log.Errorf("%s", err.Error())
			return nil, err
		}
		spec.IP = ip.String()
	}

	opts := types.EndpointCreateOptions{
		IP:            spec.IP,
		Ports:         spec.Ports,
		Policy:        spec.Policy,
		BindStrategy:  spec.Strategy.Bind,
		RouteStrategy: spec.Strategy.Route,
		Domain:        domain,
	}

	e, err := em.Create(namespace, service, &opts)
	if err != nil {
		log.Errorf("%s> get endpoint error: %s", logPrefix, err.Error())
		return nil, err
	}
	return e, nil
}

func EndpointUpdate(e *types.Endpoint, spec types.SpecNetwork) (*types.Endpoint, error) {

	em := distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())

	opts := types.EndpointUpdateOptions{
		Ports:         spec.Ports,
		Policy:        spec.Policy,
		BindStrategy:  spec.Strategy.Bind,
		RouteStrategy: spec.Strategy.Route,
	}

	et, err := em.Update(e, &opts)
	if err != nil {
		log.Errorf("%s> get endpoint error: %s", logPrefix, err.Error())
		return nil, err
	}

	return et, nil
}

func EndpointRemove(e *types.Endpoint) error {
	em := distribution.NewEndpointModel(context.Background(), envs.Get().GetStorage())
	return em.Remove(e)
}
