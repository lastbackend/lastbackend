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

package state

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"
)

const logSubnetPrefix = "state:subnet:>"

type SubnetState struct {
	lock    sync.RWMutex
	subnets map[string]models.NetworkState
}

func (n *SubnetState) GetSubnets() map[string]models.NetworkState {
	return n.subnets
}

func (n *SubnetState) AddSubnet(cidr string, sn *models.NetworkState) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s add subnet: %s", logSubnetPrefix, cidr)
	n.SetSubnet(cidr, sn)
}

func (n *SubnetState) SetSubnet(cidr string, sn *models.NetworkState) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s set subnet: %s", logSubnetPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}

	n.subnets[cidr] = *sn
}

func (n *SubnetState) GetSubnet(cidr string) *models.NetworkState {
	log := logger.WithContext(context.Background())
	log.Debugf("%s get subnet: %s", logSubnetPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.subnets[cidr]
	if !ok {
		return nil
	}
	return &s
}

func (n *SubnetState) DelSubnet(cidr string) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s del subnet: %s", logSubnetPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}
}
