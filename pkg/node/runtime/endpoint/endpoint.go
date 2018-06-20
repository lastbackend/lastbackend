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

func Manage(ctx context.Context, key string, manifest *types.EndpointManifest) error {
	log.Debugf("%s manage: %s", logEndpointPrefix, key)

	state := envs.Get().GetState().Endpoints().GetEndpoint(key)

	if state != nil {
		if manifest.State == types.StateDestroy {
			Destroy(ctx, key, state)
			envs.Get().GetState().Endpoints().DelEndpoint(key)
			return nil
		}

		log.Debugf("%s check endpoint: %s", logEndpointPrefix, key)
		if !equal(manifest, state) {
			state, err := Update(ctx, key, state, manifest)
			if err != nil {
				log.Errorf("%s update error: %s", logEndpointPrefix, err.Error())
				return err
			}
			envs.Get().GetState().Endpoints().SetEndpoint(key, state)
			return nil
		}

		return nil
	}

	if manifest.State == types.StateDestroy {
		return nil
	}

	state, err := Create(ctx, key, manifest)
	if err != nil {
		log.Errorf("%s create error: %s", logEndpointPrefix, err.Error())
		return err
	}

	envs.Get().GetState().Endpoints().SetEndpoint(key, state)
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

func Create(ctx context.Context, key string, manifest *types.EndpointManifest) (*types.EndpointState, error) {
	log.Debugf("%s create %s", logEndpointPrefix, key)
	cpi := envs.Get().GetCPI()
	return cpi.Create(ctx, manifest)
}

func Update(ctx context.Context, endpoint string, state *types.EndpointState, manifest *types.EndpointManifest) (*types.EndpointState, error) {
	log.Debugf("%s update %s", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Update(ctx, state, manifest)
}

func Destroy(ctx context.Context, endpoint string, state *types.EndpointState) error {
	log.Debugf("%s destroy", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Destroy(ctx, state)
}

func equal(manifest *types.EndpointManifest, state *types.EndpointState) bool {
	if status.IP != spec.IP {
		log.Debugf("%s ips not match %s != %s", logEndpointPrefix, manifest.IP, state.IP)
		return false
	}

	if manifest.Strategy.Route != manifest.Strategy.Route {
		log.Debugf("%s route strategy not match %s != %s", logEndpointPrefix, manifest.Strategy.Route, state.Strategy.Route)
		return false
	}

	if manifest.Strategy.Bind != manifest.Strategy.Bind {
		log.Debugf("%s bind strategy not match %s != %s", logEndpointPrefix, manifest.Strategy.Bind, state.Strategy.Bind)
		return false
	}

	for port, pm := range manifest.PortMap {

		if _, ok := state.PortMap[port]; !ok {
			log.Debugf("%s portmap not found %#v", logEndpointPrefix, pm)
			return false
		}

		if state.PortMap[port] != pm {
			log.Debugf("%s portmap not match %#v != %#v", logEndpointPrefix, pm, state.PortMap[port])
			return false
		}
	}

	for port, pm := range state.PortMap {
		if _, ok := manifest.PortMap[port]; !ok {
			log.Debugf("%s portmap should be deleted %#v", logEndpointPrefix, pm)
			return false
		}
	}

	if len(manifest.Upstreams) != len(state.Upstreams) {
		log.Debugf("%s upstreams count changed %d != %d", logEndpointPrefix, len(manifest.Upstreams), len(state.Upstreams))
		return false
	}

	for _, up := range state.Upstreams {
		var f = false
		for _, stup := range state.Upstreams {
			if up == stup {
				f = true
			}
		}
		if !f {
			log.Debugf("%s upstream not found: %s", logEndpointPrefix, up)
			return false
		}
	}

	for _, up := range state.Upstreams {
		var f = false
		for _, stup := range manifest.Upstreams {
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
