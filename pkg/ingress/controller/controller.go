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

package controller

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/ingress/envs"
	"github.com/lastbackend/lastbackend/pkg/ingress/runtime"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
	"time"
)

const (
	logPrefix = "controller:>"
	logLevel  = 3
)

type Controller struct {
	runtime *runtime.Runtime
	cache   struct {
		lock      sync.RWMutex
		resources types.IngressStatus
	}
}

func New(r *runtime.Runtime) *Controller {
	var c = new(Controller)
	c.runtime = r
	return c
}

func (c *Controller) Connect(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:connect:> connect init", logPrefix)

	opts := v1.Request().Ingress().IngressConnectOptions()
	opts.Info = envs.Get().GetState().Ingress().Info
	opts.Status = envs.Get().GetState().Ingress().Status
	opts.Network = *envs.Get().GetNet().Info(ctx)

	for {
		err := envs.Get().GetClient().Connect(ctx, opts)
		if err == nil {
			log.Debugf("%s connected", logPrefix)
			return nil
		}

		log.Errorf("connect err: %s", err.Error())
		time.Sleep(3 * time.Second)
	}

	return nil
}

func (c *Controller) Sync(ctx context.Context) error {

	log.Debugf("Start ingress sync")

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		opts := new(request.IngressStatusOptions)

		spec, err := envs.Get().GetClient().SetStatus(ctx, opts)
		if err != nil {
			log.Errorf("ingress:exporter:dispatch err: %s", err.Error())
		}

		if spec != nil {
			c.runtime.Sync(ctx, spec.Decode())
		} else {
			log.Debug("received spec is nil, skip apply changes")
		}
	}

	return nil
}
