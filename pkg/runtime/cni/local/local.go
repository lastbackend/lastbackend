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

package local

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"net"
)

const NetworkType = "local"
const DefaultContainerDevice = "docker0"
const localIP = "127.0.0.1"

type Network struct {
	cni.CNI

	ExtIface *NetworkInterface
	Network  *net.IPNet
	CIDR     *net.IPNet
	IP       net.IP
}

type NetworkInterface struct {
	Iface     *net.Interface
	IfaceAddr net.IP
}

func New() (*Network, error) {
	ip, _ := system.GetHostIP(types.EmptyString)

	iface := getInterface()

	nt := &Network{
		ExtIface: &NetworkInterface{
			IfaceAddr: net.ParseIP(ip),
			Iface:     iface,
		},
	}

	if iface != nil {

		nt.CIDR = &net.IPNet{
			IP:   net.ParseIP(localIP),
			Mask: net.ParseIP(localIP).DefaultMask(),
		}

		nt.Network = &net.IPNet{
			IP:   net.ParseIP(localIP),
			Mask: net.CIDRMask(8, 32),
		}
	}

	return nt, nil
}

func (n *Network) Info(ctx context.Context) *types.NetworkState {

	state := types.NetworkState{}
	state.Type = NetworkType
	state.Addr = n.ExtIface.IfaceAddr.String()
	if n.CIDR != nil {
		state.CIDR = n.CIDR.String()
	}

	if n.ExtIface.Iface != nil {
		state.IFace = types.NetworkInterface{
			Index: n.ExtIface.Iface.Index,
			Name:  n.ExtIface.Iface.Name,
			HAddr: n.ExtIface.Iface.HardwareAddr.String(),
			Addr:  net.ParseIP(localIP).String(),
		}
	}

	return &state
}

func (n *Network) Create(ctx context.Context, network *types.SubnetManifest) (*types.NetworkState, error) {
	return n.Info(ctx), nil
}

func (n *Network) Destroy(ctx context.Context, network *types.NetworkState) error {
	return nil
}

func (n *Network) Replace(ctx context.Context, state *types.NetworkState, manifest *types.SubnetManifest) (*types.NetworkState, error) {
	return n.Info(ctx), nil
}

func (n *Network) Subnets(ctx context.Context) (map[string]*types.NetworkState, error) {
	return nil, nil
}
