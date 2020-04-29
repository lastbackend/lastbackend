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

package network

import (
	"context"
	"github.com/lastbackend/lastbackend/tools/logger"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/network/state"
	)

const (
	logEndpointPrefix = "network:endpoint:>"
)

func (n *Network) Endpoints() *state.EndpointState {
	return n.state.Endpoints()
}

func (n *Network) EndpointManage(ctx context.Context, key string, manifest *models.EndpointManifest) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s manage: %s", logEndpointPrefix, key)

	state := n.state.Endpoints().GetEndpoint(key)

	if state != nil {
		if manifest.State == models.StateDestroy {
			n.EndpointDestroy(ctx, key, state)
			n.state.Endpoints().DelEndpoint(key)
			return nil
		}

		log.Debugf("%s check endpoint: %s", logEndpointPrefix, key)
		if !endpointEqual(manifest, state) {
			state, err := n.EndpointUpdate(ctx, key, state, manifest)
			if err != nil {
				log.Errorf("%s update error: %s", logEndpointPrefix, err.Error())
				return err
			}
			n.state.Endpoints().SetEndpoint(key, state)
			return nil
		}

		return nil
	}

	if manifest.State == models.StateDestroy {
		return nil
	}

	state, err := n.EndpointCreate(ctx, key, manifest)
	if err != nil {
		log.Errorf("%s create error: %s", logEndpointPrefix, err.Error())
		return err
	}

	n.state.Endpoints().SetEndpoint(key, state)
	return nil
}

func (n *Network) EndpointRestore(ctx context.Context) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s restore", logEndpointPrefix)
	cpi := n.cpi
	endpoints, err := cpi.Info(ctx)
	if err != nil {
		log.Errorf("%s restore error: %s", logEndpointPrefix, err.Error())
		return err
	}
	n.state.Endpoints().SetEndpoints(endpoints)
	return nil
}

func (n *Network) EndpointCreate(ctx context.Context, key string, manifest *models.EndpointManifest) (*models.EndpointState, error) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s create %s", logEndpointPrefix, key)
	cpi := n.cpi
	return cpi.Create(ctx, manifest)
}

func (n *Network) EndpointUpdate(ctx context.Context, endpoint string, state *models.EndpointState, manifest *models.EndpointManifest) (*models.EndpointState, error) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s update %s", logEndpointPrefix, endpoint)
	cpi := n.cpi
	return cpi.Update(ctx, state, manifest)
}

func (n *Network) EndpointDestroy(ctx context.Context, endpoint string, state *models.EndpointState) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s destroy: %s", logEndpointPrefix, endpoint)
	cpi := n.cpi
	return cpi.Destroy(ctx, state)
}

func endpointEqual(manifest *models.EndpointManifest, state *models.EndpointState) bool {
	log := logger.WithContext(context.Background())
	if state.IP != manifest.IP {
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
			log.Debugf("%s portmap not found %s", logEndpointPrefix, pm)
			return false
		}

		if state.PortMap[port] != pm {
			log.Debugf("%s portmap not match %s != %s", logEndpointPrefix, pm, state.PortMap[port])
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
