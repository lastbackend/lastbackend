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
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/network/state"
)

const (
	logResolverPrefix   = "network:resolver"
	resolverEndpointKey = "resolver"
)

func (n *Network) Resolvers() *state.ResolverState {
	return n.state.Resolvers()
}

func (n *Network) GetResolverIP() string {
	return n.resolver.ip
}

func (n *Network) GetResolverEndpointKey() string {
	return resolverEndpointKey
}

func (n *Network) GetExternalDNS() []string {
	return n.resolver.external
}

func (n *Network) ResolverManage(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:> create resolver", logResolverPrefix)

	manifest := new(types.EndpointManifest)
	manifest.IP = n.resolver.ip
	manifest.PortMap = make(map[uint16]string)

	resolvers := n.state.Resolvers().GetResolvers()

	if len(resolvers) == 0 {
		manifest.Upstreams = n.resolver.external
		manifest.PortMap[53] = "53/udp"
	} else {
		var port uint16
		for _, r := range resolvers {
			if port == 0 {
				manifest.Upstreams = append(manifest.Upstreams, r.IP)
				port = r.Port
			} else if port == r.Port {
				manifest.Upstreams = append(manifest.Upstreams, r.IP)
			} else {
				continue
			}
		}
		if port == 0 {
			return errors.New("can not create endpoint: reason: resolver port can not be 0")
		}
		manifest.PortMap[53] = fmt.Sprintf("%d/udp", port)
	}

	if err := n.EndpointManage(ctx, resolverEndpointKey, manifest); err != nil {
		log.Errorf("%s:> can not create endpoint", logResolverPrefix)
		return err
	}

	return nil
}
