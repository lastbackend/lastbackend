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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/node/runtime"
	"github.com/spf13/viper"
	"io/ioutil"
	"sync"
	"time"
)

const (
	logPrefix = "client:>"
	logLevel  = 3
)

type Controller struct {
	runtime *runtime.Runtime
	cache   struct {
		lock      sync.RWMutex
		resources types.NodeStatus
		pods      map[string]*types.PodStatus
		volumes   map[string]*types.VolumeStatus
	}
}

func New(r *runtime.Runtime) *Controller {
	var c = new(Controller)
	c.runtime = r
	c.cache.pods = make(map[string]*types.PodStatus)
	c.cache.volumes = make(map[string]*types.VolumeStatus)

	for p, st := range envs.Get().GetState().Pods().GetPods() {
		c.cache.pods[p] = st
	}
	return c
}

func (c *Controller) Connect(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:connect:> connect init", logPrefix)

	opts := v1.Request().Node().NodeConnectOptions()
	opts.Info = envs.Get().GetState().Node().Info

	log.Info(opts.Info)

	opts.Status = envs.Get().GetState().Node().Status
	var network = envs.Get().GetNet()
	if network != nil {
		opts.Network = *envs.Get().GetNet().Info(ctx)
	}

	if viper.IsSet("node.tls") {
		opts.TLS = !viper.GetBool("node.tls.insecure")

		if opts.TLS {
			caData, err := ioutil.ReadFile(viper.GetString("node.tls.ca"))
			if err != nil {
				log.Errorf("%s:connect_event:> read ca cert file err: %v", logPrefix, err)
				return err
			}

			certData, err := ioutil.ReadFile(viper.GetString("node.tls.client_cert"))
			if err != nil {
				log.Errorf("%s:connect_event:> read client cert file err: %v", logPrefix, err)
				return err
			}

			keyData, err := ioutil.ReadFile(viper.GetString("node.tls.client_key"))
			if err != nil {
				log.Errorf("%s:connect_event:> read client key file err: %v", logPrefix, err)
				return err
			}

			opts.SSL = new(request.SSL)
			opts.SSL.CA = caData
			opts.SSL.Key = keyData
			opts.SSL.Cert = certData
		}
	}

	for {
		if err := envs.Get().GetNodeClient().Connect(ctx, opts); err == nil {
			return nil
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *Controller) Sync(ctx context.Context) error {

	log.Debugf("%s start node sync", logPrefix)

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		opts := new(request.NodeStatusOptions)
		opts.Pods = make(map[string]*request.NodePodStatusOptions)
		opts.Volumes = make(map[string]*request.NodeVolumeStatusOptions)

		opts.Resources.Capacity = envs.Get().GetState().Node().Status.Capacity
		opts.Resources.Allocated = envs.Get().GetState().Node().Status.Allocated

		c.cache.lock.Lock()
		var i = 0
		for p, status := range c.cache.pods {
			i++
			if i > 10 {
				break
			}

			if !envs.Get().GetState().Pods().IsLocal(p) && status != nil {
				opts.Pods[p] = getPodOptions(status)
			} else {
				delete(c.cache.pods, p)
			}
		}

		var iv = 0
		for v, status := range c.cache.volumes {
			iv++
			if iv > 10 {
				break
			}

			if !envs.Get().GetState().Volumes().IsLocal(v) && status != nil {
				opts.Volumes[v] = getVolumeOptions(status)
			} else {
				delete(c.cache.volumes, v)
			}
		}

		for p := range opts.Pods {
			delete(c.cache.pods, p)
		}

		for v := range opts.Volumes {
			delete(c.cache.volumes, v)
		}

		c.cache.lock.Unlock()

		spec, err := envs.Get().GetNodeClient().SetStatus(ctx, opts)
		if err != nil {
			log.Errorf("%s node:exporter:dispatch err: %s", logPrefix, err.Error())
		}

		if spec != nil {
			if err := c.runtime.Sync(ctx, spec.Decode()); err != nil {
				log.Errorf("%s runtime sync err: %s", logPrefix, err.Error())
			}
		} else {
			log.Debugf("%s received spec is nil, skip apply changes", logPrefix)
		}
	}

	return nil
}

func (c *Controller) Subscribe() {
	var (
		pods    = make(chan string)
		volumes = make(chan string)
		done    = make(chan bool)
	)

	go func() {
		log.Debugf("%s subscribe state", logPrefix)

		for {
			select {
			case p := <-pods:
				log.Debugf("%s pod changed: %s", logPrefix, p)
				c.cache.lock.Lock()
				c.cache.pods[p] = envs.Get().GetState().Pods().GetPod(p)
				c.cache.lock.Unlock()
				break
			case v := <-volumes:
				log.Debugf("%s volume changed: %s", logPrefix, v)
				c.cache.lock.Lock()
				c.cache.volumes[v] = envs.Get().GetState().Volumes().GetVolume(v)
				c.cache.lock.Unlock()
				break
			}
		}

	}()

	go envs.Get().GetState().Pods().Watch(pods, done)
	go envs.Get().GetState().Volumes().Watch(volumes, done)

	<-done
}

func getPodOptions(p *types.PodStatus) *request.NodePodStatusOptions {
	opts := v1.Request().Node().NodePodStatusOptions()
	opts.State = p.State
	opts.Status = p.Status
	opts.Running = p.Running
	opts.Message = p.Message
	opts.Runtime = p.Runtime
	opts.Network = p.Network
	opts.Steps = p.Steps
	return opts
}

func getVolumeOptions(p *types.VolumeStatus) *request.NodeVolumeStatusOptions {
	opts := v1.Request().Node().NodeVolumeStatusOptions()
	opts.State = p.State
	opts.Message = p.Message
	return opts
}
