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

package exporter

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

const (
	logPrefix = "controller:>"
	logLevel  = 3
)

type Exporter struct {
	cache struct {
		lock        sync.RWMutex
		cluster     *types.Cluster
		nodes       map[string]*types.Node
		namespaces  map[string]*types.Namespace
		services    map[string]*types.Service
		deployments map[string]*types.Deployment
		pods        map[string]*types.Pod
		jobs				map[string]*types.Job
		tasks				map[string]*types.Task
		volumes     map[string]*types.Volume
		routes      map[string]*types.Route
	}
}

func New() *Exporter {
	var c = new(Exporter)
	return c
}

func (c *Exporter) Connect(ctx context.Context) error {
	log.V(logLevel).Debugf("%s:connect:> connect init", logPrefix)

	return nil
}

func (c *Exporter) SendClusterState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendNodeState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendNamespaceState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendServiceState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendDeploymentState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendPodState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendVolumeState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendJobState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendTaskState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendRouteState(ctx context.Context) error {
	return nil
}
