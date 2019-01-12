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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/discovery/envs"
	"github.com/lastbackend/lastbackend/pkg/discovery/runtime"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni/utils"
	"github.com/spf13/viper"
	"net"
	"sync"
	"time"
)

const (
	logPrefix   = "controller:>"
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

	if envs.Get().GetNet() != nil {
		opts.Network = *envs.Get().GetNet().Info(ctx)
	}

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

	log.Debugf("Start discovery sync")

	var (
		err    error
		ip     = net.IP{}
		iiface = viper.GetString("runtime.cpi.interface.internal")
	)

	if iiface == types.EmptyString {
		iiface = ifaceDocker
	}

	log.Debugf("find interface to traffic route by name: %s", iiface)
	_, ip, err = utils.GetIfaceByName(iiface)
	if err != nil {
		return err
	}

	log.Debugf("internal route ip net: %s", ip.String())

	log.V(logLevel).Debugf("%s:loop:> update current discovery service info", logPrefix)
	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		opts := new(request.DiscoveryStatusOptions)

		status := envs.Get().GetState().Discovery().Status
		opts.IP = ip.String()
		opts.Port = status.Port
		opts.Ready = status.Ready

		spec, err := envs.Get().GetClient().SetStatus(ctx, opts)
		if err != nil {
			log.Errorf("discovery:exporter:dispatch err: %s", err.Error())
		}

		if spec != nil {
			c.runtime.Sync(ctx, spec.Decode())
		} else {
			log.Debug("received spec is nil, skip apply changes")
		}

	}

	return nil
}
