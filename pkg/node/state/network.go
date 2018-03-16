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

package state

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func (n *NetworkState) GetSubnets() map[string]types.Subnet {
	return n.subnets
}

func (n *NetworkState) AddSubnet(sn *types.Subnet) {
	log.V(logLevel).Debugf("State: NetworkState: add subnet: %v", sn)
	n.SetSubnet(sn)
}

func (n *NetworkState) SetSubnet(sn *types.Subnet) {
	log.V(logLevel).Debugf("State: NetworkState: set subnet: %v", sn)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.subnets[sn.Subnet]; ok {
		delete(n.subnets, sn.Subnet)
	}

	n.subnets[sn.Subnet] = *sn
}

func (n *NetworkState) GetSubnet(sn string) *types.Subnet {
	log.V(logLevel).Debugf("State: NetworkState: get subnet: %s", sn)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.subnets[sn]
	if !ok {
		return nil
	}
	return &s
}

func (n *NetworkState) DelSubnet(sn *types.Subnet) {
	log.V(logLevel).Debugf("State: NetworkState: del subnet: %v", sn)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.subnets[sn.Subnet]; ok {
		delete(n.subnets, sn.Subnet)
	}
}
