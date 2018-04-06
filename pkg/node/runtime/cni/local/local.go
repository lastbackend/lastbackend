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

package local

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/cni"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"net"
)

type Network struct {
	cni.CNI

	ExtIface *NetworkInterface
	Network  *net.IPNet
	Subnet   *net.IPNet
	IP       net.IP
}

type NetworkInterface struct {
	Iface     *net.Interface
	IfaceAddr net.IP
}

func New() (*Network, error) {
	ip, _ := system.GetNodeIP()
	return &Network{
		ExtIface: &NetworkInterface{
			IfaceAddr: net.ParseIP(ip),
		},
	}, nil
}

func (n *Network) Info(ctx context.Context) *types.NetworkSpec {
	return &types.NetworkSpec{
		Type: "local",
		Addr: n.ExtIface.IfaceAddr.String(),
	}
}

func (n *Network) Create(ctx context.Context, network *types.NetworkSpec) error {
	return nil
}

func (n *Network) Destroy(ctx context.Context, network *types.NetworkSpec) error {
	return nil
}

func (n *Network) Replace(ctx context.Context, current *types.NetworkSpec, proposal *types.NetworkSpec) error {
	return nil
}

func (n *Network) Subnets(ctx context.Context) (map[string]*types.NetworkSpec, error) {
	return nil, nil
}
