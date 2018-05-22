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


// +build linux

package ipvs

import (
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cpi"
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/network"
	"fmt"
)

const logIPVSPrefix = "cpi:ipvs:proxy:>"

// Proxy balancer
type Proxy struct {
	cpi cpi.CPI
	// IVPS cmd path
	ipvs *IPVS
}

func (p *Proxy) Info(ctx context.Context) (map[string]*types.EndpointStatus, error) {
	el := make(map[string]*types.EndpointStatus)

	svcs, err := p.ipvs.GetServices(ctx)
	if err != nil {
		log.Errorf("%s info error: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	for _, svc := range svcs {

		// check if endpoint exists
		if _, ok := el[svc.Host]; !ok {
			el[svc.Host] = new(types.EndpointStatus)
			el[svc.Host].Upstreams = make([]string, 0)
		}
	}

	return el, nil
}

// Create new proxy rules
func (p *Proxy) Create(ctx context.Context, spec *types.EndpointSpec) (*types.EndpointStatus, error) {

	var (
		err  error
		status = new(types.EndpointStatus)
		csvcs = make([]*Service, 0)
	)

	status.IP = spec.IP
	svcs, err := specToServices(spec)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			for _, svc := range csvcs {
				p.ipvs.DelService(ctx, svc)
			}
		}
	}()

	for _, svc := range svcs {
		if err := p.ipvs.AddService(ctx, &svc); err != nil {
			status.State = types.StateError
			status.Message = err.Error()
		}
		csvcs = append(csvcs, &svc)
		return status, err
	}

	state, err := getStateByIP(ctx, spec.IP)
	if err != nil {
		log.Errorf("%s get state by ip err: %s", logIPVSPrefix, err.Error())
		return status, err
	}

	status.PortMap = state.PortMap
	status.Upstreams = state.Upstreams
	status.Strategy = state.Strategy

	status.State = types.StateReady
	status.Message = ""

	return status, nil
}

func (p *Proxy) Destroy(ctx context.Context, state *types.EndpointStatus) error {

	var (
		err  error
	)

	svcs, err := stateToServices(state)
	if err != nil {
		return err
	}

	for _, svc := range svcs {
		if err = p.ipvs.DelService(ctx, &svc); err != nil {
			log.Errorf("%s can not delete service: %s", logIPVSPrefix, err.Error())
		}
	}

	return err
}

func (p *Proxy) Update(ctx context.Context, state *types.EndpointStatus, spec *types.EndpointSpec) (*types.EndpointStatus, error) {

	var (
		status = new(types.EndpointStatus)
	)

	psvc, err := specToServices(spec)
	if err != nil {
		log.Errorf("%s can not convert spec to services: %s", logIPVSPrefix, err.Error())
		return state, err
	}

	csvc, err := stateToServices(state)
	if err != nil {
		log.Errorf("%s can not convert state to services: %s", logIPVSPrefix, err.Error())
		return state, err
	}

	for id, svc := range csvc {

		// remove service which not exists in new spec
		if _, ok := psvc[id]; !ok {
			if err := p.ipvs.DelService(ctx, &svc); err != nil {
				log.Errorf("%s can not remove service: %s", logIPVSPrefix, err.Error())
			}
		}

		// check service upstreams for removing
		for host, bknd := range svc.Backends {
			pu := psvc[id].Backends
			if _, ok := pu[host]; !ok {
				if err := p.ipvs.DelBackend(ctx, &svc, &bknd); err != nil {
					log.Errorf("%s can not remove backend: %s", logIPVSPrefix, err.Error())
				}
			}
		}

		// check service upstreams for creating
		for host, bknd := range psvc[id].Backends {
			if _, ok := svc.Backends[host]; !ok {
				if err := p.ipvs.AddBackend(ctx, &svc, &bknd); err != nil {
					log.Errorf("%s can not add backend: %s", logIPVSPrefix, err.Error())
				}
			}
		}
	}

	for id, svc := range psvc {
		if _, ok := csvc[id]; !ok {
			if err := p.ipvs.AddService(ctx, &svc); err != nil {
				log.Errorf("%s can not create service: %s", logIPVSPrefix, err.Error())
			}
		}
	}

	st, err := getStateByIP(ctx, spec.IP)
	if err != nil {
		log.Errorf("%s get state by ip err: %s", logIPVSPrefix, err.Error())
		return status, err
	}

	status.PortMap = st.PortMap
	status.Upstreams = st.Upstreams
	status.Strategy = st.Strategy

	status.State = types.StateReady
	status.Message = ""

	return status, nil
}

func specToServices(spec *types.EndpointSpec) (map[string]Service, error) {

	var svcs = make(map[string]Service, 0)

	for ext, pm := range spec.PortMap {

		port, proto, err := network.ParsePortMap(pm)
		if err != nil {
			err = errors.New("Invalid port map declaration")
			return svcs, err
		}

		svc := Service{
			Host: spec.IP,
			Port: ext,
		}

		for _, host := range spec.Upstreams {
			svc.Backends[host] = Backend{
				Host: host,
				Port: port,
			}
		}

		switch proto {
		case "tcp":
			svc.Type = proxyTCPProto
			svcs[fmt.Sprintf("%s_%d_%d_%s", svc.Host, svc.Port, port, proxyTCPProto)] = svc
			break
		case "udp":
			svc.Type = proxyUDPProto
			svcs[fmt.Sprintf("%s_%d_%d_%s", svc.Host, svc.Port, port, proxyUDPProto)] = svc
			break
		case "*":
			svcc := svc
			svc.Type = proxyTCPProto
			svcc.Type = proxyUDPProto

			svcs[fmt.Sprintf("%s_%d_%d_%s", svc.Host, svc.Port, port, proxyTCPProto)] = svc
			svcs[fmt.Sprintf("%s_%d_%d_%s", svcc.Host, svcc.Port, port, proxyUDPProto)] = svcc
			break
		}
	}

	return svcs, nil
}

func stateToServices(status *types.EndpointStatus) (map[string]Service, error) {

	var svcs = make(map[string]Service, 0)

	for ext, pm := range status.PortMap {

		port, proto, err := network.ParsePortMap(pm)
		if err != nil {
			err = errors.New("Invalid port map declaration")
			return svcs, err
		}

		svc := Service{
			Host: status.IP,
			Port: ext,
		}

		for _, host := range status.Upstreams {
			svc.Backends[host] = Backend{
				Host: host,
				Port: port,
			}
		}

		switch proto {
		case "tcp":
			svc.Type = proxyTCPProto
			svcs[fmt.Sprintf("%s_%d_%d_%s", svc.Host, svc.Port, port, proxyTCPProto)] = svc
			break
		case "udp":
			svc.Type = proxyUDPProto
			svcs[fmt.Sprintf("%s_%d_%d_%s", svc.Host, svc.Port, port, proxyUDPProto)] = svc
			break
		case "*":
			svcc := svc
			svc.Type = proxyTCPProto
			svcc.Type = proxyUDPProto

			svcs[fmt.Sprintf("%s_%d_%d_%s", svc.Host, svc.Port, port, proxyTCPProto)] = svc
			svcs[fmt.Sprintf("%s_%d_%d_%s", svcc.Host, svcc.Port, port, proxyUDPProto)] = svcc
			break
		}
	}

	return svcs, nil
}

func getStateByIP(ctx context.Context, ip string) (*types.EndpointStatus, error) {
	var status = new(types.EndpointStatus)
	return status, nil
}

func New() (*Proxy, error) {
	prx := new(Proxy)
	// TODO: Check ipvs proxy mode is available on host
	return prx, nil
}
