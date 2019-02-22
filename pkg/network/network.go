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

package network

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/network/state"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni"
	ni "github.com/lastbackend/lastbackend/pkg/runtime/cni/cni"
	"github.com/lastbackend/lastbackend/pkg/runtime/cpi"
	pi "github.com/lastbackend/lastbackend/pkg/runtime/cpi/cpi"
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

func New() (*Network, error) {

	var err error

	net := new(Network)

	if viper.GetString("runtime.cni.type") == types.EmptyString &&
		viper.GetString("runtime.cpi.type") == types.EmptyString {
		log.Debug("run without network management")
		return nil, nil
	}

	net.state = state.New()
	if net.cni, err = ni.New(); err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
		return nil, err
	}

	if net.cpi, err = pi.New(); err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
		return nil, err
	}

	rip := viper.GetString("network.resolver.ip")
	if rip == types.EmptyString {
		rip = DefaultResolverIP
	}

	net.resolver.ip = rip
	net.resolver.external = viper.GetStringSlice("network.resolver.external")
	if len(net.resolver.external) == 0 {
		net.resolver.external = []string{"8.8.8.8", "8.8.4.4"}
	}

	return net, nil
}
