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
// +build linux

package vxlan

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni"
	"net"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni/utils"
	"github.com/vishvananda/netlink"
)

const NetworkType = "vxlan"
const DefaultContainerDevice = "docker0"

type Network struct {
	cni.CNI

	ExtIface *NetworkInterface
	IntIface *NetworkInterface
	Device   *Device
	Network  *net.IPNet
	CIDR     *net.IPNet
	IP       net.IP
}

type NetworkInterface struct {
	Iface     *net.Interface
	IfaceAddr net.IP
}

func New(iface string) (*Network, error) {

	var (
		nt  = new(Network)
		err error
	)

	nt.ExtIface = new(NetworkInterface)

	if iface == types.EmptyString {
		log.Debug("Add network to default interface")
		if nt.ExtIface.Iface, nt.ExtIface.IfaceAddr, err = utils.GetDefaultInterface(); err != nil {
			log.Errorf("Can not get default interface: %s", err.Error())
			return nt, err
		}
	} else {
		log.Debugf("Add network to interface: %s", iface)
		if nt.ExtIface.Iface, nt.ExtIface.IfaceAddr, err = utils.GetIfaceByName(iface); err != nil {
			log.Errorf("Can not get interface [%s]: %s", iface, err.Error())
			return nt, err
		}

	}

	log.Debugf("external interface: %s:%s", nt.ExtIface.Iface.Name, nt.ExtIface.IfaceAddr.String())

	if err := nt.SetSubnetFromDevice(DefaultContainerDevice); err != nil {
		log.Errorf("Can not set subnet: %s", err.Error())
		return nil, err
	}

	if err := nt.AddInterface(); err != nil {
		log.Errorf("Can not add interface: %s", err.Error())
		return nt, err
	}

	// Add forward rules for network range
	go utils.SetupAndEnsureIPTables(utils.ForwardRules(nt.Network.String()), 5)

	return nt, nil
}

func (n *Network) SetSubnetFromDevice(name string) error {

	iface, _, err := utils.GetIfaceByName(name)
	if err != nil {
		log.Errorf("Can not find interface by name %s", name)
		return err
	}

	if iface == nil {
		log.Errorf("Can not find interface by name %s", name)
		return errors.New("can not find interface")
	}

	addrs, err := netlink.AddrList(&netlink.Device{
		LinkAttrs: netlink.LinkAttrs{
			Index: iface.Index,
		},
	}, syscall.AF_INET)

	if err != nil {
		log.Errorf("Can not locate docker interface ips: %s", err.Error())
	}

	n.IntIface = new(NetworkInterface)
	n.IntIface.Iface = iface

	if len(addrs) == 0 {
		log.Error("docker interface has not IP address")
		panic(0)
	}

	sip := make(net.IP, len(addrs[0].IPNet.IP))
	smk := make(net.IPMask, len(addrs[0].Mask))

	n.IntIface.IfaceAddr = addrs[0].IP

	copy(sip, addrs[0].IPNet.IP)
	copy(smk, addrs[0].Mask)

	sip[3] = byte(0)
	n.CIDR = &net.IPNet{
		IP:   sip,
		Mask: smk,
	}

	n.Network = &net.IPNet{
		IP:   sip.Mask(sip.DefaultMask()),
		Mask: net.CIDRMask(8, 32),
	}

	return nil
}

func (n *Network) AddInterface() error {

	var err error

	if n.Device, err = NewDevice(DeviceCreateOpts{
		vni:   DeviceDefaultVNI,
		name:  fmt.Sprintf("%s%d", DeviceDefaultName, DeviceDefaultVNI),
		index: n.ExtIface.Iface.Index,
		addr:  n.ExtIface.IfaceAddr,
		port:  DeviceDefaultPort,
	}); err != nil {
		log.Errorf("Can not create xvlan interface: %s", err.Error())
		return err
	}

	n.Device.SetIP(*n.CIDR)
	return nil
}

func (n *Network) Info(ctx context.Context) *types.NetworkState {
	state := types.NetworkState{}

	state.Type = NetworkType
	state.CIDR = n.CIDR.String()
	state.IFace = types.NetworkInterface{
		Index: n.Device.GetIndex(),
		Name:  n.Device.GetName(),
		HAddr: n.Device.GetHardware(),
		Addr:  n.Device.GetAddr(),
	}
	state.Addr = n.ExtIface.IfaceAddr.String()
	state.IP = n.IntIface.IfaceAddr.String()
	return &state
}

func (n *Network) Destroy(ctx context.Context, network *types.NetworkState) error {

	return nil
}

