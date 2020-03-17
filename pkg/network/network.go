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

package network

import (
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/pkg/network/state"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni"
	"github.com/lastbackend/lastbackend/pkg/runtime/cpi"
	"github.com/lastbackend/lastbackend/tools/log"
	"github.com/spf13/viper"
)

const (
	DefaultResolverIP = "172.17.0.1"
)

type Network struct {
	state    *state.State
	cni      cni.CNI
	cpi      cpi.CPI
	resolver struct {
		ip       string
		external []string
	}
}

func New(v *viper.Viper) (*Network, error) {

	var err error

	net := new(Network)

	if v.GetString("network.cni.type") == types.EmptyString &&
		v.GetString("network.cpi.type") == types.EmptyString {
		log.Debug("run without network management")
		return nil, nil
	}

	net.state = state.New()
	if net.cni, err = cni.New(v); err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
		return nil, err
	}

	if net.cpi, err = cpi.New(v); err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
		return nil, err
	}

	rip := v.GetString("network.resolver.ip")
	if rip == types.EmptyString {
		rip = DefaultResolverIP
	}

	net.resolver.ip = rip
	net.resolver.external = v.GetStringSlice("network.resolver.servers")
	if len(net.resolver.external) == 0 {
		net.resolver.external = []string{"8.8.8.8", "8.8.4.4"}
	}

	return net, nil
}
