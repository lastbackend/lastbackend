//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni/utils"
	"net"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/vishvananda/netlink"
)

type Device struct {
	link *netlink.Vxlan
	addr net.IP
}

type DeviceCreateOpts struct {
	vni   int
	name  string
	index int
	addr  net.IP
	port  int
}

const DeviceDefaultVNI = 1
const DeviceDefaultName = "lb"
const DeviceDefaultPort = 8472

func NewDevice(opts DeviceCreateOpts) (*Device, error) {

	d := new(Device)

	link := netlink.Vxlan{
		LinkAttrs: netlink.LinkAttrs{
			Name: opts.name,
		},
		VxlanId:      opts.vni,
		VtepDevIndex: opts.index,
		SrcAddr:      opts.addr,
		Port:         opts.port,
		Learning:     true,
		GBP:          true,
	}

	d.link = &link

	if err := d.Create(); err != nil {
		return d, err
	}

	if err := d.SetUp(); err != nil {
		return d, err
	}

	return d, nil
}

func (d *Device) Create() error {

	log.Debug("Create new vxlan interface")

	err := netlink.LinkAdd(d.link)
	if err == syscall.EEXIST {
		log.Debugf("Device already exists: %s", d.link.Name)

		l, err := netlink.LinkByName(d.link.Name)
		if err != nil {
			log.Debugf("Link by name: %s", err.Error())
		}

		d.link = l.(*netlink.Vxlan)
		return nil
	}

	if err != nil {
		return err
	}

	link, err := netlink.LinkByIndex(d.link.Index)
	if err != nil {
		return fmt.Errorf("can't locate created vxlan device with index %v", d.link.Index)
	}

	l, ok := link.(*netlink.Vxlan)
	if !ok {
		return fmt.Errorf("created vxlan device with index %v is not vxlan", d.link.Index)
	}

	d.link = l
	return nil
}

func (d *Device) SetIP(nt net.IPNet) error {

	log.Debug("Set IP for device")

	ip := make(net.IP, len(nt.IP))
	copy(ip, nt.IP)

	ipn := net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(32, 32),
	}


	ipn.IP[3] = byte(0)
	addr := netlink.Addr{IPNet: &ipn, Broadcast: net.ParseIP("0.0.0.0")}

	existingAddrs, err := netlink.AddrList(d.link, netlink.FAMILY_V4)
	if err != nil {
		return err
	}

	// flannel will never make this happen. This situation can only be caused by a user, so get them to sort it out.
	if len(existingAddrs) > 1 {
		return fmt.Errorf("link has incompatible addresses. Remove additional addresses and try again. %#v", d.link)
	}

	// If the device has an incompatible address then delete it. This can happen if the lease changes for example.
	if len(existingAddrs) == 1 && !existingAddrs[0].Equal(addr) {
		if err := netlink.AddrDel(d.link, &existingAddrs[0]); err != nil {
			return fmt.Errorf("failed to remove IP address %s from %s: %s", ipn.String(), d.link.Attrs().Name, err)
		}
		existingAddrs = []netlink.Addr{}
	}

	// Actually add the desired address to the interface if needed.
	if len(existingAddrs) == 0 {
		if err := netlink.AddrAdd(d.link, &addr); err != nil {
			return fmt.Errorf("failed to add IP address %s to %s: %s", ipn.String(), d.link.Attrs().Name, err)
		}
	}

	d.addr = addr.IP
	return nil
}

func (d *Device) SetUp() error {
	log.Debug("Set vxlan interface up")
	if err := netlink.LinkSetUp(d.link); err != nil {
		return fmt.Errorf("failed to set interface %s to UP state: %s", d.link.Attrs().Name, err)
	}

	return nil
}

func (d *Device) AddFDB(MAC net.HardwareAddr, IP net.IP) error {
	log.Debugf("Add FDB: %v, %v", MAC, IP)

	rules, err := utils.BridgeFDBList()
	if err != nil {
		log.Errorf("Can not FDB rules: %s", err.Error())
		return err
	}

	for _, rule := range rules {
		if rule.DST == IP.String() && rule.Device == d.link.Name {
			mac, err :=net.ParseMAC(rule.Mac)
			if err != nil {
				continue
			}
			d.DelFDB(mac, IP)
		}
	}

	return netlink.NeighSet(&netlink.Neigh{
		LinkIndex:    d.link.Index,
		State:        netlink.NUD_PERMANENT,
		Family:       syscall.AF_BRIDGE,
		Flags:        netlink.NTF_SELF,
		IP:           IP,
		HardwareAddr: MAC,
	})
}

func (d *Device) DelFDB(MAC net.HardwareAddr, IP net.IP) error {
	log.Debugf("Del FDB: %v, %v", MAC, IP)
	return netlink.NeighDel(&netlink.Neigh{
		LinkIndex:    d.link.Index,
		Family:       syscall.AF_BRIDGE,
		Flags:        netlink.NTF_SELF,
		IP:           IP,
		HardwareAddr: MAC,
	})
}

func (d *Device) AddARP(MAC net.HardwareAddr, IP net.IP) error {
	log.Debugf("Add ARP: %v, %v", MAC, IP)
	return netlink.NeighSet(&netlink.Neigh{
		LinkIndex:    d.link.Index,
		State:        netlink.NUD_PERMANENT,
		Type:         syscall.RTN_UNICAST,
		IP:           IP,
		HardwareAddr: MAC,
	})
}

func (d *Device) DelARP(MAC net.HardwareAddr, IP net.IP) error {
	log.Debugf("Del ARP: %v, %v", MAC, IP)
	return netlink.NeighDel(&netlink.Neigh{
		LinkIndex:    d.link.Index,
		State:        netlink.NUD_PERMANENT,
		Type:         syscall.RTN_UNICAST,
		IP:           IP,
		HardwareAddr: MAC,
	})
	return nil
}

func (d *Device) GetIndex() int {
	return d.link.Index
}

func (d *Device) GetHardware() string {
	return d.link.HardwareAddr.String()
}

func (d *Device) GetName() string {
	return d.link.Name
}

func (d *Device) GetAddr() string {
	return d.addr.String()
}