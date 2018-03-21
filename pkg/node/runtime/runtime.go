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
	"github.com/lastbackend/lastbackend/pkg/node/runtime/network"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/events"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/pod"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/router"
	"github.com/lastbackend/lastbackend/pkg/node/runtime/volume"
)

type Runtime struct {
	ctx  context.Context
	spec chan *types.NodeSpec
}

func (r *Runtime) Restore() {
	log.Debug("node:runtime:restore:> restore init")
	network.Restore(r.ctx)
	router.Restore(r.ctx)
	volume.Restore(r.ctx)
	pod.Restore(r.ctx)
}

func (r *Runtime) Provision(ctx context.Context, spec *types.NodeSpec) error {
	log.Debug("node:runtime:provision:> provision init")
	for _, r := range spec.Routes {
		log.Debugf("route: %v", r)
	}

	for _, p := range spec.Pods {
		log.Debugf("pod: %v", p)
	}

	for _, v := range spec.Volumes {
		log.Debugf("volume: %v", v)
	}

	return nil
}

func (r *Runtime) Subscribe() {

	log.Debug("node:runtime:subscribe:> subscribe init")
	pc := make(chan *types.Pod)

	go func() {

		for {
			select {
			case p := <-pc:
				log.Debugf("node:runtime:subscribe:> new pod state event: %#v", p)
				events.NewPodStatusEvent(r.ctx, p)
			}
		}
	}()

	envs.Get().GetCri().Subscribe(r.ctx, envs.Get().GetState().Pods(), pc)
}

func (r *Runtime) GetSpec(ctx context.Context) error {

	log.Debug("node:runtime:getspec:> getspec request init")

	var (
		c = envs.Get().GetClient()
	)

	spec, err := c.GetSpec(ctx)
	if err != nil {
		log.Debug("node:runtime:getspec:> request err: %s", err.Error())
		return err
	}

	if spec == nil {
		log.Debug("node:runtime:getspec:> new spec is nil")
		return nil
	}

	//r.spec <- spec

	return nil
}

func (r *Runtime) Loop() {
	log.Debug("node:runtime:loop:> start runtime loop")

	go func(ctx context.Context) {
		for {
			select {
			case spec := <-r.spec:
				log.Debug("node:runtime:loop:> provision new spec")
				if err := r.Provision(ctx, spec); err != nil {
					log.Errorf("node:runtime:loop:> provision new spec err: %s", err.Error())
				}

			}

			log.Debug("node:runtime:loop:> request new spec")
			err := r.GetSpec(ctx)
			if err != nil {
				log.Debug("node:runtime:loop:> new spec request err: %s", err.Error())
			}
		}
	}(r.ctx)

	err := r.GetSpec(r.ctx)
	if err != nil {
		log.Debug("node:runtime:loop:> new spec request err: %s", err.Error())
	}
}

func NewRuntime(ctx context.Context) *Runtime {
	r := Runtime{
		ctx:  ctx,
		spec: make(chan *types.NodeSpec),
	}

	return &r
}
