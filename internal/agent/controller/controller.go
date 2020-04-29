//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"fmt"
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/internal/agent/config"
	"github.com/lastbackend/lastbackend/internal/agent/runtime"
	"github.com/lastbackend/lastbackend/internal/agent/state"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/client/cluster"
	client "github.com/lastbackend/lastbackend/pkg/client/cluster"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "client:>"
)

type Controller struct {
	ctx        context.Context
	runtime    *runtime.Runtime
	state      *state.State
	restClient cluster.IClient
	network    *network.Network
	cache      struct {
		lock      sync.RWMutex
		resources models.NodeStatus
		pods      map[string]*models.PodStatus
		volumes   map[string]*models.VolumeStatus
	}
}

func New(r *runtime.Runtime) (*Controller, error) {

	cfg := client.NewConfig()
	cfg.BearerToken = r.GetConfig().Security.Token

	if r.GetConfig().API.TLS.Verify {
		cfg.TLS = client.NewTLSConfig()
		cfg.TLS.Insecure = !r.GetConfig().API.TLS.Verify
		cfg.TLS.CertFile = r.GetConfig().API.TLS.FileCert
		cfg.TLS.KeyFile = r.GetConfig().API.TLS.FileKey
		cfg.TLS.CAFile = r.GetConfig().API.TLS.FileCA
	}

	endpoint := r.GetConfig().API.Address

	rest, err := client.New(client.ClientHTTP, endpoint, cfg)
	if err != nil {
		return nil, fmt.Errorf("Can not initialize http client: %v", err)
	}

	var c = new(Controller)
	c.ctx = context.Background()
	c.runtime = r
	c.restClient = rest
	c.state = r.GetState()
	c.network = r.GetNetwork()
	c.cache.pods = make(map[string]*models.PodStatus)
	c.cache.volumes = make(map[string]*models.VolumeStatus)

	pods := c.state.Pods().GetPods()
	for p, st := range pods {
		c.cache.pods[p] = st
	}

	return c, nil
}

func (c *Controller) Connect(cfg config.Config) error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s:connect:> connect init", logPrefix)

	opts := v1.Request().Node().NodeConnectOptions()
	opts.Info = c.state.Node().Info

	opts.Status = c.state.Node().Status

	if c.network != nil {
		opts.Network = c.network.Info(c.ctx)
	}

	// TODO: Implement tls logic for rest client
	//if v.IsSet("node.tls") {
	//	opts.TLS = !v.GetBool("node.tls.insecure")
	//
	//	if opts.TLS {
	//		caData, err := ioutil.ReadFile(v.GetString("node.tls.ca"))
	//		if err != nil {
	//			log.Errorf("%s:connect_event:> read ca cert file err: %v", logPrefix, err)
	//			return err
	//		}
	//
	//		certData, err := ioutil.ReadFile(v.GetString("node.tls.client_cert"))
	//		if err != nil {
	//			log.Errorf("%s:connect_event:> read client cert file err: %v", logPrefix, err)
	//			return err
	//		}
	//
	//		keyData, err := ioutil.ReadFile(v.GetString("node.tls.client_key"))
	//		if err != nil {
	//			log.Errorf("%s:connect_event:> read client key file err: %v", logPrefix, err)
	//			return err
	//		}
	//
	//		opts.SSL = new(request.SSL)
	//		opts.SSL.CA = caData
	//		opts.SSL.Key = keyData
	//		opts.SSL.Cert = certData
	//	}
	//}

	for {
		log.Debugf("%s:connect:> establish connection", logPrefix)
		if err := c.restClient.V1().Cluster().Node(c.state.Node().Info.Hostname).Connect(c.ctx, opts); err == nil {
			return nil
		} else {
			log.Errorf("%s:connect:> establish connection err: %s", logPrefix, err.Error())
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *Controller) Sync() error {
	log := logger.WithContext(context.Background())
	log.Debugf("%s start node sync", logPrefix)

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		opts := new(request.NodeStatusOptions)
		opts.Pods = make(map[string]*request.NodePodStatusOptions)
		opts.Volumes = make(map[string]*request.NodeVolumeStatusOptions)

		opts.Resources.Capacity = c.state.Node().Status.Capacity
		opts.Resources.Allocated = c.state.Node().Status.Allocated

		c.cache.lock.Lock()
		var i = 0
		for p, status := range c.cache.pods {
			i++
			if i > 10 {
				break
			}

			if !c.state.Pods().IsLocal(p) && status != nil {
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

			if !c.state.Volumes().IsLocal(v) && status != nil {
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

		spec, err := c.restClient.V1().Cluster().Node(c.state.Node().Info.Hostname).SetStatus(c.ctx, opts)
		if err != nil {
			log.Errorf("%s node:exporter:dispatch err: %s", logPrefix, err.Error())
		}

		if spec != nil {
			if err := c.runtime.Sync(spec.Decode()); err != nil {
				log.Errorf("%s runtime sync err: %s", logPrefix, err.Error())
			}
		} else {
			log.Debugf("%s received spec is nil, skip apply changes", logPrefix)
		}
	}

	return nil
}

func (c *Controller) Subscribe() {
	log := logger.WithContext(context.Background())
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
				c.cache.pods[p] = c.state.Pods().GetPod(p)
				c.cache.lock.Unlock()
				break
			case v := <-volumes:
				log.Debugf("%s volume changed: %s", logPrefix, v)
				c.cache.lock.Lock()
				c.cache.volumes[v] = c.state.Volumes().GetVolume(v)
				c.cache.lock.Unlock()
				break
			}
		}

	}()

	go c.state.Pods().Watch(pods, done)
	go c.state.Volumes().Watch(volumes, done)

	<-done
}

func getPodOptions(p *models.PodStatus) *request.NodePodStatusOptions {
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

func getVolumeOptions(p *models.VolumeStatus) *request.NodeVolumeStatusOptions {
	opts := v1.Request().Node().NodeVolumeStatusOptions()
	opts.State = p.State
	opts.Message = p.Message
	return opts
}
