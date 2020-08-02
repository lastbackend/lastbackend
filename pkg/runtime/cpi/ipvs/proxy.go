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

// +build linux

package ipvs

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/runtime/cni/utils"
	"github.com/spf13/viper"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/network"
	libipvs "github.com/moby/ipvs"
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
	// IVPS cmd path
	ipvs *libipvs.Handle
	link netlink.Link
	dest struct {
		external net.IP
		internal net.IP
	}
}

type Service struct {
	srvc *libipvs.Service
	dest map[string]*libipvs.Destination
}

func (p *Proxy) Info(ctx context.Context) (map[string]*models.EndpointState, error) {
	return p.getState(ctx)
}

// Create new proxy rules
func (p *Proxy) Create(ctx context.Context, manifest *models.EndpointManifest) (*models.EndpointState, error) {

	log := logger.WithContext(ctx)
	log.Debugf("%s create ipvs virtual server with ip %s: and upstreams %v", logIPVSPrefix, manifest.IP, manifest.Upstreams)

	var (
		err   error
		csvcs = make([]*Service, 0)
	)

	svcs, err := specToServices(manifest)
	if err != nil {
		log.Errorf("%s can not be get services from manifest: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	defer func() {
		if err != nil {
			for _, svc := range csvcs {
				log.Debugf("%s delete service: %s", logIPVSPrefix, svc.srvc.Address.String())
				p.ipvs.DelService(svc.srvc)
			}
		}
	}()

	for _, svc := range svcs {

		if len(svc.dest) == 0 {
			log.Debugf("%s skip creating service, destinations not exists", logIPVSPrefix)
			return nil, nil
		}

		log.Debugf("%s create new service: %s", logIPVSPrefix, svc.srvc.Address.String())
		if err := p.ipvs.NewService(svc.srvc); err != nil {
			log.Errorf("%s create service err: %s", logIPVSPrefix, err.Error())
		}

		for _, dest := range svc.dest {
			log.Debugf("%s create new destination %s for service: %s", logIPVSPrefix,
				dest.Address.String(), svc.srvc.Address.String())

			if err := p.ipvs.NewDestination(svc.srvc, dest); err != nil {
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
func (p *Proxy) Destroy(ctx context.Context, state *models.EndpointState) error {

	var (
		err error
		log = logger.WithContext(ctx)
	)

	if state == nil {
		return nil
	}

	mf := models.EndpointManifest{}
	mf.EndpointSpec = state.EndpointSpec
	mf.Upstreams = state.Upstreams

	svcs, err := specToServices(&mf)
	if err != nil {
		return err
	}

	for _, svc := range svcs {
		if err = p.ipvs.DelService(svc.srvc); err != nil {
			log.Errorf("%s can not be delete service: %s", logIPVSPrefix, err.Error())
		}

		if err := p.delIpBindToLink(svc.srvc.Address.String()); err != nil {
			log.Errorf("%s can not be unbind ip from link: %s", logIPVSPrefix, err.Error())
		}
	}

	return err
}

// Update proxy rules
func (p *Proxy) Update(ctx context.Context, state *models.EndpointState, spec *models.EndpointManifest) (*models.EndpointState, error) {

	log := logger.WithContext(ctx)

	psvc, err := specToServices(spec)
	if err != nil {
		log.Errorf("%s can not be convert spec to services: %s", logIPVSPrefix, err.Error())
		return state, err
	}

	mf := models.EndpointManifest{}
	mf.EndpointSpec = state.EndpointSpec
	mf.Upstreams = state.Upstreams

	csvc, err := specToServices(&mf)
	if err != nil {
		log.Errorf("%s can not be convert state to services: %s", logIPVSPrefix, err.Error())
		return state, err
	}

	for id, svc := range csvc {

		log.Debugf("%s check old service: %s", logIPVSPrefix, id)
		// remove service which not exists in new spec
		if _, ok := psvc[id]; !ok {
			log.Debugf("%s delete service: %s", logIPVSPrefix, id)
			if err := p.ipvs.DelService(svc.srvc); err != nil {
				log.Errorf("%s can not be remove service: %s", logIPVSPrefix, err.Error())
			}
			continue
		}
	}

	for id, svc := range psvc {
		log.Debugf("%s check new service: %s", logIPVSPrefix, id)

		if _, ok := csvc[id]; !ok {
			log.Debugf("%s create service: %s", logIPVSPrefix, id)
			if err := p.ipvs.NewService(svc.srvc); err != nil {
				log.Errorf("%s can not be create service: %s", logIPVSPrefix, err.Error())
			}
		} else {
			// check service upstreams for removing
			for did, dest := range csvc[id].dest {
				log.Debugf("%s check service %s old backend exists %s", logIPVSPrefix, id, did)
				if _, ok := svc.dest[did]; !ok {
					log.Debugf("%s service %s backend delete %s", logIPVSPrefix, id, did)
					if err := p.ipvs.DelDestination(svc.srvc, dest); err != nil {
						log.Errorf("%s can not be remove backend: %s", logIPVSPrefix, err.Error())
					}
				}
			}
		}

		// check service upstreams for creating
		for did, dest := range svc.dest {
			log.Debugf("%s check service %s new backend exists %s", logIPVSPrefix, id, did)

			if _, ok := csvc[id]; !ok {
				if err := p.ipvs.NewDestination(svc.srvc, dest); err != nil {
					log.Errorf("%s can not be add backend: %s", logIPVSPrefix, err.Error())
				}
			} else {
				if _, ok := csvc[id].dest[did]; !ok {
					log.Debugf("%s service %s backend create %s", logIPVSPrefix, id, did)
					if err := p.ipvs.NewDestination(svc.srvc, dest); err != nil {
						log.Errorf("%s can not be add backend: %s", logIPVSPrefix, err.Error())
					}
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
func (p *Proxy) getStateByIP(ctx context.Context, ip string) (*models.EndpointState, error) {

	log := logger.WithContext(ctx)

	state, err := p.getState(ctx)
	if err != nil {
		log.Errorf("%s get state err: %s", logIPVSPrefix, err.Error())
		return nil, err
	}

	return state[ip], nil
}

// getStateByIp returns current proxy state
func (p *Proxy) getState(ctx context.Context) (map[string]*models.EndpointState, error) {

	log := logger.WithContext(ctx)
	el := make(map[string]*models.EndpointState)

	if out, err := exec.Command("modprobe", "-va", "ip_vs").CombinedOutput(); err != nil {
		return nil, fmt.Errorf("%s running modprobe ip_vs failed with message: `%s`, error: %s", logIPVSPrefix, strings.TrimSpace(string(out)), err.Error())
	}

	svcs, err := p.ipvs.GetServices()
	if err != nil {
		log.Errorf("%s info error: %s", logIPVSPrefix, err.Error())
		return el, err
	}

	var ips = make(map[string]bool, 0)

	for _, svc := range svcs {

		ips[svc.Address.String()] = true

		// check if endpoint exists
		var host = svc.Address.String()

		endpoint := el[host]
		if endpoint == nil {
			endpoint = new(models.EndpointState)
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

		for _, dest := range dests {

			var (
				f = false
			)

			if prt == 0 {
				prt = dest.Port
			}

			if prt != 0 && prt != dest.Port {
				log.Debugf("%s dest port mismatch %d != %d", logIPVSPrefix, prt, dest.Port)
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


	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	ipn := net.ParseIP(ip)
	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/32", ipn.String()))
	if err != nil {
		log.Errorf("%s can not be parse IP %s; %s", logIPVSPrefix, ip, err.Error())
		return err
	}

	addrs, err := netlink.AddrList(p.link, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("%s can not be fetch IPs: %s", logIPVSPrefix, err.Error())
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
		log.Errorf("%s can not be get routes for ip", logIPVSPrefix)
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
			log.Errorf("%s can not be replace route: %s", logIPVSPrefix, err.Error())
			return err
		}
	}

	return nil
}

func (p *Proxy) delIpBindToLink(ip string) error {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	ipn := net.ParseIP(ip)
	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/32", ipn.String()))
	if err != nil {
		log.Errorf("%s can not be parse IP %s; %s", logIPVSPrefix, ip, err.Error())
		return err
	}

	addrs, err := netlink.AddrList(p.link, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("%s can not be fetch IPs:%s", logIPVSPrefix, err.Error())
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
			log.Errorf("%s can not be remove link: %s", logIPVSPrefix, err.Error())
		}
	}

	return nil
}

func New(v *viper.Viper) (*Proxy, error) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	prx := new(Proxy)
	handler, err := libipvs.New("")
	if err != nil {
		log.Errorf("%s can not be initialize ipvs: %s", logIPVSPrefix, err.Error())
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
				log.Debugf("%s device already exists: %s", logIPVSPrefix, link.Name)

				l, err := netlink.LinkByName(link.Name)
				if err != nil {
					log.Errorf("%s link by name: %s", logIPVSPrefix, err.Error())
				}

				prx.link = l.(*netlink.Vxlan)
			} else {
				log.Errorf("%s can not be create ipvs dummy interface: %s", logIPVSPrefix, err.Error())
				return nil, err
			}
		}

		prx.link = &link
	}

	var (
		eiface = v.GetString("network.cpi.interface.external")
		iiface = v.GetString("network.cpi.interface.internal")
	)

	if eiface == models.EmptyString {
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

	if iiface == models.EmptyString {
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

func specToServices(spec *models.EndpointManifest) (map[string]*Service, error) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	var svcs = make(map[string]*Service, 0)

	for ext, pm := range spec.PortMap {

		port, proto, err := network.ParsePortMap(pm)
		if err != nil {
			err = errors.New("Invalid port map declaration")
			return svcs, err
		}

		svc := new(Service)
		svc.srvc = &libipvs.Service{
			Address:       net.ParseIP(spec.IP),
			Port:          ext,
			AddressFamily: nl.FAMILY_V4,
			SchedName:     "rr",
		}
		svc.dest = make(map[string]*libipvs.Destination, 0)

		for _, host := range spec.Upstreams {
			log.Debugf("%s: add new destination to spec for: %s", logIPVSPrefix, host)

			dest := new(libipvs.Destination)

			dest.Address = net.ParseIP(host)
			dest.Port = port
			dest.Weight = 1
			svc.dest[fmt.Sprintf("%s_%d", dest.Address.String(), dest.Port)] = dest
			log.Debugf("%s: added new destination %s_%d", logIPVSPrefix, dest.Address.String(), dest.Port)
		}

		switch proto {
		case "tcp":
			svc.srvc.Protocol = syscall.IPPROTO_TCP
			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svc.srvc.Port, port, proxyTCPProto)] = svc
			break
		case "udp":
			svc.srvc.Protocol = syscall.IPPROTO_UDP
			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svc.srvc.Port, port, proxyUDPProto)] = svc
			break
		case "*":
			svcc := *svc
			svc.srvc.Protocol = syscall.IPPROTO_TCP
			svcc.srvc.Protocol = syscall.IPPROTO_UDP

			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svc.srvc.Port, port, proxyTCPProto)] = svc
			svcs[fmt.Sprintf("%s_%d_%d_%s", spec.IP, svcc.srvc.Port, port, proxyUDPProto)] = &svcc
			break
		}
	}

	return svcs, nil
}
