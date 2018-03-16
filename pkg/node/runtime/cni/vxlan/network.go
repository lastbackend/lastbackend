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

package vxlan

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni"
	"net"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni/utils"
	"github.com/vishvananda/netlink"
)

const NetworkType = "vxlan"
const DefaultContainerDevice = "docker0"

type Network struct {
	cni.CNI

	ExtIface *NetworkInterface
	Device   *Device
	Network  *net.IPNet
	Subnet   *net.IPNet
	IP       net.IP
}

type NetworkInterface struct {
	Iface     *net.Interface
	IfaceAddr net.IP
}

func New() (*Network, error) {

	var (
		nt  = new(Network)
		err error
	)

	nt.ExtIface = new(NetworkInterface)

	if nt.ExtIface.Iface, nt.ExtIface.IfaceAddr, err = utils.GetDefaultInterface(); err != nil {
		log.Errorf("Can not get default interface: %s", err.Error())
		return nt, err
	}

	nt.SetSubnetFromDevice(DefaultContainerDevice)

	if err := nt.AddInterface(); err != nil {
		log.Errorf("Can not add interface: %s", err.Error())
		return nt, err
	}

	// Add forward rules for network range
	go utils.SetupAndEnsureIPTables(utils.ForwardRules(nt.Network.String()), 5)

	return nt, nil
}

func (n *Network) SetSubnetFromDevice(name string) error {

	iface, err := utils.GetIfaceByName(name)
	if err != nil {
		log.Errorf("Can not find interface by name %s", name)
		return err
	}

	addrs, err := netlink.AddrList(&netlink.Device{
		LinkAttrs: netlink.LinkAttrs{
			Index: iface.Index,
		},
	}, syscall.AF_INET)

	if err != nil {
		log.Errorf("Can not locate docker interface ips: %s", err.Error())
	}

	if len(addrs) == 0 {
		log.Error("docker interface has not IP address")
		panic(0)
	}

	sip := make(net.IP, len(addrs[0].IPNet.IP))
	smk := make(net.IPMask, len(addrs[0].Mask))

	copy(sip, addrs[0].IPNet.IP)
	copy(smk, addrs[0].Mask)

	sip[3] = byte(0)
	n.Subnet = &net.IPNet{
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
		name:  fmt.Sprintf("%s.%d", DeviceDefaultName, DeviceDefaultVNI),
		index: n.ExtIface.Iface.Index,
		addr:  n.ExtIface.IfaceAddr,
		port:  DeviceDefaultPort,
	}); err != nil {
		log.Errorf("Can not create xvlan interface: %s", err.Error())
		return err
	}

	n.Device.SetIP(*n.Subnet)
	return nil
}

func (n *Network) Info(ctx context.Context) *types.Subnet {
	return &types.Subnet{
		Type:   NetworkType,
		Subnet: n.Subnet.String(),
		IFace: types.NetworkInterface{
			Index: n.Device.GetIndex(),
			Name:  n.Device.GetName(),
			HAddr: n.Device.GetHardware(),
			Addr:  n.Device.GetAddr(),
		},
		Addr: n.ExtIface.IfaceAddr.String(),
	}
}

func (n *Network) Destroy(ctx context.Context, network *types.Subnet) error {

	return nil
}

func (n *Network) Create(ctx context.Context, network *types.Subnet) error {
	log.Debugf("Connect to node to network: %v > %v", network.Subnet, network.IFace.Addr)

	if n.Subnet.String() == network.Subnet {
		log.Debug("Skip local network provision")
		return nil
	}

	// Parse MAC address from string
	lladdr, err := net.ParseMAC(network.IFace.HAddr)
	if err != nil {
		log.Errorf("Can-not parse MAC addres %v: %s", network.IFace.HAddr, err.Error())
		return err
	}

	// Add ARP record
	log.Debugf("Add new ARP record to %v :> %v", lladdr, network.Addr)
	if err := n.Device.AddARP(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
		log.Errorf("Can not add ARP record: %s", err.Error())
		return err
	}

	// Add FDB record
	log.Debugf("Add new FDB record to %v :> %v", network.IFace.HAddr, network.Addr)
	if err := n.Device.AddFDB(lladdr, net.ParseIP(network.Addr)); err != nil {
		log.Errorf("Can not add FDB record: %s", err.Error())
		return n.Device.DelARP(lladdr, net.ParseIP(network.IFace.Addr))
	}

	// Add route
	log.Debugf("Add new route record for %v :> %v", network.Subnet, network.IFace.Addr)

	_, ipn, err := net.ParseCIDR(network.Subnet)
	if err != nil {
		log.Errorf("Can-not parse subnet %v: %s", network.Subnet, err.Error())
		return err
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
		log.Debug("Clean up added before records")

		if err := n.Device.DelARP(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
			log.Errorf("Can not del ARP record: %s", err.Error())
			return err
		}

		if err := n.Device.DelFDB(lladdr, net.ParseIP(network.IFace.Addr)); err != nil {
			log.Errorf("Can not del FBD record: %s", err.Error())
			return err
		}

		return err
	}

	return nil
}

func (n *Network) Replace(ctx context.Context, current *types.Subnet, proposal *types.Subnet) error {

	if current != nil {
		if err := n.Destroy(ctx, current); err != nil {
			return err
		}
	}

	if proposal != nil {
		if err := n.Create(ctx, proposal); err != nil {
			return err
		}
	}

	return nil
}

func (n *Network) Subnets(ctx context.Context) (map[string]*types.Subnet, error) {

	log.Debug("Get current subnets list")

	var (
		subnets = make(map[string]*types.Subnet)
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

		sn := types.Subnet{
			Type:   n.Device.link.Type(),
			Subnet: r.Dst.String(),
			IFace: types.NetworkInterface{
				Index: n.Device.link.Index,
				Name:  n.Device.link.Name,
				Addr:  r.Gw.String(),
				HAddr: neighs[r.Gw.String()],
			},
		}

		for _, rule := range rules {
			if rule.Mac == sn.IFace.HAddr && rule.DST != "" {
				sn.Addr = rule.DST
			}
		}

		subnets[r.Dst.String()] = &sn
	}

	for r, sn := range subnets {
		log.Debugf("Subnet [%s]: %v", r, sn)
	}

	return subnets, nil
}
