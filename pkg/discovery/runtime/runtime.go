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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/discovery/runtime/endpoint"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logRuntimePrefix = "discovery:runtime"
	logLevel         = 3
)

type Runtime struct {
	spec chan *types.DiscoveryManifest
	ctx  context.Context
}

func (r *Runtime) Restore(ctx context.Context) {
	log.V(logLevel).Debugf("%s:restore:> restore init", logRuntimePrefix)
	var network = envs.Get().GetNet()

	if network != nil {
		if err := envs.Get().GetNet().SubnetRestore(ctx); err != nil {
			log.Errorf("%s:> can not restore network: %s", logRuntimePrefix, err.Error())
		}
	}

}

// Sync discovery runtime with new spec
func (r *Runtime) Sync(ctx context.Context, spec *types.DiscoveryManifest) error {
	log.V(logLevel).Debugf("%s:sync:> sync runtime state", logRuntimePrefix)
	r.spec <- spec
	return nil
}

func (r *Runtime) Loop(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:loop:> watch endpoint start", logRuntimePrefix)
	go endpoint.Watch(r.ctx)

	log.V(logLevel).Debugf("%s:loop:> start runtime loop", logRuntimePrefix)
	var network = envs.Get().GetNet()

	go func(ctx context.Context) {

		for {
			select {
			case spec := <-r.spec:

				log.V(logLevel).Debugf("%s:loop:> provision new spec", logRuntimePrefix)

				if spec.Meta.Initial && network != nil {
					log.V(logLevel).Debugf("%s> clean up networks", logRuntimePrefix)
					nets := network.Subnets().GetSubnets()

					for cidr := range nets {
						if _, ok := spec.Network[cidr]; !ok {
							network.SubnetDestroy(ctx, cidr)
						}
					}

				}

				log.V(logLevel).Debugf("%s> provision init", logRuntimePrefix)

				if network != nil {
					log.V(logLevel).Debugf("%s> provision networks", logRuntimePrefix)
					for cidr, n := range spec.Network {
						log.V(logLevel).Debugf("network: %v", n)
						if err := network.SubnetManage(ctx, cidr, n); err != nil {
							log.Errorf("Subnet [%s] create err: %s", n.CIDR, err.Error())
						}
					}
				}
			}
		}
	}(ctx)

	return nil
}

func NewRuntime(ctx context.Context) *Runtime {
	return &Runtime{
		ctx:  ctx,
		spec: make(chan *types.DiscoveryManifest),
	}
}
