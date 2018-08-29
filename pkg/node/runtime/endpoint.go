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

package runtime

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
)

const (
	logEndpointPrefix = "runtime:endpoint:>"
)

func EndpointManage(ctx context.Context, key string, manifest *types.EndpointManifest) error {
	log.V(logLevel).Debugf("%s manage: %s", logEndpointPrefix, key)

	state := envs.Get().GetState().Endpoints().GetEndpoint(key)

	if state != nil {
		if manifest.State == types.StateDestroy {
			EndpointDestroy(ctx, key, state)
			envs.Get().GetState().Endpoints().DelEndpoint(key)
			return nil
		}

		log.V(logLevel).Debugf("%s check endpoint: %s", logEndpointPrefix, key)
		if !endpointEqual(manifest, state) {
			state, err := EndpointUpdate(ctx, key, state, manifest)
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

	state, err := EndpointCreate(ctx, key, manifest)
	if err != nil {
		log.Errorf("%s create error: %s", logEndpointPrefix, err.Error())
		return err
	}

	envs.Get().GetState().Endpoints().SetEndpoint(key, state)
	return nil
}

func EndpointRestore(ctx context.Context) error {
	log.V(logLevel).Debugf("%s restore", logEndpointPrefix)
	cpi := envs.Get().GetCPI()
	endpoints, err := cpi.Info(ctx)
	if err != nil {
		log.Errorf("%s restore error: %s", err.Error())
		return err
	}
	envs.Get().GetState().Endpoints().SetEndpoints(endpoints)
	return nil
}

func EndpointCreate(ctx context.Context, key string, manifest *types.EndpointManifest) (*types.EndpointState, error) {
	log.V(logLevel).Debugf("%s create %s", logEndpointPrefix, key)
	cpi := envs.Get().GetCPI()
	return cpi.Create(ctx, manifest)
}

func EndpointUpdate(ctx context.Context, endpoint string, state *types.EndpointState, manifest *types.EndpointManifest) (*types.EndpointState, error) {
	log.V(logLevel).Debugf("%s update %s", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Update(ctx, state, manifest)
}

func EndpointDestroy(ctx context.Context, endpoint string, state *types.EndpointState) error {
	log.V(logLevel).Debugf("%s destroy", logEndpointPrefix, endpoint)
	cpi := envs.Get().GetCPI()
	return cpi.Destroy(ctx, state)
}

func endpointEqual(manifest *types.EndpointManifest, state *types.EndpointState) bool {

	if state.IP != manifest.IP {
		log.V(logLevel).Debugf("%s ips not match %s != %s", logEndpointPrefix, manifest.IP, state.IP)
		return false
	}

	if manifest.Strategy.Route != manifest.Strategy.Route {
		log.V(logLevel).Debugf("%s route strategy not match %s != %s", logEndpointPrefix, manifest.Strategy.Route, state.Strategy.Route)
		return false
	}

	if manifest.Strategy.Bind != manifest.Strategy.Bind {
		log.V(logLevel).Debugf("%s bind strategy not match %s != %s", logEndpointPrefix, manifest.Strategy.Bind, state.Strategy.Bind)
		return false
	}

	for port, pm := range manifest.PortMap {

		if _, ok := state.PortMap[port]; !ok {
			log.V(logLevel).Debugf("%s portmap not found %#v", logEndpointPrefix, pm)
			return false
		}

		if state.PortMap[port] != pm {
			log.V(logLevel).Debugf("%s portmap not match %#v != %#v", logEndpointPrefix, pm, state.PortMap[port])
			return false
		}
	}

	for port, pm := range state.PortMap {
		if _, ok := manifest.PortMap[port]; !ok {
			log.V(logLevel).Debugf("%s portmap should be deleted %#v", logEndpointPrefix, pm)
			return false
		}
	}

	if len(manifest.Upstreams) != len(state.Upstreams) {
		log.V(logLevel).Debugf("%s upstreams count changed %d != %d", logEndpointPrefix, len(manifest.Upstreams), len(state.Upstreams))
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
			log.V(logLevel).Debugf("%s upstream not found: %s", logEndpointPrefix, up)
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
			log.V(logLevel).Debugf("%s upstream should be deleted: %s", logEndpointPrefix, up)
			return false
		}
	}

	return true
}
