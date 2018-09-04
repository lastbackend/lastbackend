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

package runtime

import (
	"context"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
)

const (
	logNodeRuntimePrefix = "node:runtime"
	logLevel             = 3
)

type Runtime struct {
	ctx  context.Context
	spec chan *types.NodeManifest
}

func (r *Runtime) Restore() {
	log.V(logLevel).Debugf("%s:restore:> restore init", logNodeRuntimePrefix)

	NetworkRestore(r.ctx)
	VolumeRestore(r.ctx)
	PodRestore(r.ctx)
	EndpointRestore(r.ctx)

	if envs.Get().GetModeIngress() {
		envs.Get().GetIngress().Provision(context.Background())
	}
}

func (r *Runtime) Provision(ctx context.Context, spec *types.NodeManifest) error {

	log.V(logLevel).Debugf("%s> provision init", logNodeRuntimePrefix)

	log.V(logLevel).Debugf("%s> provision networks", logNodeRuntimePrefix)
	for cidr, n := range spec.Network {
		log.V(logLevel).Debugf("network: %v", n)
		if err := NetworkManage(ctx, cidr, n); err != nil {
			log.Errorf("Subnet [%s] create err: %s", n.CIDR, err.Error())
		}
	}

	log.V(logLevel).Debugf("%s> update secrets", logNodeRuntimePrefix)
	for s, spec := range spec.Secrets {
		log.V(logLevel).Debugf("secret: %s > %s", s, spec.State)
	}

	log.V(logLevel).Debugf("%s> provision pods", logNodeRuntimePrefix)
	for p, spec := range spec.Pods {
		log.V(logLevel).Debugf("pod: %v", p)
		if err := PodManage(ctx, p, spec); err != nil {
			log.Errorf("Pod [%s] manage err: %s", p, err.Error())
		}
	}

	log.V(logLevel).Debugf("%s> provision endpoints", logNodeRuntimePrefix)
	for e, spec := range spec.Endpoints {
		log.V(logLevel).Debugf("endpoint: %v", e)
		if err := EndpointManage(ctx, e, spec); err != nil {
			log.Errorf("Endpoint [%s] manage err: %s", e, err.Error())
		}
	}

	if envs.Get().GetModeIngress() {
		log.V(logLevel).Debugf("%s> provision routes", logNodeRuntimePrefix)
		for e, spec := range spec.Routes {
			log.V(logLevel).Debugf("route: %v", e)
			if err := RouteManage(ctx, e, spec); err != nil {
				log.Errorf("Route [%s] manage err: %s", e, err.Error())
			}
		}

	}


	log.V(logLevel).Debugf("%s> provision volumes", logNodeRuntimePrefix)
	for _, v := range spec.Volumes {
		log.V(logLevel).Debugf("volume: %v", v)
	}

	return nil
}

func (r *Runtime) Subscribe() {

	log.V(logLevel).Debugf("%s:subscribe:> subscribe init", logNodeRuntimePrefix)

	if err := containerSubscribe(r.ctx); err != nil {
		log.Errorf("container subscribe err: %v", err)
	}
}

func (r *Runtime) Connect(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:connect:> connect init", logNodeRuntimePrefix)
	if err := events.NewConnectEvent(ctx); err != nil {
		log.Errorf("%s:connect:> connect err: %s", logNodeRuntimePrefix, err.Error())
		return err
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			if err := events.NewStatusEvent(ctx); err != nil {
				log.Errorf("%s:connect:> send status err: %s", logNodeRuntimePrefix, err.Error())
			}
		}
	}(ctx)

	return nil
}

func (r *Runtime) GetSpec(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:getspec:> getspec request init", logNodeRuntimePrefix)

	var (
		c = envs.Get().GetNodeClient()
	)

	spec, err := c.GetSpec(ctx)
	if err != nil {
		log.Errorf("%s:getspec:> request err: %s", logNodeRuntimePrefix, err.Error())
		return err
	}

	if spec == nil {
		log.Warnf("%s:getspec:> new spec is nil", logNodeRuntimePrefix)
		return nil
	}

	r.spec <- spec.Decode()
	return nil
}

func (r *Runtime) Clean(ctx context.Context, manifest *types.NodeManifest) error {

	log.V(logLevel).Debugf("%s> clean up endpoints", logNodeRuntimePrefix)
	endpoints := envs.Get().GetState().Endpoints().GetEndpoints()
	for e := range endpoints {
		if _, ok := manifest.Endpoints[e]; !ok {
			EndpointDestroy(context.Background(), e, endpoints[e])
		}
	}

	log.V(logLevel).Debugf("%s> clean up pods", logNodeRuntimePrefix)
	pods := envs.Get().GetState().Pods().GetPods()

	for k := range pods {
		if _, ok := manifest.Pods[k]; !ok {
			if !envs.Get().GetState().Pods().IsLocal(k) {
				PodDestroy(context.Background(), k, pods[k])
			}
		}
	}

	log.V(logLevel).Debugf("%s> clean up networks", logNodeRuntimePrefix)
	nets := envs.Get().GetState().Networks().GetSubnets()

	for cidr := range nets {
		if _, ok := manifest.Network[cidr]; !ok {
			NetworkDestroy(ctx, cidr)
		}
	}

	return nil
}

func (r *Runtime) Loop() {
	log.V(logLevel).Debugf("%s:loop:> start runtime loop", logNodeRuntimePrefix)

	var clean = true

	go func(ctx context.Context) {
		for {
			select {
			case spec := <-r.spec:
				log.V(logLevel).Debugf("%s:loop:> provision new spec", logNodeRuntimePrefix)

				if clean {
					if err := r.Clean(ctx, spec); err != nil {
						log.Errorf("%s:loop:> clean err: %s", logEndpointPrefix, err.Error())
						continue
					}
					clean = false
				}

				if err := r.Provision(ctx, spec); err != nil {
					log.Errorf("%s:loop:> provision new spec err: %s", logNodeRuntimePrefix, err.Error())
				}
			}
		}
	}(r.ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			err := r.GetSpec(r.ctx)
			if err != nil {
				log.V(logLevel).Debugf("%s:loop:> new spec request err: %s", logNodeRuntimePrefix, err.Error())
			}
		}
	}(context.Background())

	err := r.GetSpec(r.ctx)
	if err != nil {
		log.V(logLevel).Debugf("%s:loop:> new spec request err: %s", logNodeRuntimePrefix, err.Error())
	}
}

func (r *Runtime) Ingress() {

}

func NewRuntime(ctx context.Context) *Runtime {
	r := Runtime{
		ctx:  ctx,
		spec: make(chan *types.NodeManifest),
	}

	return &r
}
