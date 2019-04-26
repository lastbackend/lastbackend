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

package controller

import (
	"context"
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/discovery/runtime"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logPrefix   = "controller"
	logLevel    = 3
	ifaceDocker = "docker0"
)

type Controller struct {
	runtime *runtime.Runtime
	cache   struct {
		lock      sync.RWMutex
		resources types.DiscoveryStatus
	}
}

func New(r *runtime.Runtime) *Controller {
	var c = new(Controller)
	c.runtime = r
	return c
}

func (c *Controller) Connect(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:connect:> connect init", logPrefix)

	opts := v1.Request().Discovery().DiscoveryConnectOptions()
	opts.Info = envs.Get().GetState().Discovery().Info
	opts.Status = envs.Get().GetState().Discovery().Status

	for {
		err := envs.Get().GetClient().Connect(ctx, opts)
		if err == nil {
			log.Debugf("%s connected", logPrefix)
			return nil
		}

		log.Errorf("connect err: %s", err.Error())
		time.Sleep(3 * time.Second)
	}
}
