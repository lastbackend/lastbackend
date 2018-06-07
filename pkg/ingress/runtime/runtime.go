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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/ingress/envs"
	"github.com/lastbackend/lastbackend/pkg/ingress/events"
	"github.com/lastbackend/lastbackend/pkg/log"
	"time"
)

type Runtime struct {
	ctx  context.Context
	spec chan *types.IngressSpec
}

func (r *Runtime) Provision(ctx context.Context, spec *types.IngressSpec, clean bool) error {

	var (
		msg = "node:runtime:provision:"
	)

	log.Debugf("%s> provision init", msg)

	if clean {
		log.Debugf("%s> clean up current routes", msg)
	}

	log.Debugf("%s> provision routes", msg)
	for _, r := range spec.Routes {
		log.Debugf("route: %v", r)
	}

	return nil
}

func (r *Runtime) Connect(ctx context.Context) error {

	log.Debug("ingress:runtime:connect:> connect init")
	if err := events.NewConnectEvent(ctx); err != nil {
		log.Errorf("ingress:runtime:connect:> connect err: %s", err.Error())
		return err
	}

	return nil
}

func (r *Runtime) GetSpec(ctx context.Context) error {

	log.Debug("ingress:runtime:getspec:> getspec request init")

	var (
		c = envs.Get().GetClient()
	)

	spec, err := c.GetSpec(ctx)
	if err != nil {
		log.Errorf("ingress:runtime:getspec:> request err: %s", err.Error())
		return err
	}

	if spec == nil {
		log.Warnf("ingress:runtime:getspec:> new spec is nil")
		return nil
	}

	r.spec <- spec.Decode()
	return nil
}

func (r *Runtime) Loop() {
	log.Debug("ingress:runtime:loop:> start runtime loop")

	var clean = true

	go func(ctx context.Context) {
		for {
			select {
			case spec := <-r.spec:
				log.Debug("ingress:runtime:loop:> provision new spec")
				if err := r.Provision(ctx, spec, clean); err != nil {
					log.Errorf("ingress:runtime:loop:> provision new spec err: %s", err.Error())
				}
				clean = false
			}
		}
	}(r.ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			err := r.GetSpec(r.ctx)
			if err != nil {
				log.Debugf("ingress:runtime:loop:> new spec request err: %s", err.Error())
			}
		}
	}(context.Background())

	err := r.GetSpec(r.ctx)
	if err != nil {
		log.Debugf("ingress:runtime:loop:> new spec request err: %s", err.Error())
	}
}

func NewRuntime(ctx context.Context) *Runtime {
	r := Runtime{
		ctx:  ctx,
		spec: make(chan *types.IngressSpec),
	}

	return &r
}
