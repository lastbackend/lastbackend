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

package utils

import (
	"fmt"
	"net"
	"syscall"

	"github.com/lastbackend/lastbackend/tools/log"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/vishvananda/netlink"
)

func GetIfaceByName(name string) (*net.Interface, net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, iface := range ifaces {
		if iface.Name == name {
			ifaceAddr, err := GetIfaceIP4Addr(&iface)
			if err != nil {
				return nil, nil, errors.New(fmt.Sprintf("failed to find IPv4 address for interface %s", iface.Name))
			}
			return &iface, ifaceAddr, nil
		}
	}

	return nil, nil, nil

}

func GetDefaultInterface() (*net.Interface, net.IP, error) {
	var iface *net.Interface
	var ifaceAddr net.IP
	var err error

	log.Info("Determining IP address of default interface")
	if iface, err = GetDefaultGatewayIface(); err != nil {
		return nil, nil, errors.New(fmt.Sprintf("failed to get default interface: %s", err))
	}

	if iface == nil {
		return nil, nil, errors.New(fmt.Sprintf("failed to get default interface"))
	}

	if ifaceAddr == nil {
		ifaceAddr, err = GetIfaceIP4Addr(iface)
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("failed to find IPv4 address for interface %s", iface.Name))
		}
	}

	return iface, ifaceAddr, nil

}

func GetDefaultGatewayIface() (*net.Interface, error) {
	routes, err := netlink.RouteList(nil, syscall.AF_INET)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.Dst == nil || route.Dst.String() == "0.0.0.0/0" {
			if route.LinkIndex <= 0 {
				return nil, errors.New("Found default route but could not determine interface")
			}
			return net.InterfaceByIndex(route.LinkIndex)
		}
	}

	return nil, errors.New("Unable to find default route")
}

func GetIfaceIP4Addr(iface *net.Interface) (net.IP, error) {
	addrs, err := getIfaceAddrs(iface)
	if err != nil {
		return nil, err
	}

	// prefer non link-local addr
	var ll net.IP

	for _, addr := range addrs {
		if addr.IP.To4() == nil {
			continue
		}

		if addr.IP.IsGlobalUnicast() {
			return addr.IP, nil
		}

		if addr.IP.IsLinkLocalUnicast() {
			ll = addr.IP
		}
	}

	if ll != nil {
		// didn't find global but found link-local. it'll do.
		return ll, nil
	}

	return nil, errors.New("No IPv4 address found for given interface")
}

func getIfaceAddrs(iface *net.Interface) ([]netlink.Addr, error) {
	link := &netlink.Device{
		LinkAttrs: netlink.LinkAttrs{
			Index: iface.Index,
		},
	}

	return netlink.AddrList(link, syscall.AF_INET)
}
