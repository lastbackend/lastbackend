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
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni/utils"
	"github.com/spf13/viper"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"net"
	"os/exec"
	"strings"
	"syscall"

	libipvs "github.com/docker/libnetwork/ipvs"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/runtime/cpi"
	"github.com/lastbackend/lastbackend/pkg/util/network"
	"github.com/vishvananda/netlink/nl"
)

const (
	logIPVSPrefix = "cpi:ipvs:proxy:>"
	logLevel      = 3
	ifaceName     = "lb-ipvs"
	ifaceDocker   = "docker0"
)

// Proxy balancer
type Proxy struct {
	cpi cpi.CPI
	// IVPS cmd path
	ipvs *libipvs.Handle
	link netlink.Link
	dest struct {
		external net.IP
		internal net.IP
	}
}

func (p *Proxy) Info(ctx context.Context) (map[string]*types.EndpointState, error) {
	return p.getState(ctx)
}

// Create new proxy rules
func (p *Proxy) Create(ctx context.Context, manifest *types.EndpointManifest) (*types.EndpointState, error) {

	log.V(logLevel).Debugf("%s create ipvs virtual server with ip %s: and upstreams %v", logIPVSPrefix, manifest.IP, manifest.Upstreams)

	var (
		err   error
		csvcs = make([]*libipvs.Service, 0)
	)

	svcs, dests, err := specToServices(manifest)
	if err != nil {
		log.Errorf("%s can not get services from manifest: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	if len(dests) == 0 {
		log.V(logLevel).Debugf("%s skip creating service, destinations not exists", logIPVSPrefix)
		return nil, nil
	}

	defer func() {
		if err != nil {
			for _, svc := range csvcs {
				log.V(logLevel).Debugf("%s delete service: %s", logIPVSPrefix, svc.Address.String())
				p.ipvs.DelService(svc)
			}
		}
	}()

	for _, svc := range svcs {
		log.V(logLevel).Debugf("%s create new service: %s", logIPVSPrefix, svc.Address.String())
		if err := p.ipvs.NewService(svc); err != nil {
			log.Errorf("%s create service err: %s", logIPVSPrefix, err.Error())
		}

		for _, dest := range dests {
			log.V(logLevel).Debugf("%s create new destination %s for service: %s", logIPVSPrefix,
				dest.Address.String(), svc.Address.String())

			if err := p.ipvs.NewDestination(svc, dest); err != nil {
				log.Errorf("%s create destination for service err: %s", logIPVSPrefix, err.Error())
			}
		}

		csvcs = append(csvcs, svc)
	}

	log.Debugf("%s check ip %s is binded to link %s", logIPVSPrefix, manifest.IP, p.link.Attrs().Name)

	var dest net.IP

	if manifest.External {
		dest = p.dest.external
	} else {
		dest = p.dest.internal
	}

	if err := p.addIpBindToLink(manifest.IP, dest); err != nil {
		log.Warnf("%s failed bind ip to link err: %s", logIPVSPrefix, err.Error())
	}

	state, err := p.getStateByIP(ctx, manifest.IP)
	if err != nil {
		log.Errorf("%s get state by ip err: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	return state, nil
}

// Destroy proxy rules
func (p *Proxy) Destroy(ctx context.Context, state *types.EndpointState) error {

	var (
		err error
	)

	mf := types.EndpointManifest{}
	mf.EndpointSpec = state.EndpointSpec
	mf.Upstreams = state.Upstreams

	svcs, _, err := specToServices(&mf)
	if err != nil {
		return err
	}

	for _, svc := range svcs {
		if err = p.ipvs.DelService(svc); err != nil {
			log.Errorf("%s can not delete service: %s", logIPVSPrefix, err.Error())
		}

		if err := p.delIpBindToLink(svc.Address.String()); err != nil {
			log.Errorf("%s can not unbind ip from link: %s", logIPVSPrefix, err.Error())
		}
	}

	return err
}

// Update proxy rules
func (p *Proxy) Update(ctx context.Context, state *types.EndpointState, spec *types.EndpointManifest) (*types.EndpointState, error) {

	psvc, pdest, err := specToServices(spec)
	if err != nil {
		log.Errorf("%s can not convert spec to services: %s", logIPVSPrefix, err.Error())
		return state, err
	}

	mf := types.EndpointManifest{}
	mf.EndpointSpec = state.EndpointSpec
	mf.Upstreams = state.Upstreams

	csvc, cdest, err := specToServices(&mf)
	if err != nil {
		log.Errorf("%s can not convert state to services: %s", logIPVSPrefix, err.Error())
		return state, err
	}

	for id, svc := range csvc {

		log.Debugf("%s check old service: %s", logIPVSPrefix, id)
		// remove service which not exists in new spec
		if _, ok := psvc[id]; !ok {
			log.Debugf("%s delete service: %s", logIPVSPrefix, id)
			if err := p.ipvs.DelService(svc); err != nil {
				log.Errorf("%s can not remove service: %s", logIPVSPrefix, err.Error())
			}
			continue
		}
	}

	for id, svc := range psvc {
		log.Debugf("%s check new service: %s", logIPVSPrefix, id)

		if _, ok := csvc[id]; !ok {
			log.Debugf("%s create service: %s", logIPVSPrefix, id)
			if err := p.ipvs.NewService(svc); err != nil {
				log.Errorf("%s can not create service: %s", logIPVSPrefix, err.Error())
			}
		}

		// check service upstreams for removing
		for did, dest := range cdest {
			log.Debugf("%s check service %s old backend exists %s", logIPVSPrefix, id, did)
			if _, ok := pdest[did]; !ok {
				log.Debugf("%s service %s backend delete %s", logIPVSPrefix, id, did)
				if err := p.ipvs.DelDestination(svc, dest); err != nil {
					log.Errorf("%s can not remove backend: %s", logIPVSPrefix, err.Error())
				}
			}
		}

		// check service upstreams for creating
		for did, dest := range pdest {
			log.Debugf("%s check service %s new backend exists %s", logIPVSPrefix, id, did)
			if _, ok := cdest[did]; !ok {
				log.Debugf("%s service %s backend create %s", logIPVSPrefix, id, did)
				if err := p.ipvs.NewDestination(svc, dest); err != nil {
					log.Errorf("%s can not add backend: %s", logIPVSPrefix, err.Error())
				}
			}
		}
	}

	log.Debugf("Check ip %s is binded to link %s", spec.IP, p.link.Attrs().Name)

	var dest net.IP

	if spec.External {
		dest = p.dest.external
	} else {
		dest = p.dest.internal
	}

	if err := p.addIpBindToLink(spec.IP, dest); err != nil {
		log.Warnf("%s failed bind ip to link err: %s", logIPVSPrefix, err.Error())
	}

	st, err := p.getStateByIP(ctx, spec.IP)
	if err != nil {
		log.Errorf("%s get state by ip err: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	return st, nil
}

// getStateByIp returns current proxy state filtered by endpoint ip
func (p *Proxy) getStateByIP(ctx context.Context, ip string) (*types.EndpointState, error) {

	state, err := p.getState(ctx)
	if err != nil {
		log.Errorf("%s get state err: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	return state[ip], nil
}

// getStateByIp returns current proxy state
func (p *Proxy) getState(ctx context.Context) (map[string]*types.EndpointState, error) {

	el := make(map[string]*types.EndpointState)

	if out, err := exec.Command("modprobe", "-va", "ip_vs").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("%s running modprobe ip_vs failed with message: `%s`, error: %s", logIPVSPrefix, strings.TrimSpace(string(out)), err.Error())
	}

	svcs, err := p.ipvs.GetServices()
	if err != nil {
		log.Errorf("%s info error: %s", logIPVSPrefix, err.Error())
		return el, err
	}

	log.V(logLevel).Debugf("%s services list: %#v", logIPVSPrefix, svcs)

	var ips = make(map[string]bool, 0)

	for _, svc := range svcs {

		ips[svc.Address.String()] = true

		// check if endpoint exists
		var host = svc.Address.String()

		log.V(logLevel).Debugf("%s add service %s to current state", logIPVSPrefix, svc.Address.String())

		endpoint := el[host]
		if endpoint == nil {
			endpoint = new(types.EndpointState)
			endpoint.IP = host
			endpoint.PortMap = make(map[uint16]string)
			endpoint.Upstreams = make([]string, 0)
		}

		var prt uint16

		dests, err := p.ipvs.GetDestinations(svc)
		if err != nil {
			log.Errorf("%s get destinations err: %s", logIPVSPrefix, err.Error())
			continue
		}

		log.V(logLevel).Debugf("%s found %d destinations for service: %s", logIPVSPrefix, len(dests), svc.Address.String())

		for _, dest := range dests {

			var (
				f = false
			)

			if prt == 0 {
				prt = dest.Port
			}

			if prt != 0 && prt != dest.Port {
				log.V(logLevel).Debugf("%s dest port mismatch %d != %d", logIPVSPrefix, prt, dest.Port)
				break
			}

			for _, hst := range endpoint.Upstreams {
				if dest.Address.String() == hst {
					f = true
					break
				}
			}

			if !f {
				endpoint.Upstreams = append(endpoint.Upstreams, dest.Address.String())
			}
		}

		if prt != 0 {
			if _, ok := endpoint.PortMap[svc.Port]; ok {
				if svc.Protocol == syscall.IPPROTO_TCP && (endpoint.PortMap[svc.Port] == fmt.Sprintf("%d/%s", prt, proxyUDPProto)) {
					endpoint.PortMap[svc.Port] = fmt.Sprintf("%d/*", prt)
				}

				if svc.Protocol == syscall.IPPROTO_UDP && (endpoint.PortMap[svc.Port] == fmt.Sprintf("%d/%s", prt, proxyTCPProto)) {
					endpoint.PortMap[svc.Port] = fmt.Sprintf("%d/*", prt)
				}
			} else {

				if svc.Protocol == syscall.IPPROTO_TCP {
					endpoint.PortMap[svc.Port] = fmt.Sprintf("%d/%s", prt, proxyTCPProto)
				}

				if svc.Protocol == syscall.IPPROTO_UDP {
					endpoint.PortMap[svc.Port] = fmt.Sprintf("%d/%s", prt, proxyUDPProto)
				}
			}
		}

		el[host] = endpoint
	}

	return el, nil
}

func (p *Proxy) addIpBindToLink(ip string, dest net.IP) error {

	ipn := net.ParseIP(ip)
	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/32", ipn.String()))
	if err != nil {
		log.Errorf("%s can not parse IP %s; %s", logIPVSPrefix, ip, err.Error())
		return err
	}

	addrs, err := netlink.AddrList(p.link, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("%s can not fetch IPs: %s", logIPVSPrefix, err.Error())
		return err
	}

	var exists = false
	for _, a := range addrs {
		if a.IP.String() == addr.IP.String() {
			exists = true
			break
		}
	}
	if !exists {
		netlink.AddrAdd(p.link, addr)
	}

	routes, err := netlink.RouteGet(ipn)
	if err != nil {
		log.Errorf("%s can not get routes for ip", logIPVSPrefix)
		return err
	}

	for _, route := range routes {

		if route.Dst.IP.Equal(ipn) {
			log.Debugf("%s replace route destination %s > %s", logIPVSPrefix, route.Dst.IP.String(), dest.String())
			route.Src = dest
			route.LinkIndex = p.link.Attrs().Index
			route.Scope = netlink.SCOPE_HOST
			route.Table = unix.RT_TABLE_LOCAL
		}

		if err := netlink.RouteReplace(&route); err != nil {
			log.Errorf("%s can not replace route: %s", logIPVSPrefix, err.Error())
			return err
		}
	}

	return nil
}

func (p *Proxy) delIpBindToLink(ip string) error {

	ipn := net.ParseIP(ip)
	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/32", ipn.String()))
	if err != nil {
		log.Errorf("%s can not parse IP %s; %s", logIPVSPrefix, ip, err.Error())
		return err
	}

	addrs, err := netlink.AddrList(p.link, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("%s can not fetch IPs:%s", logIPVSPrefix, err.Error())
		return err
	}

	var exists = false
	for _, a := range addrs {

		if a.IP.String() == addr.IP.String() {
			exists = true
			break
		}
	}
	if exists {
		if err := netlink.AddrDel(p.link, addr); err != nil {
			log.Errorf("%s can not remove link: %s", logIPVSPrefix, err.Error())
		}
	}

	return nil
}

func New() (*Proxy, error) {

	prx := new(Proxy)
	handler, err := libipvs.New("")
	if err != nil {
		log.Errorf("%s can not initialize ipvs: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	prx.ipvs = handler

	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	f := false
	for _, link := range links {
		if link.Attrs().Name == ifaceName {
			log.Debugf("ipvs interface found %s", link.Attrs().Name)
			f = true
			prx.link = link
			break
		}
	}

	if !f {
		link := netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{
				Name: ifaceName,
			},
		}

		log.Debugf("%s ipvs interface not found: create new", logIPVSPrefix)
		if err := netlink.LinkAdd(&link); err != nil {
			if err == syscall.EEXIST {
				log.V(logLevel).Debugf("%s device already exists: %s", logIPVSPrefix, link.Name)

				l, err := netlink.LinkByName(link.Name)
				if err != nil {
					log.Errorf("%s link by name: %s", logIPVSPrefix, err.Error())
				}

				prx.link = l.(*netlink.Vxlan)
			} else {
				log.Errorf("%s can not create ipvs dummy interface: %s", logIPVSPrefix, err.Error())
				return nil, err
			}
		}

		prx.link = &link
	}

	var (
		eiface = viper.GetString("runtime.cpi.interface.external")
		iiface = viper.GetString("runtime.cpi.interface.internal")
	)

	if eiface == types.EmptyString {
		log.Debugf("%s find default interface to traffic route by name", logIPVSPrefix)
		_, prx.dest.external, err = utils.GetDefaultInterface()
		if err != nil {
			return nil, err
		}

		log.Debugf("%s external route ip net: %s", logIPVSPrefix, prx.dest.external.String())
	} else {
		log.Debugf("%s find interface to traffic route by name: %s", logIPVSPrefix, eiface)
		_, prx.dest.external, err = utils.GetIfaceByName(eiface)
		if err != nil {
			return nil, err
		}
		log.Debugf("%s external route ip net: %s", logIPVSPrefix, prx.dest.external.String())
	}

	if iiface == types.EmptyString {
		iiface = ifaceDocker
	}

	log.Debugf("%s find interface to traffic route by name: %s", logIPVSPrefix, iiface)
	_, prx.dest.internal, err = utils.GetIfaceByName(iiface)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s internal route ip net: %s", logIPVSPrefix, prx.dest.internal.String())

	// TODO: Check ipvs proxy mode is available on host
	return prx, nil
}

func specToServices(spec *types.EndpointManifest) (map[string]*libipvs.Service, map[string]*libipvs.Destination, error) {

	var svcs = make(map[string]*libipvs.Service, 0)
	dests := make(map[string]*libipvs.Destination, 0)

	for ext, pm := range spec.PortMap {

		port, proto, err := network.ParsePortMap(pm)
		if err != nil {
			err = errors.New("Invalid port map declaration")
			return svcs, dests, err
		}

		svc := libipvs.Service{
			Address:       net.ParseIP(spec.IP),
			Port:          ext,
			AddressFamily: nl.FAMILY_V4,
			SchedName:     "rr",
		}

		for _, host := range spec.Upstreams {
			log.Debugf("%s: add new destination to spec for: %s", logIPVSPrefix, host)
			dest := new(libipvs.Destination)

			dest.Address = net.ParseIP(host)
			dest.Port = port
			dest.Weight = 1
			dests[fmt.Sprintf("%s_%d", dest.Address.String(), dest.Port)] = dest
			log.Debugf("%s: added new destination %s_%d", logIPVSPrefix, dest.Address.String(), dest.Port)
		}

		switch proto {
		case "tcp":
			svc.Protocol = syscall.IPPROTO_TCP
			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svc.Port, port, proxyTCPProto)] = &svc
			break
		case "udp":
			svc.Protocol = syscall.IPPROTO_UDP
			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svc.Port, port, proxyUDPProto)] = &svc
			break
		case "*":
			svcc := svc
			svc.Protocol = syscall.IPPROTO_TCP
			svcc.Protocol = syscall.IPPROTO_UDP

			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svc.Port, port, proxyTCPProto)] = &svc
			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svcc.Port, port, proxyUDPProto)] = &svcc
			break
		}
	}

	return svcs, dests, nil
}
