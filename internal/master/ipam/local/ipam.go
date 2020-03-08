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

package local

import (
	"context"
	"net"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logIPAMPrefix         = "controller:ipam:>"
	IPAMLeaseNotAvailable = "IPAMLeaseNotAvailable"
	defaultCIDR           = "172.17.0.0/16"
)

// IPAM - IP address management
type IPAM struct {
	leased    map[string]bool
	released  map[string]bool
	available int
	reserved  int
	storage   storage.Storage
}

// Lease IP from range
func (i *IPAM) Lease() (*net.IP, error) {

	var (
		lease string
	)

	// Check available leases count
	if i.available == 0 {
		return nil, errors.New(IPAMLeaseNotAvailable)
	}

	// Get first map element
	for i := range i.released {
		lease = i
		break
	}

	// Mark lease IP as leased
	delete(i.released, lease)
	i.available--
	// Find new lease for next reservation

	// Decrease available count
	i.leased[lease] = true
	i.reserved++

	ip := net.ParseIP(lease)

	if err := i.save(); err != nil {
		return nil, err
	}

	return &ip, nil
}

// release IP
func (i *IPAM) Release(ip *net.IP) error {

	var (
		lease = ip.String()
	)

	// Mark IP as released
	delete(i.leased, lease)
	i.reserved--
	// Add IP as released and increase available count
	i.released[lease] = true
	i.available++

	return i.save()
}

// Available ips count
func (i *IPAM) Available() int {
	return i.available
}

// Reserved ips count
func (i *IPAM) Reserved() int {
	return i.reserved
}

func (i *IPAM) save() error {

	var (
		ips = make([]string, 0)
	)

	for ip := range i.leased {
		ips = append(ips, ip)
	}

	opts := storage.GetOpts()
	opts.Force = true
	return i.storage.Set(context.Background(), i.storage.Collection().System(), "ipam", &ips, opts)
}

// New IPAM object initializing and returning
func New(stg storage.Storage, cidr string) (*IPAM, error) {

	var (
		skip = true
		ipam = new(IPAM)
	)

	ipam.storage = stg
	ipam.leased = make(map[string]bool, 0)
	ipam.released = make(map[string]bool, 0)

	if cidr == "" {
		cidr = defaultCIDR
	}

	// Get IP range by network CIDR
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		if skip {
			// remove network address and broadcast address
			skip = false
			continue
		}

		ipam.released[ip.String()] = true
		ipam.available++
	}

	ips := make([]string, 0)

	// Get IP list from database storage
	err = stg.Get(context.Background(), stg.Collection().System(), "ipam", &ips, nil)
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s get context error: %s", logIPAMPrefix, err.Error())
			return nil, err
		}
	}

	// Mark IPs as leased
	for _, item := range ips {
		if _, ok := ipam.released[item]; ok {
			delete(ipam.released, item)
			ipam.available--
			ipam.leased[item] = true
			ipam.reserved++
		}
	}

	return ipam, nil
}