func (n *Network) Create(ctx context.Context, network *types.SubnetManifest) (*types.NetworkState, error) {

	log.V(logLevel).Debugf("Connect to node to network: %v > %v", network.CIDR, network.IFace.Addr)

	if n.CIDR.String() == network.CIDR {
		log.V(logLevel).Debug("Skip local network provision")
		return n.Info(ctx), nil
	}

	// Parse MAC address from string
	lladdr, err := net.ParseMAC(network.IFace.HAddr)
	if err != nil {
		log.Errorf("Can-not parse MAC addres %v: %s", network.IFace.HAddr, err.Error())
		return nil, err
	}

	// Add ARP record
	log.V(logLevel).Debugf("Add new ARP record to %v :> %v", lladdr, network.Addr)
	if err := n.Device.AddARP(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
		log.Errorf("Can not add ARP record: %s", err.Error())
		return nil, err
	}

	// Add FDB record
	log.V(logLevel).Debugf("Add new FDB record to %v :> %v", network.IFace.HAddr, network.Addr)
	if err := n.Device.AddFDB(lladdr, net.ParseIP(network.Addr)); err != nil {
		log.Errorf("Can not add FDB record: %s", err.Error())
		if err := n.Device.DelARP(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
			return nil, err
		}
	}

	// Add route
	log.V(logLevel).Debugf("Add new route record for %v :> %v", network.CIDR, network.IFace.Addr)

	_, ipn, err := net.ParseCIDR(network.CIDR)
	if err != nil {
		log.Errorf("Can-not parse subnet %v: %s", network.CIDR, err.Error())
		return nil, err
	}

	vxlanRoute := netlink.Route{
		LinkIndex: n.Device.link.Attrs().Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Dst:       ipn,
		Gw:        net.ParseIP(network.IFace.Addr),
	}
	vxlanRoute.SetFlag(syscall.RTNH_F_ONLINK)

	if err := netlink.RouteReplace(&vxlanRoute); err != nil {
		log.Errorf("Add xvlan route err: %s", err.Error())
		log.V(logLevel).Debug("Clean up added before records")

		if err := n.Device.DelARP(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
			log.Errorf("Can not del ARP record: %s", err.Error())
			return nil, err
		}

		if err := n.Device.DelFDB(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
			log.Errorf("Can not del FBD record: %s", err.Error())
			return nil, err
		}

		return nil, err
	}

	state := types.NetworkState{}

	state.Type = NetworkType
	state.CIDR = n.CIDR.String()
	state.IFace = types.NetworkInterface{
		Index: n.Device.GetIndex(),
		Name:  n.Device.GetName(),
		HAddr: n.Device.GetHardware(),
		Addr:  n.Device.GetAddr(),
	}
	state.Addr = n.ExtIface.IfaceAddr.String()

	return &state, nil
}

func (n *Network) Replace(ctx context.Context, state *types.NetworkState, manifest *types.SubnetManifest) (*types.NetworkState, error) {

	if state != nil {
		if err := n.Destroy(ctx, state); err != nil {
			return nil, err
		}
	}

	if manifest == nil {
		return nil, nil
	}

	state, err := n.Create(ctx, manifest)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (n *Network) Subnets(ctx context.Context) (map[string]*types.NetworkState, error) {

	log.V(logLevel).Debug("Get current subnets list")

	var (
		subnets = make(map[string]*types.NetworkState)
		neighs  = make(map[string]string)
	)

	arps, err := netlink.NeighList(n.Device.link.Index, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("Can not get arps: %s", err.Error())
		return subnets, err
	}

	for _, arp := range arps {
		neighs[arp.IP.String()] = arp.HardwareAddr.String()
	}

	rules, err := utils.BridgeFDBList()
	if err != nil {
		log.Errorf("Can not FDB rules: %s", err.Error())
		return subnets, err
	}

	routes, err := netlink.RouteList(n.Device.link, netlink.FAMILY_V4)
	if err != nil {
		log.Errorf("Can not get routes: %s", err.Error())
		return subnets, err
	}

	for _, r := range routes {

		sn := types.NetworkState{}
		sn.Type = n.Device.link.Type()
		sn.CIDR = r.Dst.String()
		sn.IFace = types.NetworkInterface{
			Index: n.Device.link.Index,
			Name:  n.Device.link.Name,
			Addr:  r.Gw.String(),
			HAddr: neighs[r.Gw.String()],
		}

		for _, rule := range rules {
			if rule.Mac == sn.IFace.HAddr && rule.DST != "" {
				sn.Addr = rule.DST
			}
		}

		subnets[r.Dst.String()] = &sn
	}

	for r, sn := range subnets {
		log.V(logLevel).Debugf("SubnetSpec [%s]: %v", r, sn)
	}

	return subnets, nil
}
