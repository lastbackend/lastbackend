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

package endpoint

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

const logEndpointPrefix = "runtime:endpoint:>"

func Manage(ctx context.Context, key string, spec *types.EndpointSpec) error {
	log.Debugf("%s manage: %s", logEndpointPrefix, key)

	ep := envs.Get().GetState().Endpoints().GetEndpoint(key)

	if ep != nil {
		log.Debugf("%s check endpoint: %s", logEndpointPrefix, key)
		if !equal(spec, ep) {
			status, err := Update(ctx, key, ep, spec)
			if err != nil {
				log.Errorf("%s update error: %s", logEndpointPrefix, err.Error())
				return err
			}
			envs.Get().GetState().Endpoints().SetEndpoint(key, status)
			return nil
		}

		return nil
	}

	status, err := Create(ctx, key, spec)
	if err != nil {
		log.Errorf("%s create error: %s", logEndpointPrefix, err.Error())
		return err
	}

	envs.Get().GetState().Endpoints().SetEndpoint(key, status)
	return nil
}

func Restore(ctx context.Context) error {
	log.Debugf("%s restore", logEndpointPrefix)
	cpi := envs.Get().GetCPI()
	endpoints, err := cpi.Info(ctx)
	if err != nil {
		log.Errorf("%s restore error: %s", err.Error())
		return err
	}
	envs.Get().GetState().Endpoints().SetEndpoints(endpoints)
	return nil
}

func Create(ctx context.Context, key string, spec *types.EndpointSpec) (*types.EndpointStatus, error) {
	log.Debugf("%s create %s", logEndpointPrefix, key)
	cpi := envs.Get().GetCPI()
	return cpi.Create(ctx, spec)
}

func Update(ctx context.Context, endpoint string, status *types.EndpointStatus, spec *types.EndpointSpec) (*types.EndpointStatus, error) {
	log.Debugf("%s update %s", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Update(ctx, status, spec)
}

func Destroy(ctx context.Context, endpoint string, status *types.EndpointStatus) error {
	log.Debugf("%s destroy", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Destroy(ctx, status)
}

func equal(spec *types.EndpointSpec, status *types.EndpointStatus) bool {
	if status.IP != spec.IP {
		log.Debugf("%s ips not match %s != %s", logEndpointPrefix, spec.IP, status.IP)
		return false
	}

	if spec.Strategy.Route != spec.Strategy.Route {
		log.Debugf("%s route strategy not match %s != %s", logEndpointPrefix, spec.Strategy.Route, status.Strategy.Route)
		return false
	}

	if spec.Strategy.Bind != spec.Strategy.Bind {
		log.Debugf("%s bind strategy not match %s != %s", logEndpointPrefix, spec.Strategy.Bind, status.Strategy.Bind)
		return false
	}

	for port, pm := range spec.PortMap {

		if _, ok := status.PortMap[port]; !ok {
			log.Debugf("%s portmap not found %#v", logEndpointPrefix, pm)
			return false
		}

		if status.PortMap[port] != pm {
			log.Debugf("%s portmap not match %#v != %#v", logEndpointPrefix, pm, status.PortMap[port])
			return false
		}
	}

	for port, pm := range status.PortMap {
		if _, ok := spec.PortMap[port]; !ok {
			log.Debugf("%s portmap should be deleted %#v", logEndpointPrefix, pm)
			return false
		}
	}

	if len(spec.Upstreams) != len(status.Upstreams) {
		log.Debugf("%s upstreams count changed %d != %d", logEndpointPrefix, len(spec.Upstreams), len(status.Upstreams))
		return false
	}

	for _, up := range spec.Upstreams {
		var f = false
		for _, stup := range status.Upstreams {
			if up == stup {
				f = true
			}
		}
		if !f {
			log.Debugf("%s upstream not found: %s", logEndpointPrefix, up)
			return false
		}
	}

	for _, up := range status.Upstreams {
		var f = false
		for _, stup := range spec.Upstreams {
			if up == stup {
				f = true
			}
		}
		if !f {
			log.Debugf("%s upstream should be deleted: %s", logEndpointPrefix, up)
			return false
		}
	}

	return true
}
